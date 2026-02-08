package service

import (
	"errors"
	"time"

	"pinche/internal/model"
	"pinche/internal/repository"
)

type MessageService struct {
	repo     *repository.MessageRepository
	userRepo *repository.UserRepository
}

func NewMessageService() *MessageService {
	return &MessageService{
		repo:     repository.NewMessageRepository(),
		userRepo: repository.NewUserRepository(),
	}
}

// SendMessage creates a new message from sender to receiver
func (s *MessageService) SendMessage(senderID uint64, req *model.MessageSendReq) (*model.Message, error) {
	// validate content length for text messages
	if req.MsgType == model.MsgTypeText && len(req.Content) > 2000 {
		return nil, errors.New("文本消息内容不能超过2000字符")
	}

	// get receiver by open_id
	receiver, err := s.userRepo.GetByOpenID(req.ReceiverID)
	if err != nil {
		return nil, err
	}
	if receiver == nil {
		return nil, errors.New("接收者不存在")
	}

	// cannot send message to self
	if senderID == receiver.ID {
		return nil, errors.New("不能给自己发送消息")
	}

	// get sender info
	sender, err := s.userRepo.GetByID(senderID)
	if err != nil {
		return nil, err
	}

	msg := &model.Message{
		SenderID:       senderID,
		ReceiverID:     receiver.ID,
		SenderOpenID:   sender.OpenID,
		ReceiverOpenID: receiver.OpenID,
		Content:        req.Content,
		MsgType:        req.MsgType,
		IsRead:         0,
		CreatedAt:      time.Now(),
	}

	if err := s.repo.Create(msg); err != nil {
		return nil, err
	}

	// attach user info
	msg.Sender = sender
	msg.Receiver = receiver

	return msg, nil
}

// GetConversationMessages retrieves messages between current user and peer
func (s *MessageService) GetConversationMessages(userID uint64, req *model.MessageListReq) (*model.MessageListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// get peer by open_id
	peer, err := s.userRepo.GetByOpenID(req.PeerID)
	if err != nil {
		return nil, err
	}
	if peer == nil {
		return nil, errors.New("用户不存在")
	}

	messages, total, err := s.repo.GetConversationMessages(userID, peer.ID, req.Page, req.PageSize)
	if err != nil {
		return nil, err
	}

	// attach user info to messages
	userCache := make(map[uint64]*model.User)
	for _, msg := range messages {
		if _, ok := userCache[msg.SenderID]; !ok {
			user, _ := s.userRepo.GetByID(msg.SenderID)
			userCache[msg.SenderID] = user
		}
		msg.Sender = userCache[msg.SenderID]
	}

	return &model.MessageListResp{
		List:  messages,
		Total: total,
	}, nil
}

// GetConversations retrieves all conversations for a user
func (s *MessageService) GetConversations(userID uint64) (*model.ConversationListResp, error) {
	conversations, err := s.repo.GetConversations(userID)
	if err != nil {
		return nil, err
	}

	// attach peer user info
	for _, conv := range conversations {
		peer, _ := s.userRepo.GetByID(conv.PeerID)
		conv.Peer = peer
	}

	return &model.ConversationListResp{
		List: conversations,
	}, nil
}

// MarkAsRead marks all messages from peer as read
func (s *MessageService) MarkAsRead(userID uint64, peerOpenID string) error {
	peer, err := s.userRepo.GetByOpenID(peerOpenID)
	if err != nil {
		return err
	}
	if peer == nil {
		return errors.New("用户不存在")
	}
	return s.repo.MarkAsRead(userID, peer.ID)
}

// GetUnreadCount returns total unread message count for user
func (s *MessageService) GetUnreadCount(userID uint64) (int, error) {
	return s.repo.GetUnreadCount(userID)
}

// GetMessageByID retrieves a message by ID with permission check
func (s *MessageService) GetMessageByID(userID, msgID uint64) (*model.Message, error) {
	msg, err := s.repo.GetByID(msgID)
	if err != nil {
		return nil, err
	}
	if msg == nil {
		return nil, errors.New("消息不存在")
	}

	// only sender or receiver can view the message
	if msg.SenderID != userID && msg.ReceiverID != userID {
		return nil, errors.New("无权查看该消息")
	}

	return msg, nil
}
