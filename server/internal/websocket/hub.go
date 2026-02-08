package websocket

import (
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"pinche/internal/logger"
)

type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}

type Client struct {
	UserID uint64
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients    map[uint64]*Client
	register   chan *Client
	unregister chan *Client
	mu         sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[uint64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			h.mu.Unlock()
			logger.Info("WebSocket client registered", "user_id", client.UserID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				close(client.Send)
			}
			h.mu.Unlock()
			logger.Info("WebSocket client unregistered", "user_id", client.UserID)
		}
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Unregister(client *Client) {
	h.unregister <- client
}

func (h *Hub) SendToUser(userID uint64, msg Message) {
	h.mu.RLock()
	client, ok := h.clients[userID]
	h.mu.RUnlock()

	if !ok {
		logger.Debug("WebSocket SendToUser: user not connected", "user_id", userID)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("WebSocket SendToUser: failed to marshal message", "user_id", userID, "error", err)
		return
	}

	logger.Debug("WebSocket SendToUser: sending message", "user_id", userID, "type", msg.Type)

	select {
	case client.Send <- data:
		logger.Debug("WebSocket SendToUser: message sent", "user_id", userID)
	default:
		logger.Warn("WebSocket SendToUser: channel full, removing client", "user_id", userID)
		h.mu.Lock()
		delete(h.clients, userID)
		close(client.Send)
		h.mu.Unlock()
	}
}

func (c *Client) WritePump() {
	defer func() {
		c.Conn.Close()
	}()

	for message := range c.Send {
		if err := c.Conn.WriteMessage(websocket.TextMessage, message); err != nil {
			return
		}
	}
}

func (c *Client) ReadPump(hub *Hub) {
	defer func() {
		hub.Unregister(c)
		c.Conn.Close()
	}()

	for {
		_, _, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}
		// heartbeat or other messages can be handled here
	}
}
