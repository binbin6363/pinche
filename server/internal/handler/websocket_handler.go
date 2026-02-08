package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"pinche/internal/logger"
	"pinche/internal/service"
	ws "pinche/internal/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WebSocketHandler struct {
	hub         *ws.Hub
	userService *service.UserService
}

func NewWebSocketHandler(hub *ws.Hub, userService *service.UserService) *WebSocketHandler {
	return &WebSocketHandler{
		hub:         hub,
		userService: userService,
	}
}

func (h *WebSocketHandler) HandleConnection(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		logger.Warn("WebSocket connection failed: missing token", "client_ip", c.ClientIP())
		c.JSON(http.StatusUnauthorized, gin.H{"error": "missing token"})
		return
	}

	userID, err := h.userService.ParseToken(token)
	if err != nil {
		logger.Warn("WebSocket connection failed: invalid token",
			"client_ip", c.ClientIP(),
			"token", logger.MaskToken(token),
			"error", err)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		logger.Error("WebSocket upgrade failed", "user_id", userID, "error", err)
		return
	}

	logger.Info("WebSocket connection established", "user_id", userID, "client_ip", c.ClientIP())

	client := &ws.Client{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan []byte, 256),
	}

	h.hub.Register(client)

	go client.WritePump()
	go client.ReadPump(h.hub)
}
