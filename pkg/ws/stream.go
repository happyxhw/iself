// Package ws: websocket
package ws

import (
	"errors"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var (
	// ErrClientOffline client offline
	ErrClientOffline = errors.New("client offline")
)

// Msg info sent to Client
type Msg struct {
	Timestamp int64 `json:"timestamp"`
}

// Stream websocket srv
type Stream struct {
	m       sync.RWMutex
	clients map[int64]map[string]*Client

	upgrader *websocket.Upgrader
}

// NewStream creates a new instance of steam.
func NewStream() *Stream {
	return &Stream{
		clients: make(map[int64]map[string]*Client),
		upgrader: &websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin: func(r *http.Request) bool {
				return true
			},
		},
	}
}

// Upgrade conn to websocket
func (s *Stream) Upgrade(w http.ResponseWriter, r *http.Request, h http.Header) (*websocket.Conn, error) {
	return s.upgrader.Upgrade(w, r, h)
}

// SendMsg notifies the clients with the given userID that a new messages was created.
func (s *Stream) SendMsg(userID int64, msg *Msg) error {
	s.m.RLock()
	defer s.m.RUnlock()
	clients, ok := s.clients[userID]
	if !ok || len(clients) == 0 {
		return ErrClientOffline
	}
	for _, c := range clients {
		_ = c.send(msg)
	}
	return nil
}

// Register a client
func (s *Stream) Register(c *Client) {
	s.m.Lock()
	defer s.m.Unlock()
	// close first
	if clients, ok := s.clients[c.userID]; ok {
		if cli, cOk := clients[c.id]; cOk {
			cli.close()
		}
	} else {
		s.clients[c.userID] = make(map[string]*Client)
	}
	s.clients[c.userID][c.id] = c
}

// Remove a client
func (s *Stream) Remove(c *Client) {
	s.m.Lock()
	defer s.m.Unlock()
	delete(s.clients[c.userID], c.id)
	if len(s.clients[c.userID]) == 0 {
		delete(s.clients, c.userID)
	}
}
