package repository

import (
	"database/sql"
	"pinche/internal/database"
	"pinche/internal/model"
)

type MessageRepository struct{}

func NewMessageRepository() *MessageRepository {
	return &MessageRepository{}
}

// Create inserts a new message into the database
func (r *MessageRepository) Create(msg *model.Message) error {
	query := `INSERT INTO messages (sender_id, receiver_id, content, msg_type, is_read) VALUES (?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, msg.SenderID, msg.ReceiverID, msg.Content, msg.MsgType, msg.IsRead)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	msg.ID = uint64(id)
	return nil
}

// GetByID retrieves a message by its ID
func (r *MessageRepository) GetByID(id uint64) (*model.Message, error) {
	query := `SELECT id, sender_id, receiver_id, content, msg_type, is_read, created_at FROM messages WHERE id = ?`
	msg := &model.Message{}
	err := database.DB.QueryRow(query, id).Scan(
		&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.MsgType, &msg.IsRead, &msg.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return msg, nil
}

// GetConversationMessages retrieves messages between two users with pagination
func (r *MessageRepository) GetConversationMessages(userID, peerID uint64, page, pageSize int) ([]*model.Message, int64, error) {
	// count total messages in conversation
	countQuery := `SELECT COUNT(*) FROM messages WHERE (sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)`
	var total int64
	err := database.DB.QueryRow(countQuery, userID, peerID, peerID, userID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// get messages with pagination, ordered by created_at desc
	offset := (page - 1) * pageSize
	query := `
		SELECT m.id, m.sender_id, m.receiver_id, m.content, m.msg_type, m.is_read, m.created_at,
		       s.open_id as sender_open_id, r.open_id as receiver_open_id
		FROM messages m
		LEFT JOIN users s ON s.id = m.sender_id
		LEFT JOIN users r ON r.id = m.receiver_id
		WHERE (m.sender_id = ? AND m.receiver_id = ?) OR (m.sender_id = ? AND m.receiver_id = ?)
		ORDER BY m.created_at DESC
		LIMIT ? OFFSET ?
	`
	rows, err := database.DB.Query(query, userID, peerID, peerID, userID, pageSize, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var messages []*model.Message
	for rows.Next() {
		msg := &model.Message{}
		err := rows.Scan(&msg.ID, &msg.SenderID, &msg.ReceiverID, &msg.Content, &msg.MsgType, &msg.IsRead, &msg.CreatedAt,
			&msg.SenderOpenID, &msg.ReceiverOpenID)
		if err != nil {
			return nil, 0, err
		}
		messages = append(messages, msg)
	}

	return messages, total, nil
}

// GetConversations retrieves all conversations for a user
func (r *MessageRepository) GetConversations(userID uint64) ([]*model.Conversation, error) {
	// use subquery to get latest message for each conversation
	query := `
		SELECT 
			peer_id,
			COALESCE((SELECT open_id FROM users WHERE id = peer_id), '') as peer_open_id,
			COALESCE((SELECT content FROM messages m2 WHERE 
				((m2.sender_id = ? AND m2.receiver_id = peer_id) OR (m2.sender_id = peer_id AND m2.receiver_id = ?))
				ORDER BY m2.created_at DESC LIMIT 1), '') as last_content,
			COALESCE((SELECT msg_type FROM messages m3 WHERE 
				((m3.sender_id = ? AND m3.receiver_id = peer_id) OR (m3.sender_id = peer_id AND m3.receiver_id = ?))
				ORDER BY m3.created_at DESC LIMIT 1), 1) as last_msg_type,
			(SELECT created_at FROM messages m4 WHERE 
				((m4.sender_id = ? AND m4.receiver_id = peer_id) OR (m4.sender_id = peer_id AND m4.receiver_id = ?))
				ORDER BY m4.created_at DESC LIMIT 1) as last_message_at,
			(SELECT COUNT(*) FROM messages m5 WHERE m5.sender_id = peer_id AND m5.receiver_id = ? AND m5.is_read = 0) as unread_count
		FROM (
			SELECT DISTINCT CASE 
				WHEN sender_id = ? THEN receiver_id 
				ELSE sender_id 
			END as peer_id
			FROM messages 
			WHERE sender_id = ? OR receiver_id = ?
		) AS peers
		ORDER BY last_message_at DESC
	`
	rows, err := database.DB.Query(query,
		userID, userID, userID, userID, userID, userID, userID, userID, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var conversations []*model.Conversation
	for rows.Next() {
		conv := &model.Conversation{
			LastMessage: &model.Message{},
		}
		var lastMessageAt sql.NullTime
		err := rows.Scan(
			&conv.PeerID,
			&conv.PeerOpenID,
			&conv.LastMessage.Content,
			&conv.LastMessage.MsgType,
			&lastMessageAt,
			&conv.UnreadCount,
		)
		if err != nil {
			return nil, err
		}
		if lastMessageAt.Valid {
			conv.LastMessageAt = lastMessageAt.Time
		}
		conversations = append(conversations, conv)
	}

	return conversations, nil
}

// MarkAsRead marks all messages from a sender to a receiver as read
func (r *MessageRepository) MarkAsRead(receiverID, senderID uint64) error {
	query := `UPDATE messages SET is_read = 1 WHERE sender_id = ? AND receiver_id = ? AND is_read = 0`
	_, err := database.DB.Exec(query, senderID, receiverID)
	return err
}

// GetUnreadCount returns the count of unread messages for a user
func (r *MessageRepository) GetUnreadCount(userID uint64) (int, error) {
	query := `SELECT COUNT(*) FROM messages WHERE receiver_id = ? AND is_read = 0`
	var count int
	err := database.DB.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
