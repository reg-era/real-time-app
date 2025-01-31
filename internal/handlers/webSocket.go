package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"forum/internal/database"
	"forum/internal/utils"
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

	// Create client
	newclient := &utils.Client{
		Id:   userid,
		Conn: conn,
	}
	hub.Register <- newclient
	defer func() {
		hub.Unregister <- newclient
		if mssg, err := getuserslist(newclient, hub, db); err != nil {
			fmt.Printf("Failed to send users list: %v", err)
		} else {
			hub.Broadcast <- mssg
		}
	}()

	if mssg, err := getuserslist(newclient, hub, db); err != nil {
		fmt.Printf("Failed to send users list: %v", err)
	} else {
		hub.Broadcast <- mssg
	}

	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

func getuserslist(client *utils.Client, hub *utils.Hub, db *sql.DB) ([]byte, error) {
	allFriends, err := database.GetFriends(db, client.Id)
	fmt.Println("all friends:", allFriends)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %v", err)
	}

	var onlineFriends, offlineFriends []string

	hub.Mutex.Lock()
	for c := range hub.Clients {
		clientname, err := database.GetUserName(c.Id, db)
		if err != nil {
			continue
		}
		onlineFriends = append(onlineFriends, clientname)
	}
	for _, client := range allFriends {
		for _, onligne := range onlineFriends {
			if client != onligne {
				offlineFriends = append(offlineFriends, client)
			}
		}
	}
	hub.Mutex.Unlock()

	response := WebSocketMessage{
		Type: "onlineusers",
		Users: users{
			Online:  onlineFriends,
			Offline: offlineFriends,
		},
	}

	jsonResponse, err := json.Marshal(response)
	if err != nil {
		return nil, fmt.Errorf("JSON error: %v", err)
	}
	return jsonResponse, nil
}

// go readmessages()
// read from this connection and store in db and if receiver connected send it to writemssg in channel
// go writemessages()
// write this message to this clients if connected
