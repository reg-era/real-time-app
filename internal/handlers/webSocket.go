package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

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
	hub.Mutex.Lock()
	if exist, client := checkForValue(userid, hub.Clients); exist {
		hub.Unregister <- client
	}
	hub.Mutex.Unlock()
	// if exist , ok := hub.Clients[]
	// Create client
	newclient := &utils.Client{
		Id:   userid,
		Conn: conn,
	}
	hub.Register <- newclient

	// locking map to check if user exist slow the regitration process
	time.Sleep(100 * time.Millisecond)

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

	type mssge struct {
		ReceiverName string `json:"ReceiverName"`
		Data         string `json:"Data"`
	}

	for {
		_, mssg, err := conn.ReadMessage()
		if err != nil {
			break
		}
		newmssg := utils.Message{}
		received := mssge{}
		json.Unmarshal(mssg, &received)
		newmssg.Message = received.Data
		newmssg.SenderID = userid
		newmssg.CreatedAt = time.Now().String()
		id, err := database.GetUserIdByName(received.ReceiverName, db)
		if err != nil {
			break
		}
		newmssg.ReceiverID = id
		database.CreateMessage(&newmssg, db)
		hub.Message <- newmssg
	}
}

func getuserslist(client *utils.Client, hub *utils.Hub, db *sql.DB) ([]byte, error) {
	allFriends, err := database.GetFriends(db, client.Id)
	if err != nil {
		return nil, fmt.Errorf("failed to get friends: %v", err)
	}

	onlineMap := make(map[string]bool)
	var onlineFriends []string

	hub.Mutex.Lock()
	for c := range hub.Clients {
		clientname, err := database.GetUserName(c.Id, db)
		if err != nil {
			continue
		}
		onlineFriends = append(onlineFriends, clientname)
		onlineMap[clientname] = true
	}
	hub.Mutex.Unlock()

	// Build offline friends list without duplicates
	var offlineFriends []string
	for _, friend := range allFriends {
		if !onlineMap[friend] {
			offlineFriends = append(offlineFriends, friend)
		}
	}

	response := WebSocketMessage{
		Type: "onlineusers",
		Users: users{
			Online:  onlineFriends,
			Offline: offlineFriends,
		},
	}

	return json.Marshal(response)
}

func checkForValue(userValue int, users map[*utils.Client]int) (bool, *utils.Client) {
	for c, value := range users {
		if value == userValue {
			return true, c
		}
	}

	return false, nil
}

// go readmessages()
// read from this connection and store in db and if receiver connected send it to writemssg in channel
// go writemessages()
// write this message to this clients if connected
