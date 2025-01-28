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

func HandleWs(w http.ResponseWriter, r *http.Request, userid int, db *sql.DB, hub *utils.Hub) {
	conn, err := utils.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade: %v", err)
		return
	}
	defer conn.Close()
	type users struct {
		Online  []string `json:"online"`
		Offline []string `json:"offline"`
	}
	newclient := &utils.Client{
		Id:   userid,
		Conn: conn,
	}
	type WebSocketMessage struct {
		Type  string `json:"Type"`
		Users users  `json:"users"`
	}
	hub.Mu.Lock()
	hub.Clients[newclient] = true
	hub.Mu.Unlock()

	// Get all friends
	allFriends, err := database.GetFriends(db, userid)
	if err != nil {
		log.Printf("Failed to get friends: %v", err)
		return
	}

	// Create maps for O(1) lookups
	onlineFriends := make([]string, 0)
	offlineFriends := make([]string, 0)

	hub.Mu.Lock()
	fmt.Println(hub.Clients)
	// Check each friend's online status
	for _, friend := range allFriends {
		isOnline := false
		for client := range hub.Clients {
			clientname, err := database.GetUserName(client.Id, db)
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
	// Create response structure
	Users1 := users{
		Online:  onlineFriends,
		Offline: offlineFriends,
	}

	jsonResponse, err := json.Marshal(WebSocketMessage{
		Type:  "onlineusers",
		Users: Users1,
	})
	fmt.Println(string(jsonResponse))
	if err != nil {
		log.Printf("JSON error: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonResponse); err != nil {
		log.Printf("Write error: %v", err)
		return
	}
}

// go readmessages()
// read from this connection and store in db and if receiver connected send it to writemssg in channel
// go writemessages()
// write this message to this clients if connected
