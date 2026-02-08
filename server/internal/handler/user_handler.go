package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"pinche/internal/logger"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

func (h *UserHandler) Register(c *gin.Context) {
	var req model.UserRegisterReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	user, err := h.service.Register(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(user))
}

func (h *UserHandler) Login(c *gin.Context) {
	var req model.UserLoginReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	resp, err := h.service.Login(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	user, err := h.service.GetByID(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取用户信息失败"))
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, model.Error(model.ErrCodeNotFound, "用户不存在"))
		return
	}
	c.JSON(http.StatusOK, model.Success(user))
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.UserUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	user, err := h.service.Update(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(user))
}

// Admin handlers

func (h *UserHandler) AdminListUsers(c *gin.Context) {
	var req model.AdminUserListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "参数错误"))
		return
	}

	resp, err := h.service.AdminListUsers(&req)
	if err != nil {
		logger.Error("Admin list users failed", "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "获取用户列表失败"))
		return
	}

	logger.Debug("Admin listed users", "page", req.Page, "page_size", req.PageSize, "total", resp.Total)
	c.JSON(http.StatusOK, model.Success(resp))
}

func (h *UserHandler) AdminBanUser(c *gin.Context) {
	openID := c.Param("id") // now using open_id
	if openID == "" {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "无效的用户ID"))
		return
	}

	if err := h.service.AdminBanUser(openID); err != nil {
		logger.Error("Admin ban user failed", "open_id", openID, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "封禁用户失败"))
		return
	}

	logger.Info("Admin banned user", "open_id", openID)
	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *UserHandler) AdminUnbanUser(c *gin.Context) {
	openID := c.Param("id") // now using open_id
	if openID == "" {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "无效的用户ID"))
		return
	}

	if err := h.service.AdminUnbanUser(openID); err != nil {
		logger.Error("Admin unban user failed", "open_id", openID, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "解封用户失败"))
		return
	}

	logger.Info("Admin unbanned user", "open_id", openID)
	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *UserHandler) AdminGetStats(c *gin.Context) {
	stats, err := h.service.AdminGetStats()
	if err != nil {
		logger.Error("Admin get stats failed", "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "获取统计数据失败"))
		return
	}

	logger.Debug("Admin retrieved stats")
	c.JSON(http.StatusOK, model.Success(stats))
}
