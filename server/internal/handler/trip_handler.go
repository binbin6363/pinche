package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pinche/internal/logger"
	"pinche/internal/middleware"
	"pinche/internal/model"
	"pinche/internal/service"
)

type TripHandler struct {
	service *service.TripService
}

func NewTripHandler(service *service.TripService) *TripHandler {
	return &TripHandler{service: service}
}

func (h *TripHandler) Create(c *gin.Context) {
	userID := middleware.GetUserID(c)
	var req model.TripCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	trip, err := h.service.Create(userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(trip))
}

func (h *TripHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	// get viewer ID if authenticated
	viewerID := middleware.GetUserID(c)

	trip, err := h.service.GetByIDAndIncrementView(id, viewerID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取行程失败: "+err.Error()))
		return
	}
	if trip == nil {
		c.JSON(http.StatusNotFound, model.Error(model.ErrCodeNotFound, "行程不存在"))
		return
	}

	c.JSON(http.StatusOK, model.Success(trip))
}

func (h *TripHandler) List(c *gin.Context) {
	var req model.TripListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误: "+err.Error()))
		return
	}

	resp, err := h.service.List(&req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取行程列表失败"))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}

func (h *TripHandler) GetMyTrips(c *gin.Context) {
	userID := middleware.GetUserID(c)
	trips, err := h.service.GetMyTrips(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, model.Error(model.ErrCodeInternal, "获取行程列表失败"))
		return
	}
	c.JSON(http.StatusOK, model.Success(trips))
}

// GetMyTripDetail returns trip detail with grabbers list for trip owner
func (h *TripHandler) GetMyTripDetail(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	trip, err := h.service.GetMyTripDetail(id, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(trip))
}

// UpdateTrip handles trip update
func (h *TripHandler) UpdateTrip(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	var req model.TripUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "参数错误"))
		return
	}

	needsReview, message, err := h.service.UpdateTrip(id, userID, &req)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(map[string]interface{}{
		"needs_review": needsReview,
		"message":      message,
	}))
}

func (h *TripHandler) Cancel(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	if err := h.service.Cancel(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *TripHandler) Complete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	if err := h.service.Complete(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *TripHandler) Delete(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	if err := h.service.Delete(id, userID); err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(nil))
}

// Admin handlers

func (h *TripHandler) AdminListTrips(c *gin.Context) {
	var req model.AdminTripListReq
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "参数错误"))
		return
	}

	resp, err := h.service.AdminListTrips(&req)
	if err != nil {
		logger.Error("Admin list trips failed", "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "获取行程列表失败"))
		return
	}

	logger.Debug("Admin listed trips", "page", req.Page, "page_size", req.PageSize, "total", resp.Total)
	c.JSON(http.StatusOK, model.Success(resp))
}

func (h *TripHandler) AdminBanTrip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	if err := h.service.AdminBanTrip(id); err != nil {
		logger.Error("Admin ban trip failed", "trip_id", id, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "封禁行程失败"))
		return
	}

	logger.Info("Admin banned trip", "trip_id", id)
	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *TripHandler) AdminUnbanTrip(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	if err := h.service.AdminUnbanTrip(id); err != nil {
		logger.Error("Admin unban trip failed", "trip_id", id, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "解封行程失败"))
		return
	}

	logger.Info("Admin unbanned trip", "trip_id", id)
	c.JSON(http.StatusOK, model.Success(nil))
}

func (h *TripHandler) GrabTrip(c *gin.Context) {
	userID := middleware.GetUserID(c)
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, "无效的行程ID"))
		return
	}

	var req model.GrabTripReq
	// message is optional, so ignore bind errors
	c.ShouldBindJSON(&req)

	resp, err := h.service.GrabTrip(id, userID, req.Message)
	if err != nil {
		c.JSON(http.StatusBadRequest, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	c.JSON(http.StatusOK, model.Success(resp))
}
