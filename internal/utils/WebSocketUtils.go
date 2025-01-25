package utils

import (
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

var Clients map[*Client]bool

type Hub struct {
	Clients    map[*Client]bool
	Broadcast  chan Message
	Register   chan *Client
	Unregister chan *Client
	Mu         sync.Mutex
}

type Message struct {
	Id         int    `Id`
	SenderID   int    `SenderID`
	ReceiverID int    `ReceiverID`
	Message    string `Message`
	CreatedAt  string `CreatedAt`
	IsSender   bool   `IsSender`
}

// func (h *Hub) Run() {
// 	for {
// 		select {
// 		case client := <-h.Register:
// 			h.Mu.Lock()
// 			h.Clients[client] = true
// 			h.Mu.Unlock()
// 		case client := <-h.Unregister:
// 			h.Mu.Lock()
// 			if _, ok := h.Clients[client]; ok {
// 				delete(h.Clients, client)
// 				close(client.Send)
// 			}
// 			h.Mu.Unlock()
// 		case Message := <-h.Broadcast:
// 			h.Mu.Lock()
// 			for client := range h.Clients {
// 				select {
// 				case client.Send <- Message:
// 				default:
// 					close(client.Send)
// 					delete(h.Clients, client)
// 				}
// 			}
// 			h.Mu.Unlock()
// 		}
// 	}
// }

// func (c *Client) ReadPump() {
// 	defer func() {
// 		c.Hub.Unregister <- c
// 		c.Conn.Close()
// 	}()

// 	c.Conn.SetReadLimit(maxMessageSize)
// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 	c.Conn.SetPongHandler(func(string) error {
// 		c.Conn.SetReadDeadline(time.Now().Add(pongWait))
// 		return nil
// 	})

// 	for {
// 		_, message, err := c.Conn.ReadMessage()
// 		if err != nil {
// 			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
// 				log.Printf("error: %v", err)
// 			}
// 			break
// 		}
// 		c.Hub.Broadcast <- message
// 	}
// }

func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.Conn.Close()
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := c.Conn.NextWriter(websocket.TextMessage)
			if err != nil {
				return
			}
			w.Write([]byte(message.Message))

			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}
