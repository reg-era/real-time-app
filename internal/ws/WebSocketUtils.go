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
	Id          int       `json: "Id"`
	Name        string    `json: "Name"`
	LastMessage string    `json: "LastMessage"`
	Time        time.Time `json: "Time"`
	Online      bool      `json: "Online"`
}

type Users struct {
	Friends []Friend `json:"Friends"`
}

type WebSocketMessage struct {
	Type  string `json:"Type"`
	Users Users  `json:"users"`
}

type websocketmsg struct {
	Type    string        `json: "Type"`
	Message utils.Message `json:"Message"`
}

type Hub struct {
	Clients    map[*Client]int
	Broadcast  chan *sql.DB
	Message    chan utils.Message
	Register   chan *Client
	Unregister chan *Client
	Logout     chan *Client
	Mutex      sync.RWMutex
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = client.Id
			h.Mutex.Unlock()

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.Mutex.Unlock()

		case message := <-h.Broadcast:
			h.Mutex.RLock()
			for client := range h.Clients {
				correctmessage, _ := Getuserslist(client, h, message)
				mssg, err := json.Marshal(correctmessage)
				if err != nil {
					log.Printf("Error marshling: %v", err)
					client.Conn.Close()
					delete(h.Clients, client)
				}
				err = client.Conn.WriteMessage(websocket.TextMessage, mssg)
				if err != nil {
					log.Printf("Error broadcasting to client11: %v", err)
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.Mutex.RUnlock()
		case mssg := <-h.Message:
			response := websocketmsg{
				Type:    "message",
				Message: mssg,
			}
			h.Mutex.RLock()
			for client := range h.Clients {
				if client.Id == mssg.ReceiverID {
					data, _ := json.Marshal(response)
					err := client.Conn.WriteMessage(websocket.TextMessage, data)
					if err != nil {
						log.Printf("Error broadcasting to client: %v", err)
						client.Conn.Close()
						delete(h.Clients, client)
					}
				}
			}
			h.Mutex.RUnlock()
		case tologout := <-h.Logout:
			response := websocketmsg{
				Type:    "Logout",
				Message: utils.Message{},
			}
			data, _ := json.Marshal(response)
			err := tologout.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Error broadcasting to client: %v", err)
				tologout.Conn.Close()
				delete(h.Clients, tologout)
			}
		}
	}
}

func Getuserslist(client *Client, hub *Hub, db *sql.DB) (WebSocketMessage, error) {
	friendsIds, err := database.GetFriends(db, client.Id)
	if err != nil {
		return WebSocketMessage{}, fmt.Errorf("failed to get friends: %v", err)
	}

	err, allFriends := creatfriendslist(friendsIds, client.Id, db)
	if err != nil {
		return WebSocketMessage{}, fmt.Errorf("failed to get friends1: %v", err)
	}
	SortByLastMessage(allFriends)

	for client := range hub.Clients {
		for i, friend := range allFriends {
			if client.Id == friend.Id {
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

func creatfriendslist(allFriends []int, userId int, db *sql.DB) (error, []Friend) {
	Friends := []Friend{}
	for i := 0; i < len(allFriends); i++ {
		friend := Friend{}
		err, mssg := database.Getlastmessg(userId, allFriends[i], db)
		if err != nil && err != sql.ErrNoRows {
			return err, nil
		}
		friend.LastMessage = mssg.Message
		friend.Time = mssg.CreatedAt
		friend.Online = false
		if mssg.Message != "" {
			if mssg.SenderID == userId {
				friend.Id = mssg.ReceiverID
				friend.Name, err = database.GetUserName(mssg.ReceiverID, db)
				if err != nil {
					return err, nil
				}
			} else {
				friend.Name, err = database.GetUserName(mssg.SenderID, db)
				friend.Id = mssg.SenderID
				if err != nil {
					return err, nil
				}
			}
		} else {
			friend.Name, err = database.GetUserName(allFriends[i], db)
			friend.Id = allFriends[i]
			if err != nil {
				return err, nil
			}
		}
		Friends = append(Friends, friend)
	}
	return nil, Friends
}
