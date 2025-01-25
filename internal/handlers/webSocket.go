package handlers

import (
	"database/sql"
	"encoding/json"
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

	newclient := &utils.Client{
		Id:   userid,
		Conn: conn,
	}

	hub.Mu.Lock()
	hub.Clients[newclient] = true
	hub.Mu.Unlock()

	var clients []string
	for client := range hub.Clients {
		clientname, err := database.GetUserName(client.Id, db)
		if err != nil {
			continue
		}
		clients = append(clients, clientname)
	}

	jsonNames, err := json.Marshal(clients)
	if err != nil {
		log.Printf("JSON error: %v", err)
		return
	}

	if err := conn.WriteMessage(websocket.TextMessage, jsonNames); err != nil {
		log.Printf("Write error: %v", err)
		return
	}
}

// go readmessages()
// read from this connection and store in db and if receiver connected send it to writemssg in channel
// go writemessages()
// write this message to this clients if connected
