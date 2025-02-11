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

func HandleWs(w http.ResponseWriter, r *http.Request, db *sql.DB, userid int, hub *websocket.Hub) {
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
	// time.Sleep(100 * time.Millisecond)

	defer func() {
		hub.Unregister <- newclient
		hub.Broadcast <- db
	}()

	hub.Broadcast <- db

	type mssge struct {
		ReceiverName string `json:"ReceiverName"`
		Data         string `json:"Data"`
	}

	lastmessage := time.Now()
	for {
		_, mssg, err := conn.ReadMessage()
		if err != nil {
			break
		}

		Duration := time.Since(lastmessage)
		if Duration.Milliseconds() > 1 {
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
			if _, exist := hub.Clients[id]; exist {
				newmssg.Seen = 1
			} else {
				newmssg.Seen = 0
			}
			newmssg.SenderName, err = database.GetUserName(userid, db)
			if err != nil {
				break
			}
			newmssg.ReceiverID = id
			database.CreateMessage(&newmssg, db)
			hub.Message <- newmssg
			hub.Broadcast <- db
			lastmessage = time.Now()
		} else {
			// error spaming messages
			hub.Unregister <- newclient
			hub.Broadcast <- db
		}
	}
}

func checkForValue(userValue int, users map[int]*websocket.Client) (bool, *websocket.Client) {
	for c, value := range users {
		if c == userValue {
			return true, value
		}
	}
	return false, nil
}
