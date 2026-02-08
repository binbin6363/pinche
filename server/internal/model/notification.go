package model

import "time"

type Notification struct {
	ID        uint64    `json:"id"`
	UserID    uint64    `json:"user_id"`
	MatchID   uint64    `json:"match_id"`
	TripID    uint64    `json:"trip_id"`
	Title     string    `json:"title"`
	Content   string    `json:"content"`
	IsRead    int8      `json:"is_read"`
	CreatedAt time.Time `json:"created_at"`
}

type NotificationListResp struct {
	List   []*Notification `json:"list"`
	Total  int64           `json:"total"`
	Unread int64           `json:"unread"`
}
