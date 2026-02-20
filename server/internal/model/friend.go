package model

import "time"

const (
	FriendStatusPending  = 0 // waiting for acceptance
	FriendStatusAccepted = 1 // accepted, now friends
	FriendStatusRejected = 2 // rejected
)

// Friend represents a friend relationship record
type Friend struct {
	ID           uint64    `json:"id"`
	UserID       uint64    `json:"-"`         // applicant internal ID
	UserOpenID   string    `json:"user_id"`   // applicant open ID
	FriendID     uint64    `json:"-"`         // target user internal ID
	FriendOpenID string    `json:"friend_id"` // target user open ID
	Status       int8      `json:"status"`
	Message      string    `json:"message"` // application message
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`

	// join fields
	User   *User `json:"user,omitempty"`   // applicant user info
	Friend *User `json:"friend,omitempty"` // target user info
}

// FriendRequestReq is the request for sending friend request
type FriendRequestReq struct {
	FriendOpenID string `json:"friend_id" binding:"required"` // target user open ID
	Message      string `json:"message" binding:"max=200"`    // optional application message
}

// FriendListResp is the response for friend list
type FriendListResp struct {
	List  []*Friend `json:"list"`
	Total int64     `json:"total"`
}

// FriendRequestListResp is the response for friend request list
type FriendRequestListResp struct {
	List  []*Friend `json:"list"`
	Total int64     `json:"total"`
}

// UserPublicProfile is the public profile visible to others
type UserPublicProfile struct {
	OpenID   string `json:"open_id"`
	Nickname string `json:"nickname"`
	Avatar   string `json:"avatar"`
	Gender   int8   `json:"gender"`
	City     string `json:"city"`
	Province string `json:"province"`

	// friendship info
	IsFriend      bool `json:"is_friend"`       // whether they are friends
	IsPending     bool `json:"is_pending"`      // whether there's a pending request
	IsRequester   bool `json:"is_requester"`    // whether current user is the requester
	FriendshipID  uint64 `json:"friendship_id"` // friend record ID if exists

	// friend-only visible fields
	CarBrand string `json:"car_brand,omitempty"`
	CarModel string `json:"car_model,omitempty"`
	CarColor string `json:"car_color,omitempty"`

	// recent trips (last 3 days)
	RecentTrips []*Trip `json:"recent_trips"`
}

// FriendCountResp is the response for friend and request count
type FriendCountResp struct {
	FriendCount  int64 `json:"friend_count"`
	RequestCount int64 `json:"request_count"` // pending requests received
}
