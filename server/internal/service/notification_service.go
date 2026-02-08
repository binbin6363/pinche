package service

import (
	"pinche/internal/model"
	"pinche/internal/repository"
)

type NotificationService struct {
	repo *repository.NotificationRepository
}

func NewNotificationService() *NotificationService {
	return &NotificationService{
		repo: repository.NewNotificationRepository(),
	}
}

func (s *NotificationService) GetList(userID uint64, page, pageSize int) (*model.NotificationListResp, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	list, total, unread, err := s.repo.GetByUserID(userID, page, pageSize)
	if err != nil {
		return nil, err
	}

	return &model.NotificationListResp{
		List:   list,
		Total:  total,
		Unread: unread,
	}, nil
}

func (s *NotificationService) MarkAsRead(id, userID uint64) error {
	return s.repo.MarkAsRead(id, userID)
}

func (s *NotificationService) MarkAllAsRead(userID uint64) error {
	return s.repo.MarkAllAsRead(userID)
}
