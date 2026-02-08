package model

import "time"

type Announcement struct {
	ID        uint64    `json:"id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	Type      int8      `json:"type"`       // 1-normal 2-important 3-urgent
	IsActive  int8      `json:"is_active"`  // 0-inactive 1-active
	SortOrder int       `json:"sort_order"` // higher = more priority
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type AnnouncementCreateReq struct {
	Title     string `json:"title" binding:"required,max=100"`
	Content   string `json:"content" binding:"required,max=1000"`
	Type      int8   `json:"type" binding:"oneof=1 2 3"`
	IsActive  int8   `json:"is_active" binding:"oneof=0 1"`
	SortOrder int    `json:"sort_order"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}

type AnnouncementUpdateReq struct {
	Title     string `json:"title" binding:"max=100"`
	Content   string `json:"content" binding:"max=1000"`
	Type      int8   `json:"type"`
	IsActive  int8   `json:"is_active"`
	SortOrder int    `json:"sort_order"`
	StartTime string `json:"start_time"`
	EndTime   string `json:"end_time"`
}
