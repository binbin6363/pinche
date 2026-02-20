package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
)

type FriendHandler struct {
	service *service.FriendService
}

func NewFriendHandler() *FriendHandler {
	return &FriendHandler{
		service: service.NewFriendService(),
	}
}

// SendFriendRequest handles POST /api/friends/request
func (h *FriendHandler) SendFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	var req model.FriendRequestReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	if err := h.service.SendFriendRequest(userID, &req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// GetFriendRequests handles GET /api/friends/requests
func (h *FriendHandler) GetFriendRequests(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.service.GetFriendRequests(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取好友申请失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// AcceptFriendRequest handles POST /api/friends/requests/:id/accept
func (h *FriendHandler) AcceptFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的申请ID"))
		return
	}

	if err := h.service.AcceptFriendRequest(userID, requestID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// RejectFriendRequest handles POST /api/friends/requests/:id/reject
func (h *FriendHandler) RejectFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的申请ID"))
		return
	}

	if err := h.service.RejectFriendRequest(userID, requestID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// CancelFriendRequest handles DELETE /api/friends/requests/:id
func (h *FriendHandler) CancelFriendRequest(c *gin.Context) {
	userID := middleware.GetUserID(c)

	requestID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的申请ID"))
		return
	}

	if err := h.service.CancelFriendRequest(userID, requestID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// GetFriends handles GET /api/friends
func (h *FriendHandler) GetFriends(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.service.GetFriends(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取好友列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

// DeleteFriend handles DELETE /api/friends/:id
func (h *FriendHandler) DeleteFriend(c *gin.Context) {
	userID := middleware.GetUserID(c)
	friendOpenID := c.Param("id")

	if friendOpenID == "" {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的好友ID"))
		return
	}

	if err := h.service.DeleteFriend(userID, friendOpenID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// GetUserProfile handles GET /api/users/:id/profile
func (h *FriendHandler) GetUserProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	targetOpenID := c.Param("id")

	if targetOpenID == "" {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的用户ID"))
		return
	}

	profile, err := h.service.GetUserPublicProfile(userID, targetOpenID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(profile))
}

// GetFriendCount handles GET /api/friends/count
func (h *FriendHandler) GetFriendCount(c *gin.Context) {
	userID := middleware.GetUserID(c)

	resp, err := h.service.GetFriendCount(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取好友数量失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
