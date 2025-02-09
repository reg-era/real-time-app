package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"forum/internal/database"
	"forum/internal/utils"
	websocket "forum/internal/ws"
)

func HandleWs(w http.ResponseWriter, r *http.Request, userid int, db *sql.DB, hub *websocket.Hub) {
	conn, err := websocket.Upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("Failed to upgrade: %v", err)
		return
	}
	hub.Mutex.Lock()
	if exist, client := checkForValue(userid, hub.Clients); exist {
		hub.Logout <- client
		hub.Unregister <- client
	}
	hub.Mutex.Unlock()

	newclient := &websocket.Client{
		Id:   userid,
		Conn: conn,
	}
	hub.Register <- newclient

	// locking map to check if user exist slow the regitration process
	time.Sleep(100 * time.Millisecond)

	defer func() {
		hub.Unregister <- newclient
		hub.Broadcast <- db
	}()

	hub.Broadcast <- db

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
		newmssg.CreatedAt = time.Now()
		id, err := database.GetUserIdByName(received.ReceiverName, db)
		if err != nil {
			break
		}
		newmssg.SenderName, err = database.GetUserName(userid, db)
		if err != nil {
			break
		}
		newmssg.ReceiverID = id
		database.CreateMessage(&newmssg, db)
		hub.Message <- newmssg
		hub.Broadcast <- db
	}
}

func checkForValue(userValue int, users map[*websocket.Client]int) (bool, *websocket.Client) {
	for c, value := range users {
		if value == userValue {
			return true, c
		}
	}
	return false, nil
}
