package utils

import (
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

type Connection struct {
	ID       string
	Conn     *websocket.Conn
	LastSeen time.Time
}

type Action struct {
	Type string
	Data interface{}
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
		ID:       id,
		Conn:     conn,
		LastSeen: time.Now(),
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

func (p *Pool) Broadcast(message []byte) {
	p.Mu.RLock()
	defer p.Mu.RUnlock()

	for _, conn := range p.Connections {
		if err := conn.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			log.Println("Error broadcasting to connection:", conn.ID, err)
		}
	}
}

func (p *Pool) Cleanup() {
	p.Mu.Lock()
	defer p.Mu.Unlock()

	for id, conn := range p.Connections {
		if time.Since(conn.LastSeen) > 30*time.Minute {
			conn.Conn.Close()
			delete(p.Connections, id)
			log.Println("Removed stale connection with ID:", id)
		}
	}
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