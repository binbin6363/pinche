package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
	"pinche/internal/websocket"
)

type MessageHandler struct {
	service     *service.MessageService
	userService *service.UserService
	wsHub       *websocket.Hub
}

func NewMessageHandler(service *service.MessageService, userService *service.UserService, wsHub *websocket.Hub) *MessageHandler {
	return &MessageHandler{
		service:     service,
		userService: userService,
		wsHub:       wsHub,
	}
}

// SendMessage handles POST /api/messages
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.MessageSendReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	msg, err := h.service.SendMessage(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	// get receiver internal ID for websocket
	receiver, _ := h.userService.GetByOpenID(req.ReceiverID)
	if receiver != nil {
		// push message to receiver via websocket
		h.wsHub.SendToUser(receiver.ID, websocket.Message{
			Type: "new_message",
			Data: msg,
		})
	}

	c.JSON(http.StatusOK, model.Success(msg))
}

// GetConversationMessages handles GET /api/messages
func (h *MessageHandler) GetConversationMessages(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.MessageListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	resp, err := h.service.GetConversationMessages(userID, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取消息失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// GetConversations handles GET /api/conversations
func (h *MessageHandler) GetConversations(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.service.GetConversations(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取会话列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// MarkAsRead handles PUT /api/messages/read
func (h *MessageHandler) MarkAsRead(c *gin.Context) {
	userID := middleware.GetUserID(c)
	peerOpenID := c.Query("peer_id")
	if peerOpenID == "" {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误"))
		return
	}

	if err := h.service.MarkAsRead(userID, peerOpenID); err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "标记已读失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// GetUnreadCount handles GET /api/messages/unread-count
func (h *MessageHandler) GetUnreadCount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	count, err := h.service.GetUnreadCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取未读数失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(gin.H{"count": count}))
}
