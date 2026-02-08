package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
)

type MatchHandler struct {
	service *service.MatchService
}

func NewMatchHandler(service *service.MatchService) *MatchHandler {
	return &MatchHandler{service: service}
}

func (h *MatchHandler) GetMyMatches(c *gin.Context) {
	userID := middleware.GetUserID(c)
	matches, err := h.service.GetMyMatches(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取匹配列表失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(matches))
}

func (h *MatchHandler) Confirm(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的匹配ID"))
		return
	}

	var req model.MatchConfirmReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误"))
		return
	}

	if err := h.service.Confirm(id, userID, req.Accept); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *MatchHandler) GetContactInfo(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的匹配ID"))
		return
	}

	info, err := h.service.GetContactInfo(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(info))
}

func (h *MatchHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的匹配ID"))
		return
	}

	match, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取匹配信息失败"))
		return
	}
	if match == nil {
		c.JSON(http.StatusNotFound, model.Error(model.ErrCodeNotFound, "匹配记录不存在"))
		return
	}

	c.JSON(http.StatusOK, model.Success(match))
}
