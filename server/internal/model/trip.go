package model

import "time"

const (
	TripTypeDriver    = 1
	TripTypePassenger = 2

	TripStatusPending   = 1
	TripStatusMatched   = 2
	TripStatusCompleted = 3
	TripStatusCancelled = 4
	TripStatusBanned    = 5 // admin banned

	// trip update types
	TripUpdateTypeLocation = 1 // departure/destination change
	TripUpdateTypeTime     = 2 // departure time change

	// trip update status
	TripUpdateStatusPending  = 0
	TripUpdateStatusApproved = 1
	TripUpdateStatusRejected = 2
)

type Trip struct {
	ID                  uint64    `json:"id"`
	UserID              uint64    `json:"user_internal_id"` // internal ID, needed for cache
	UserOpenID          string    `json:"user_id"`          // open_id for external
	TripType            int8      `json:"trip_type"`
	DepartureCity       string    `json:"departure_city"`
	DepartureProvince   string    `json:"departure_province"`
	DepartureAddress    string    `json:"departure_address"`
	DepartureLat        float64   `json:"departure_lat"`
	DepartureLng        float64   `json:"departure_lng"`
	DestinationCity     string    `json:"destination_city"`
	DestinationProvince string    `json:"destination_province"`
	DestinationAddress  string    `json:"destination_address"`
	DestinationLat      float64   `json:"destination_lat"`
	DestinationLng      float64   `json:"destination_lng"`
	DepartureTime       time.Time `json:"departure_time"`
	Seats               int       `json:"seats"`
	Price               float64   `json:"price"`
	Remark              string    `json:"remark"`
	Images              string    `json:"images"` // JSON array of image URLs
	Status              int8      `json:"status"`
	ViewCount           int       `json:"view_count"`
	CreatedAt           time.Time `json:"created_at"`
	UpdatedAt           time.Time `json:"updated_at"`

	// join fields
	User    *User       `json:"user,omitempty"`
	Grabbers []*TripGrab `json:"grabbers,omitempty"` // users who grabbed this trip
}

type TripCreateReq struct {
	TripType            int8    `json:"trip_type" binding:"required,oneof=1 2"`
	DepartureCity       string  `json:"departure_city" binding:"required"`
	DepartureProvince   string  `json:"departure_province"`
	DepartureAddress    string  `json:"departure_address" binding:"required"`
	DepartureLat        float64 `json:"departure_lat"`
	DepartureLng        float64 `json:"departure_lng"`
	DestinationCity     string  `json:"destination_city" binding:"required"`
	DestinationProvince string  `json:"destination_province"`
	DestinationAddress  string  `json:"destination_address" binding:"required"`
	DestinationLat      float64 `json:"destination_lat"`
	DestinationLng      float64 `json:"destination_lng"`
	DepartureTime       string  `json:"departure_time" binding:"required"`
	Seats               int     `json:"seats" binding:"required,min=1,max=7"`
	Price               float64 `json:"price" binding:"min=0"`
	Remark              string  `json:"remark" binding:"max=500"`
	Images              string  `json:"images"` // JSON array of image URLs
}

type TripListReq struct {
	TripType        int8   `form:"trip_type"`
	DepartureCity   string `form:"departure_city"`
	DestinationCity string `form:"destination_city"`
	Date            string `form:"date"`
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=20"`
	// for relevance filtering
	UserCity     string `form:"user_city"`
	UserProvince string `form:"user_province"`
	// exclude current user's trips (open_id from frontend)
	ExcludeUserOpenID string `form:"exclude_user_id"`
	// internal user_id (set by handler after lookup)
	ExcludeUserID uint64 `form:"-"`
}

type TripListResp struct {
	List  []*Trip `json:"list"`
	Total int64   `json:"total"`
}

type AdminTripListReq struct {
	Search          string `form:"search"`
	TripType        int8   `form:"trip_type"`
	DepartureCity   string `form:"departure_city"`
	DestinationCity string `form:"destination_city"`
	Status          *int8  `form:"status"`
	UserOpenID      string `form:"user_id"` // open_id
	Page            int    `form:"page,default=1"`
	PageSize        int    `form:"page_size,default=20"`
}

type GrabTripReq struct {
	Message string `json:"message"` // optional message from grabber
}

type GrabTripResp struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// TripGrab represents a user who grabbed a trip
type TripGrab struct {
	ID        uint64    `json:"id"`
	TripID    uint64    `json:"trip_id"`
	UserID    uint64    `json:"-"`
	Message   string    `json:"message"`
	CreatedAt time.Time `json:"created_at"`
	User      *User     `json:"user,omitempty"`
}

// TripUpdate represents a pending update request for a trip
type TripUpdate struct {
	ID           uint64    `json:"id"`
	TripID       uint64    `json:"trip_id"`
	UserID       uint64    `json:"-"`
	UpdateType   int8      `json:"update_type"`
	OldValue     string    `json:"old_value"`
	NewValue     string    `json:"new_value"`
	Status       int8      `json:"status"`
	RejectReason string    `json:"reject_reason"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// TripUpdateReq is the request for updating a trip
type TripUpdateReq struct {
	Images             string  `json:"images"`              // can update directly
	Remark             string  `json:"remark"`              // can update directly
	Seats              *int    `json:"seats"`               // can update directly
	Price              *float64 `json:"price"`              // can update directly
	DepartureCity      string  `json:"departure_city"`      // requires review
	DepartureAddress   string  `json:"departure_address"`   // requires review
	DestinationCity    string  `json:"destination_city"`    // requires review
	DestinationAddress string  `json:"destination_address"` // requires review
	DepartureTime      string  `json:"departure_time"`      // requires review
}
