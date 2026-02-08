package model

import "time"

const (
	ConfirmStatusPending  = 0
	ConfirmStatusAccepted = 1
	ConfirmStatusRejected = 2

	MatchStatusPending = 0
	MatchStatusSuccess = 1
	MatchStatusFailed  = 2
)

type Match struct {
	ID                uint64    `json:"id"`
	DriverTripID      uint64    `json:"driver_trip_id"`
	PassengerTripID   uint64    `json:"passenger_trip_id"`
	DriverID          uint64    `json:"-"`           // internal ID
	PassengerID       uint64    `json:"-"`           // internal ID
	DriverOpenID      string    `json:"driver_id"`   // open_id for external
	PassengerOpenID   string    `json:"passenger_id"` // open_id for external
	MatchScore        float64   `json:"match_score"`
	DriverStatus      int8      `json:"driver_status"`
	PassengerStatus   int8      `json:"passenger_status"`
	Status            int8      `json:"status"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`

	// join fields
	DriverTrip    *Trip `json:"driver_trip,omitempty"`
	PassengerTrip *Trip `json:"passenger_trip,omitempty"`
	Driver        *User `json:"driver,omitempty"`
	Passenger     *User `json:"passenger,omitempty"`
}

type MatchConfirmReq struct {
	Accept bool `json:"accept"`
}

type MatchListResp struct {
	List  []*Match `json:"list"`
	Total int64    `json:"total"`
}

type ContactInfo struct {
	DriverPhone       string `json:"driver_phone"`
	DriverNickname    string `json:"driver_nickname"`
	PassengerPhone    string `json:"passenger_phone"`
	PassengerNickname string `json:"passenger_nickname"`
}
