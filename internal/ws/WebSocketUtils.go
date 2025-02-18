package utils

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"forum/internal/database"
	"forum/internal/utils"

	"github.com/gorilla/websocket"
)

var Upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	Id   int
	Conn *websocket.Conn
}

type Friend struct {
	Id          int       `json:"Id"`
	Name        string    `json:"Name"`
	LastMessage string    `json:"LastMessage"`
	Time        time.Time `json:"Time"`
	Online      bool      `json:"Online"`
	Seen        int       `json:"Seen"`
	IsSender    bool      `json:"IsSender"`
}

type Users struct {
	Friends []Friend `json:"Friends"`
}

type WebSocketMessage struct {
	Type  string `json:"Type"`
	Users Users  `json:"users"`
}

type websocketmsg struct {
	Type    string        `json:"Type"`
	Message utils.Message `json:"Message"`
}

type Hub struct {
	Clients    map[int][]*Client
	Broadcast  chan *sql.DB
	Message    chan utils.Message
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex

	Progress chan Progresser
}

type Progresser struct {
	UserId   int
	UserName string `json:"user_name"`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client.Id] = append(h.Clients[client.Id], client)
			h.Mutex.Unlock()
		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client.Id]; ok {
				client.Conn.Close()
				h.Clients[client.Id] = DeleteUser(h.Clients[client.Id], client)
				if len(h.Clients[client.Id]) == 0 {
					delete(h.Clients, client.Id)
				}
			}
			h.Mutex.Unlock()
		case message := <-h.Broadcast:
			h.Mutex.RLock()
			for client := range h.Clients {
				correctmessage, _ := Getuserslist(h.Clients[client][0], h, message)
				mssg, err := json.Marshal(correctmessage)
				if err != nil {
					log.Printf("Error marshling: %v", err)
				}
				for _, window := range h.Clients[client] {
					err = window.Conn.WriteMessage(websocket.TextMessage, mssg)
					if err != nil {
						log.Printf("Error broadcasting to client11: %v", err)
						window.Conn.Close()
						delete(h.Clients, client)
					}
				}
				database.Updatesenn(client, message)
			}
			h.Mutex.RUnlock()
		case mssg := <-h.Message:
			response := websocketmsg{
				Type:    "message",
				Message: mssg,
			}
			h.Mutex.RLock()
			if client, ok := h.Clients[mssg.ReceiverID]; ok {
				data, err := json.Marshal(response)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err) // handling error
				}
				for _, window := range client {
					err = window.Conn.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						log.Printf("Error broadcasting to client: %v", err) // handling error
					}
				}
			}
			h.Mutex.RUnlock()

		case typer := <-h.Progress:
			// fmt.Println(len(h.Clients[typer.UserId]))
			if Client, ok := h.Clients[typer.UserId]; ok {
				for _, window := range Client {
					response := websocketmsg{
						Type: "inprogress",
						Message: utils.Message{
							SenderName: typer.UserName,
						},
					}
					err := window.Conn.WriteJSON(response)
					if err != nil {
						log.Printf("Error broadcasting to client: %v", err)
					}
				}
			}
		}
	}
}

func Getuserslist(client *Client, hub *Hub, db *sql.DB) (WebSocketMessage, error) {
	friendsIds, err := database.GetFriends(db, client.Id)
	if err != nil {
		return WebSocketMessage{}, fmt.Errorf("failed to get friends: %v", err)
	}

	allFriends, err := creatfriendslist(friendsIds, client.Id, db)
	if err != nil {
		return WebSocketMessage{}, fmt.Errorf("failed to get friends1: %v", err)
	}

	SortByLastMessage(allFriends)

	for client := range hub.Clients {
		for i, friend := range allFriends {
			if client == friend.Id {
				allFriends[i].Online = true
			}
		}
	}

	response := WebSocketMessage{
		Type: "onlineusers",
		Users: Users{
			Friends: allFriends,
		},
	}

	return response, nil
}

func SortByLastMessage(allConversations []Friend) {
	n := len(allConversations)
	for i := 0; i < n-1; i++ {
		for j := 0; j < n-1-i; j++ {
			shouldSwap := false
			if allConversations[j].LastMessage != "" && allConversations[j+1].LastMessage != "" {
				shouldSwap = allConversations[j].Time.Before(allConversations[j+1].Time)
			} else if allConversations[j].LastMessage == "" && allConversations[j+1].LastMessage == "" {
				shouldSwap = allConversations[j].Name > allConversations[j+1].Name
			} else {
				shouldSwap = allConversations[j].LastMessage == "" && allConversations[j+1].LastMessage != ""
			}

			if shouldSwap {
				allConversations[j], allConversations[j+1] = allConversations[j+1], allConversations[j]
			}
		}
	}
}

func creatfriendslist(allFriends []int, userId int, db *sql.DB) ([]Friend, error) {
	Friends := []Friend{}
	for i := 0; i < len(allFriends); i++ {
		friend := Friend{}
		err, mssg := database.Getlastmessg(userId, allFriends[i], db)
		if err != nil && err != sql.ErrNoRows {
			return nil, err
		}
		friend.LastMessage = mssg.Message
		friend.Time = mssg.CreatedAt
		friend.Online = false
		if mssg.Message == "" {
			friend.Seen = 1
		} else {
			friend.Seen = mssg.Seen
		}
		if mssg.Message != "" {
			if mssg.SenderID == userId {
				friend.Id = mssg.ReceiverID
				friend.Name, err = database.GetUserName(mssg.ReceiverID, db)
				if err != nil {
					return nil, err
				}
				friend.IsSender = true
			} else {
				friend.Name, err = database.GetUserName(mssg.SenderID, db)
				friend.Id = mssg.SenderID
				if err != nil {
					return nil, err
				}
				friend.IsSender = false
			}
		} else {
			friend.Name, err = database.GetUserName(allFriends[i], db)
			friend.Id = allFriends[i]
			if err != nil {
				return nil, err
			}
		}
		Friends = append(Friends, friend)
	}
	return Friends, nil
}

func DeleteUser(arr []*Client, user *Client) []*Client {
	for index, client := range arr {
		if client.Conn == user.Conn {
			return append(arr[:index], arr[index+1:]...)
		}
	}
	return arr
}
