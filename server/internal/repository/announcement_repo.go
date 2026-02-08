package repository

import (
	"database/sql"
	"time"

	"pinche/internal/database"
	"pinche/internal/model"
)

type AnnouncementRepository struct{}

func NewAnnouncementRepository() *AnnouncementRepository {
	return &AnnouncementRepository{}
}

func (r *AnnouncementRepository) Create(ann *model.Announcement) error {
	query := `INSERT INTO announcements (title, content, type, is_active, sort_order, start_time, end_time) 
		VALUES (?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, ann.Title, ann.Content, ann.Type, ann.IsActive, ann.SortOrder, ann.StartTime, ann.EndTime)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	ann.ID = uint64(id)
	return nil
}

func (r *AnnouncementRepository) GetByID(id uint64) (*model.Announcement, error) {
	query := `SELECT id, title, content, type, is_active, sort_order, start_time, end_time, created_at, updated_at 
		FROM announcements WHERE id = ?`
	ann := &model.Announcement{}
	err := database.DB.QueryRow(query, id).Scan(
		&ann.ID, &ann.Title, &ann.Content, &ann.Type, &ann.IsActive, &ann.SortOrder, &ann.StartTime, &ann.EndTime, &ann.CreatedAt, &ann.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return ann, nil
}

// GetActiveAnnouncements returns currently active announcements for public display
func (r *AnnouncementRepository) GetActiveAnnouncements(limit int) ([]*model.Announcement, error) {
	now := time.Now()
	query := `SELECT id, title, content, type, is_active, sort_order, start_time, end_time, created_at, updated_at 
		FROM announcements 
		WHERE is_active = 1 AND start_time <= ? AND end_time >= ?
		ORDER BY sort_order DESC, created_at DESC
		LIMIT ?`

	rows, err := database.DB.Query(query, now, now, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var announcements []*model.Announcement
	for rows.Next() {
		ann := &model.Announcement{}
		err := rows.Scan(&ann.ID, &ann.Title, &ann.Content, &ann.Type, &ann.IsActive, &ann.SortOrder, &ann.StartTime, &ann.EndTime, &ann.CreatedAt, &ann.UpdatedAt)
		if err != nil {
			return nil, err
		}
		announcements = append(announcements, ann)
	}
	return announcements, nil
}

// ListAll returns all announcements for admin panel
func (r *AnnouncementRepository) ListAll(page, pageSize int) ([]*model.Announcement, int64, error) {
	// count
	countQuery := `SELECT COUNT(*) FROM announcements`
	var total int64
	if err := database.DB.QueryRow(countQuery).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list
	offset := (page - 1) * pageSize
	query := `SELECT id, title, content, type, is_active, sort_order, start_time, end_time, created_at, updated_at 
		FROM announcements ORDER BY created_at DESC LIMIT ? OFFSET ?`

	rows, err := database.DB.Query(query, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var announcements []*model.Announcement
	for rows.Next() {
		ann := &model.Announcement{}
		err := rows.Scan(&ann.ID, &ann.Title, &ann.Content, &ann.Type, &ann.IsActive, &ann.SortOrder, &ann.StartTime, &ann.EndTime, &ann.CreatedAt, &ann.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		announcements = append(announcements, ann)
	}
	return announcements, total, nil
}

func (r *AnnouncementRepository) Update(ann *model.Announcement) error {
	query := `UPDATE announcements SET title = ?, content = ?, type = ?, is_active = ?, sort_order = ?, start_time = ?, end_time = ? WHERE id = ?`
	_, err := database.DB.Exec(query, ann.Title, ann.Content, ann.Type, ann.IsActive, ann.SortOrder, ann.StartTime, ann.EndTime, ann.ID)
	return err
}

func (r *AnnouncementRepository) Delete(id uint64) error {
	query := `DELETE FROM announcements WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}
