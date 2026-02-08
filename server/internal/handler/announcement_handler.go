package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/service"
)

type AnnouncementHandler struct {
	service *service.AnnouncementService
}

func NewAnnouncementHandler(s *service.AnnouncementService) *AnnouncementHandler {
	return &AnnouncementHandler{service: s}
}

// GetActiveAnnouncements returns currently active announcements for public display
func (h *AnnouncementHandler) GetActiveAnnouncements(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "3"))
	if limit < 1 || limit > 10 {
		limit = 5
	}

	announcements, err := h.service.GetActiveAnnouncements(limit)
	if err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to get announcements"))
		return
	}
	if announcements == nil {
		announcements = []*model.Announcement{}
	}
	c.JSON(http.StatusOK, model.Success(announcements))
}

// ListAll returns all announcements for admin panel
func (h *AnnouncementHandler) ListAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	list, total, err := h.service.ListAll(page, pageSize)
	if err != nil {
		logger.Error("Admin list announcements failed", "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to get announcements"))
		return
	}
	if list == nil {
		list = []*model.Announcement{}
	}
	logger.Debug("Admin listed announcements", "page", page, "page_size", pageSize, "total", total)
	c.JSON(http.StatusOK, model.Success(gin.H{
		"list":  list,
		"total": total,
	}))
}

// Create creates a new announcement
func (h *AnnouncementHandler) Create(c *gin.Context) {
	var req model.AnnouncementCreateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	ann, err := h.service.Create(&req)
	if err != nil {
		logger.Error("Admin create announcement failed", "title", req.Title, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to create announcement"))
		return
	}
	logger.Info("Admin created announcement", "id", ann.ID, "title", ann.Title)
	c.JSON(http.StatusOK, model.Success(ann))
}

// Update updates an announcement
func (h *AnnouncementHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "Invalid ID"))
		return
	}

	var req model.AnnouncementUpdateReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, err.Error()))
		return
	}

	ann, err := h.service.Update(id, &req)
	if err != nil {
		logger.Error("Admin update announcement failed", "id", id, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to update announcement"))
		return
	}
	if ann == nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeNotFound, "Announcement not found"))
		return
	}
	logger.Info("Admin updated announcement", "id", id, "title", ann.Title)
	c.JSON(http.StatusOK, model.Success(ann))
}

// Delete deletes an announcement
func (h *AnnouncementHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		c.JSON(http.StatusOK, model.Error(model.ErrCodeBadRequest, "Invalid ID"))
		return
	}

	if err := h.service.Delete(id); err != nil {
		logger.Error("Admin delete announcement failed", "id", id, "error", err)
		c.JSON(http.StatusOK, model.Error(model.ErrCodeInternal, "Failed to delete announcement"))
		return
	}
	logger.Info("Admin deleted announcement", "id", id)
	c.JSON(http.StatusOK, model.Success(nil))
}
