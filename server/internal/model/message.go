package model

import "time"

// Message represents a private chat message between two users
type Message struct {
	ID         uint64    `json:"id"`
	SenderID   uint64    `json:"-"`              // internal ID
	ReceiverID uint64    `json:"-"`              // internal ID
	SenderOpenID   string `json:"sender_id"`     // open_id for external
	ReceiverOpenID string `json:"receiver_id"`   // open_id for external
	Content    string    `json:"content"`
	MsgType    int8      `json:"msg_type"` // 1: text, 2: image
	IsRead     int8      `json:"is_read"`  // 0: unread, 1: read
	CreatedAt  time.Time `json:"created_at"`
	// joined fields
	Sender   *User `json:"sender,omitempty"`
	Receiver *User `json:"receiver,omitempty"`
}

// MsgType constants
const (
	MsgTypeText  int8 = 1
	MsgTypeImage int8 = 2
)

// MessageSendReq is the request body for sending a message
type MessageSendReq struct {
	ReceiverID string `json:"receiver_id" binding:"required"` // open_id
	Content    string `json:"content" binding:"required"`     // text max 2000, image is COS key
	MsgType    int8   `json:"msg_type" binding:"required,oneof=1 2"`
}

// MessageListReq is the request params for listing messages
type MessageListReq struct {
	PeerID   string `form:"peer_id" binding:"required"` // open_id
	Page     int    `form:"page" binding:"min=1"`
	PageSize int    `form:"page_size" binding:"min=1,max=100"`
}

// MessageListResp is the response for listing messages
type MessageListResp struct {
	List  []*Message `json:"list"`
	Total int64      `json:"total"`
}

// Conversation represents a chat conversation with another user
type Conversation struct {
	PeerID        uint64    `json:"-"`              // internal ID
	PeerOpenID    string    `json:"peer_id"`        // open_id for external
	Peer          *User     `json:"peer"`
	LastMessage   *Message  `json:"last_message"`
	UnreadCount   int       `json:"unread_count"`
	LastMessageAt time.Time `json:"last_message_at"`
}

// ConversationListResp is the response for listing conversations
type ConversationListResp struct {
	List []*Conversation `json:"list"`
}
