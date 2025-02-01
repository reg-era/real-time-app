package utils

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512
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
	Send chan Message
	// Target *Client
}

// var Clients map[*Client]bool

type Hub struct {
	Clients    map[*Client]int
	Broadcast  chan []byte
	Message    chan Message
	Register   chan *Client
	Unregister chan *Client
	Mutex      sync.RWMutex
}

type Message struct {
	Id         int    `Id`
	SenderID   int    `SenderID`
	ReceiverID int    `ReceiverID`
	Message    string `Message`
	CreatedAt  string `CreatedAt`
	IsSender   bool   `IsSender`
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			h.Mutex.Lock()
			h.Clients[client] = client.Id
			h.Mutex.Unlock()
			// Notify all Clients about new user
			// message := []byte("New user joined: " + string(client.Id))
			// h.Broadcast <- message

		case client := <-h.Unregister:
			h.Mutex.Lock()
			if _, ok := h.Clients[client]; ok {
				delete(h.Clients, client)
				client.Conn.Close()
			}
			h.Mutex.Unlock()

		case message := <-h.Broadcast:
			fmt.Println(string(message))
			h.Mutex.RLock()
			for client := range h.Clients {
				err := client.Conn.WriteMessage(websocket.TextMessage, message)
				if err != nil {
					log.Printf("Error broadcasting to client: %v", err)
					client.Conn.Close()
					delete(h.Clients, client)
				}
			}
			h.Mutex.RUnlock()
		}
	}
}

// func (c *Client) WritePump() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()

// 	for {
// 		select {
// 		case message, ok := <-c.Send:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if !ok {
// 				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
// 				return
// 			}

// 			w, err := c.Conn.NextWriter(websocket.TextMessage)
// 			if err != nil {
// 				return
// 			}
// 			w.Write([]byte(message.Message))

// 			if err := w.Close(); err != nil {
// 				return
// 			}
// 		case <-ticker.C:
// 			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
// 				return
// 			}
// 		}
// 	}
// }
