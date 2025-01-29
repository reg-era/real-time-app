package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/utils"

	"github.com/gorilla/websocket"
)

type WebSocketMessage struct {
	Type  string `json:"Type"`
	Users users  `json:"users"`
}

type users struct {
	Online  []string `json:"online"`
	Offline []string `json:"offline"`
}

func HandleWs(w http.ResponseWriter, r *http.Request, userid int, db *sql.DB, hub *utils.Hub) {
	conn, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade: %v", err)
		return
	}
	defer conn.Close()

	// Create client
	newclient := &utils.Client{
		Id:   userid,
		Conn: conn,
	}

	// Add client to hub
	hub.Mu.Lock()
	hub.Clients[newclient] = true
	hub.Mu.Unlock()
	// Cleanup when function returns
	defer func() {
		hub.Mu.Lock()
		delete(hub.Clients, newclient)
		hub.Mu.Unlock()
	}()

	if err := sendUsersList(newclient, hub, db); err != nil {
		fmt.Printf("Failed to send users list: %v", err)
	}

	cont := 0
	for {
		cont++
	}
}

func sendUsersList(client *utils.Client, hub *utils.Hub, db *sql.DB) error {
	allFriends, err := database.GetFriends(db, client.Id)
	if err != nil {
		return fmt.Errorf("failed to get friends: %v", err)
	}

	var onlineFriends, offlineFriends []string

	hub.Mu.Lock()
	for _, friend := range allFriends {
		isOnline := false
		for c := range hub.Clients {
			clientname, err := database.GetUserName(c.Id, db)
			if err != nil {
				continue
			}
			if friend == clientname {
				onlineFriends = append(onlineFriends, clientname)
				isOnline = true
				break
			}
		}
		if !isOnline {
			offlineFriends = append(offlineFriends, friend)
		}
	}
	hub.Mu.Unlock()

	response := WebSocketMessage{
		Type: "onlineusers",
		Users: users{
			Online:  onlineFriends,
			Offline: offlineFriends,
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return fmt.Errorf("JSON error: %v", err)
	}
	fmt.Println(string(jsonResponse))
	hub.Mu.Lock()
	for client := range hub.Clients {
		err := client.Conn.WriteMessage(websocket.TextMessage, jsonResponse)
		if err != nil {
			return err
		}
	}
	hub.Mu.Unlock()

	return nil
}

// go readmessages()
// read from this connection and store in db and if receiver connected send it to writemssg in channel
// go writemessages()
// write this message to this clients if connected
