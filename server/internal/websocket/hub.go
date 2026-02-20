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

// SignalingMessage represents call signaling messages from client
type SignalingMessage struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

// SignalingData contains common fields for call signaling
type SignalingData struct {
	TargetOpenID string      `json:"target_open_id"`
	CallID       string      `json:"call_id"`
	CallType     string      `json:"call_type,omitempty"`
	CallerInfo   interface{} `json:"caller_info,omitempty"`
	Accept       bool        `json:"accept,omitempty"`
	Reason       string      `json:"reason,omitempty"`
	SDP          interface{} `json:"sdp,omitempty"`
	Candidate    interface{} `json:"candidate,omitempty"`
}

type Client struct {
	UserID uint64
	OpenID string
	Conn   *websocket.Conn
	Send   chan []byte
}

type Hub struct {
	clients       map[uint64]*Client
	clientsByOpen map[string]*Client // open_id -> client mapping
	register      chan *Client
	unregister    chan *Client
	mu            sync.RWMutex
}

func NewHub() *Hub {
	return &Hub{
		clients:       make(map[uint64]*Client),
		clientsByOpen: make(map[string]*Client),
		register:      make(chan *Client),
		unregister:    make(chan *Client),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			h.clients[client.UserID] = client
			if client.OpenID != "" {
				h.clientsByOpen[client.OpenID] = client
			}
			h.mu.Unlock()
			logger.Info("WebSocket client registered", "user_id", client.UserID, "open_id", client.OpenID)

		case client := <-h.unregister:
			h.mu.Lock()
			if _, ok := h.clients[client.UserID]; ok {
				delete(h.clients, client.UserID)
				if client.OpenID != "" {
					delete(h.clientsByOpen, client.OpenID)
				}
				close(client.Send)
			}
			h.mu.Unlock()
			logger.Info("WebSocket client unregistered", "user_id", client.UserID, "open_id", client.OpenID)
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
		if client.OpenID != "" {
			delete(h.clientsByOpen, client.OpenID)
		}
		close(client.Send)
		h.mu.Unlock()
	}
}

// SendToUserByOpenID sends a message to a user by their open_id
func (h *Hub) SendToUserByOpenID(openID string, msg Message) {
	h.mu.RLock()
	client, ok := h.clientsByOpen[openID]
	h.mu.RUnlock()

	if !ok {
		logger.Debug("WebSocket SendToUserByOpenID: user not connected", "open_id", openID)
		return
	}

	data, err := json.Marshal(msg)
	if err != nil {
		logger.Error("WebSocket SendToUserByOpenID: failed to marshal message", "open_id", openID, "error", err)
		return
	}

	logger.Debug("WebSocket SendToUserByOpenID: sending message", "open_id", openID, "type", msg.Type)

	select {
	case client.Send <- data:
		logger.Debug("WebSocket SendToUserByOpenID: message sent", "open_id", openID)
	default:
		logger.Warn("WebSocket SendToUserByOpenID: channel full, removing client", "open_id", openID)
		h.mu.Lock()
		delete(h.clients, client.UserID)
		delete(h.clientsByOpen, openID)
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
		_, message, err := c.Conn.ReadMessage()
		if err != nil {
			break
		}

		// Parse message to check if it's a call signaling message
		var sigMsg SignalingMessage
		if err := json.Unmarshal(message, &sigMsg); err != nil {
			logger.Debug("WebSocket ReadPump: failed to parse message", "error", err)
			continue
		}

		// Handle call signaling messages
		if isCallSignaling(sigMsg.Type) {
			hub.handleCallSignaling(c, sigMsg)
		}
	}
}

// isCallSignaling checks if the message type is a call signaling type
func isCallSignaling(msgType string) bool {
	switch msgType {
	case "call_invite", "call_answer", "call_end", "webrtc_offer", "webrtc_answer", "ice_candidate":
		return true
	default:
		return false
	}
}

// handleCallSignaling processes and forwards call signaling messages
func (h *Hub) handleCallSignaling(senderClient *Client, msg SignalingMessage) {
	logger.Info("WebSocket: received call signaling",
		"type", msg.Type,
		"sender_id", senderClient.UserID,
		"sender_open_id", senderClient.OpenID,
		"raw_data", string(msg.Data),
	)

	var data SignalingData
	if err := json.Unmarshal(msg.Data, &data); err != nil {
		logger.Error("WebSocket: failed to parse signaling data",
			"error", err,
			"sender_id", senderClient.UserID,
			"raw_data", string(msg.Data),
		)
		return
	}

	targetOpenID := data.TargetOpenID
	if targetOpenID == "" {
		logger.Warn("WebSocket: signaling message missing target_open_id",
			"sender_id", senderClient.UserID,
			"type", msg.Type,
			"parsed_data", data,
		)
		return
	}

	logger.Info("WebSocket: forwarding call signaling",
		"type", msg.Type,
		"from_open_id", senderClient.OpenID,
		"to_open_id", targetOpenID,
		"call_id", data.CallID,
	)

	// Build message to forward, include sender info
	forwardData := map[string]interface{}{
		"call_id":        data.CallID,
		"from_user_id":   senderClient.UserID,
		"from_open_id":   senderClient.OpenID,
	}

	// Copy relevant fields based on message type
	switch msg.Type {
	case "call_invite":
		forwardData["call_type"] = data.CallType
		forwardData["caller_info"] = data.CallerInfo
	case "call_answer":
		forwardData["accept"] = data.Accept
	case "call_end":
		forwardData["reason"] = data.Reason
	case "webrtc_offer", "webrtc_answer":
		forwardData["sdp"] = data.SDP
	case "ice_candidate":
		forwardData["candidate"] = data.Candidate
	}

	// Forward to target user by open_id
	h.SendToUserByOpenID(targetOpenID, Message{
		Type: msg.Type,
		Data: forwardData,
	})
}
