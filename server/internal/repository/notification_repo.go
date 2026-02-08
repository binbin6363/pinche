package repository

import (
	"pinche/internal/database"
	"pinche/internal/model"
)

type NotificationRepository struct{}

func NewNotificationRepository() *NotificationRepository {
	return &NotificationRepository{}
}

func (r *NotificationRepository) Create(n *model.Notification) error {
	query := `INSERT INTO notifications (user_id, match_id, trip_id, title, content) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, n.UserID, n.MatchID, n.TripID, n.Title, n.Content)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	n.ID = uint64(id)
	return nil
}

func (r *NotificationRepository) GetByUserID(userID uint64, page, pageSize int) ([]*model.Notification, int64, int64, error) {
	// count total
	var total int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ?", userID).Scan(&total); err != nil {
		return nil, 0, 0, err
	}

	// count unread
	var unread int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM notifications WHERE user_id = ? AND is_read = 0", userID).Scan(&unread); err != nil {
		return nil, 0, 0, err
	}

	// list
	offset := (page - 1) * pageSize
	query := `SELECT id, user_id, match_id, COALESCE(trip_id, 0), title, content, is_read, created_at FROM notifications WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`
	rows, err := database.DB.Query(query, userID, pageSize, offset)
	if err != nil {
		return nil, 0, 0, err
	}
	defer rows.Close()

	var notifications []*model.Notification
	for rows.Next() {
		n := &model.Notification{}
		if err := rows.Scan(&n.ID, &n.UserID, &n.MatchID, &n.TripID, &n.Title, &n.Content, &n.IsRead, &n.CreatedAt); err != nil {
			return nil, 0, 0, err
		}
		notifications = append(notifications, n)
	}
	return notifications, total, unread, nil
}

func (r *NotificationRepository) MarkAsRead(id, userID uint64) error {
	query := `UPDATE notifications SET is_read = 1 WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, id, userID)
	return err
}

func (r *NotificationRepository) MarkAllAsRead(userID uint64) error {
	query := `UPDATE notifications SET is_read = 1 WHERE user_id = ?`
	_, err := database.DB.Exec(query, userID)
	return err
}
