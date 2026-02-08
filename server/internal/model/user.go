package model

import "time"

type User struct {
	ID        uint64    `json:"-"`                 // internal ID, not exposed
	OpenID    string    `json:"open_id"`           // public ID for external use
	Phone     string    `json:"phone"`
	Password  string    `json:"-"`
	Nickname  string    `json:"nickname"`
	Avatar    string    `json:"avatar"`
	Gender    int8      `json:"gender"`
	Status    int8      `json:"status"`    // 0-normal 1-banned
	City      string    `json:"city"`      // user's city for trip matching
	Province  string    `json:"province"`  // user's province for trip matching
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type UserRegisterReq struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required,min=6"`
	Nickname string `json:"nickname" binding:"required,min=2,max=20"`
}

type UserLoginReq struct {
	Phone    string `json:"phone" binding:"required,len=11"`
	Password string `json:"password" binding:"required"`
}

type UserLoginResp struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type UserUpdateReq struct {
	Nickname string `json:"nickname" binding:"omitempty,min=2,max=20"`
	Avatar   string `json:"avatar"`
	Gender   *int8  `json:"gender" binding:"omitempty,min=0,max=2"`
	City     string `json:"city"`
	Province string `json:"province"`
}

type AdminUserListReq struct {
	Search   string `form:"search"`
	Phone    string `form:"phone"`
	Nickname string `form:"nickname"`
	Status   *int8  `form:"status"`
	Page     int    `form:"page,default=1"`
	PageSize int    `form:"page_size,default=20"`
}

type AdminUserListResp struct {
	List  []*User `json:"list"`
	Total int64   `json:"total"`
}
