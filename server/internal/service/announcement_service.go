package service

import (
	"strings"
	"time"

	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/repository"
)

type AnnouncementService struct {
	repo *repository.AnnouncementRepository
}

func NewAnnouncementService() *AnnouncementService {
	return &AnnouncementService{
		repo: repository.NewAnnouncementRepository(),
	}
}

// parseTimeString parses time string with multiple format support
func parseTimeString(s string) time.Time {
	if s == "" {
		return time.Time{}
	}
	// try ISO 8601 format (from frontend toISOString())
	if t, err := time.Parse(time.RFC3339, s); err == nil {
		return t
	}
	// try ISO format without timezone
	if t, err := time.Parse("2006-01-02T15:04:05.000Z", s); err == nil {
		return t
	}
	// try standard datetime format
	if t, err := time.Parse("2006-01-02 15:04:05", s); err == nil {
		return t
	}
	// try datetime-local input format (from HTML input)
	if t, err := time.Parse("2006-01-02T15:04", s); err == nil {
		return t
	}
	// try date only
	if t, err := time.Parse("2006-01-02", strings.Split(s, "T")[0]); err == nil {
		return t
	}
	return time.Time{}
}

func (s *AnnouncementService) Create(req *model.AnnouncementCreateReq) (*model.Announcement, error) {
	startTime := parseTimeString(req.StartTime)
	endTime := parseTimeString(req.EndTime)

	// default time range if not provided
	if startTime.IsZero() {
		startTime = time.Now()
	}
	if endTime.IsZero() {
		endTime = startTime.AddDate(1, 0, 0) // 1 year later
	}

	ann := &model.Announcement{
		Title:     req.Title,
		Content:   req.Content,
		Type:      req.Type,
		IsActive:  req.IsActive,
		SortOrder: req.SortOrder,
		StartTime: startTime,
		EndTime:   endTime,
	}

	if ann.Type == 0 {
		ann.Type = 1
	}
	if req.IsActive == 0 {
		ann.IsActive = 1
	}

	if err := s.repo.Create(ann); err != nil {
		logger.Error("Create announcement failed", "title", req.Title, "error", err)
		return nil, err
	}
	logger.Info("Announcement created", "id", ann.ID, "title", ann.Title)
	return ann, nil
}

func (s *AnnouncementService) GetByID(id uint64) (*model.Announcement, error) {
	return s.repo.GetByID(id)
}

func (s *AnnouncementService) GetActiveAnnouncements(limit int) ([]*model.Announcement, error) {
	return s.repo.GetActiveAnnouncements(limit)
}

func (s *AnnouncementService) ListAll(page, pageSize int) ([]*model.Announcement, int64, error) {
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}
	return s.repo.ListAll(page, pageSize)
}

func (s *AnnouncementService) Update(id uint64, req *model.AnnouncementUpdateReq) (*model.Announcement, error) {
	ann, err := s.repo.GetByID(id)
	if err != nil || ann == nil {
		return nil, err
	}

	if req.Title != "" {
		ann.Title = req.Title
	}
	if req.Content != "" {
		ann.Content = req.Content
	}
	if req.Type > 0 {
		ann.Type = req.Type
	}
	ann.IsActive = req.IsActive
	if req.SortOrder > 0 {
		ann.SortOrder = req.SortOrder
	}
	if req.StartTime != "" {
		startTime := parseTimeString(req.StartTime)
		if !startTime.IsZero() {
			ann.StartTime = startTime
		}
	}
	if req.EndTime != "" {
		endTime := parseTimeString(req.EndTime)
		if !endTime.IsZero() {
			ann.EndTime = endTime
		}
	}

	if err := s.repo.Update(ann); err != nil {
		logger.Error("Update announcement failed", "id", id, "error", err)
		return nil, err
	}
	logger.Info("Announcement updated", "id", id, "title", ann.Title)
	return ann, nil
}

func (s *AnnouncementService) Delete(id uint64) error {
	if err := s.repo.Delete(id); err != nil {
		logger.Error("Delete announcement failed", "id", id, "error", err)
		return err
	}
	logger.Info("Announcement deleted", "id", id)
	return nil
}
