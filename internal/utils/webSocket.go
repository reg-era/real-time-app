package utils

import (
	"fmt"
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID   string
	Conn *websocket.Conn
}

type Action struct {
	To   string
	Data []byte
}

type Pool struct {
	Connections map[string]*Connection
	Channel     chan Action
	Mu          sync.RWMutex
}

func NewPool() *Pool {
	return &Pool{
		Connections: make(map[string]*Connection),
		Channel:     make(chan Action),
	}
}

func (p *Pool) AddConn(conn *websocket.Conn, id string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	p.Connections[id] = &Connection{
		ID:   id,
		Conn: conn,
	}
	log.Println("Added connection with ID:", id)
}

func (p *Pool) RemoveConn(id string) {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	if _, exists := p.Connections[id]; exists {
		delete(p.Connections, id)
		log.Println("Removed connection with ID:", id)
	}
}

func (p *Pool) GetConn(id string) (*Connection, bool) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	conn, exists := p.Connections[id]
	return conn, exists
}

func (p *Pool) PingConnections() {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	for id, conn := range p.Connections {
		err := conn.Conn.WriteMessage(websocket.PingMessage, nil)
		if err != nil {
			log.Println("Ping failed for connection with ID:", id, err)
			conn.Conn.Close()
			delete(p.Connections, id)
		}
	}
}

func WebSocketHandler(pool *Pool) {
	for cmd := range pool.Channel {
		pool.Mu.RLock()
		conn, exist := pool.GetConn(cmd.To)
		if exist {
			err := conn.Conn.WriteMessage(websocket.TextMessage, cmd.Data)
			if err != nil {
				fmt.Println("Error sending message:", err)
			}
		}
		pool.Mu.RUnlock()
	}
}
