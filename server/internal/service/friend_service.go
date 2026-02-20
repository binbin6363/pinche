package service

import (
	"errors"

	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/repository"
)

type FriendService struct {
	friendRepo *repository.FriendRepository
	userRepo   *repository.UserRepository
}

func NewFriendService() *FriendService {
	return &FriendService{
		friendRepo: repository.NewFriendRepository(),
		userRepo:   repository.NewUserRepository(),
	}
}

// SendFriendRequest sends a friend request from current user to target user
func (s *FriendService) SendFriendRequest(userID uint64, req *model.FriendRequestReq) error {
	// get target user by open_id
	targetUser, err := s.userRepo.GetByOpenID(req.FriendOpenID)
	if err != nil {
		logger.Error("Get target user failed", "open_id", req.FriendOpenID, "error", err)
		return err
	}
	if targetUser == nil {
		return errors.New("目标用户不存在")
	}

	// cannot add self
	if targetUser.ID == userID {
		return errors.New("不能添加自己为好友")
	}

	// check existing friendship record
	existing, err := s.friendRepo.GetFriendshipRecord(userID, targetUser.ID)
	if err != nil {
		logger.Error("Check existing friendship failed", "user_id", userID, "target_id", targetUser.ID, "error", err)
		return err
	}

	if existing != nil {
		switch existing.Status {
		case model.FriendStatusAccepted:
			return errors.New("你们已经是好友了")
		case model.FriendStatusPending:
			if existing.UserID == userID {
				return errors.New("已发送过好友申请，请等待对方同意")
			}
			// the other user sent request, we can accept it directly
			return errors.New("对方已向你发送好友申请，请在好友申请列表中同意")
		case model.FriendStatusRejected:
			// if previously rejected, allow re-sending (reset the record)
			if existing.UserID == userID {
				// user sent before and was rejected, can re-send
				if err := s.friendRepo.ResetRejectedRequest(existing.ID, req.Message); err != nil {
					logger.Error("Reset rejected request failed", "id", existing.ID, "error", err)
					return err
				}
				logger.Info("Friend request re-sent", "from", userID, "to", targetUser.ID)
				return nil
			}
			// user rejected the other's request before, need to delete and create new
			if err := s.friendRepo.Delete(existing.ID); err != nil {
				logger.Error("Delete rejected record failed", "id", existing.ID, "error", err)
				return err
			}
		}
	}

	// create new friend request
	friend := &model.Friend{
		UserID:   userID,
		FriendID: targetUser.ID,
		Status:   model.FriendStatusPending,
		Message:  req.Message,
	}

	if err := s.friendRepo.Create(friend); err != nil {
		logger.Error("Create friend request failed", "from", userID, "to", targetUser.ID, "error", err)
		return err
	}

	logger.Info("Friend request sent", "from", userID, "to", targetUser.ID)
	return nil
}

// AcceptFriendRequest accepts a friend request
func (s *FriendService) AcceptFriendRequest(userID uint64, requestID uint64) error {
	friend, err := s.friendRepo.GetByID(requestID)
	if err != nil {
		logger.Error("Get friend request failed", "id", requestID, "error", err)
		return err
	}
	if friend == nil {
		return errors.New("好友申请不存在")
	}

	// only the target user can accept
	if friend.FriendID != userID {
		return errors.New("无权操作此申请")
	}

	if friend.Status != model.FriendStatusPending {
		return errors.New("该申请已处理")
	}

	if err := s.friendRepo.UpdateStatus(requestID, model.FriendStatusAccepted); err != nil {
		logger.Error("Accept friend request failed", "id", requestID, "error", err)
		return err
	}

	logger.Info("Friend request accepted", "id", requestID, "user", userID, "friend", friend.UserID)
	return nil
}

// RejectFriendRequest rejects a friend request
func (s *FriendService) RejectFriendRequest(userID uint64, requestID uint64) error {
	friend, err := s.friendRepo.GetByID(requestID)
	if err != nil {
		logger.Error("Get friend request failed", "id", requestID, "error", err)
		return err
	}
	if friend == nil {
		return errors.New("好友申请不存在")
	}

	// only the target user can reject
	if friend.FriendID != userID {
		return errors.New("无权操作此申请")
	}

	if friend.Status != model.FriendStatusPending {
		return errors.New("该申请已处理")
	}

	if err := s.friendRepo.UpdateStatus(requestID, model.FriendStatusRejected); err != nil {
		logger.Error("Reject friend request failed", "id", requestID, "error", err)
		return err
	}

	logger.Info("Friend request rejected", "id", requestID, "user", userID, "friend", friend.UserID)
	return nil
}

// GetFriendRequests gets pending friend requests for current user
func (s *FriendService) GetFriendRequests(userID uint64) (*model.FriendRequestListResp, error) {
	friends, total, err := s.friendRepo.GetPendingRequests(userID)
	if err != nil {
		logger.Error("Get friend requests failed", "user_id", userID, "error", err)
		return nil, err
	}

	return &model.FriendRequestListResp{
		List:  friends,
		Total: total,
	}, nil
}

// GetFriends gets all friends for current user
func (s *FriendService) GetFriends(userID uint64) (*model.FriendListResp, error) {
	friends, total, err := s.friendRepo.GetFriends(userID)
	if err != nil {
		logger.Error("Get friends failed", "user_id", userID, "error", err)
		return nil, err
	}

	return &model.FriendListResp{
		List:  friends,
		Total: total,
	}, nil
}

// DeleteFriend deletes a friend relationship
func (s *FriendService) DeleteFriend(userID uint64, friendOpenID string) error {
	// get target user
	targetUser, err := s.userRepo.GetByOpenID(friendOpenID)
	if err != nil {
		logger.Error("Get target user failed", "open_id", friendOpenID, "error", err)
		return err
	}
	if targetUser == nil {
		return errors.New("用户不存在")
	}

	// check friendship exists
	isFriend, err := s.friendRepo.CheckFriendship(userID, targetUser.ID)
	if err != nil {
		logger.Error("Check friendship failed", "user_id", userID, "target_id", targetUser.ID, "error", err)
		return err
	}
	if !isFriend {
		return errors.New("你们不是好友")
	}

	// delete both directions
	if err := s.friendRepo.DeleteByUserAndFriend(userID, targetUser.ID); err != nil {
		logger.Error("Delete friend failed", "user_id", userID, "target_id", targetUser.ID, "error", err)
		return err
	}

	logger.Info("Friend deleted", "user", userID, "friend", targetUser.ID)
	return nil
}

// GetUserPublicProfile gets the public profile of a user
func (s *FriendService) GetUserPublicProfile(currentUserID uint64, targetOpenID string) (*model.UserPublicProfile, error) {
	// get target user
	targetUser, err := s.userRepo.GetByOpenID(targetOpenID)
	if err != nil {
		logger.Error("Get target user failed", "open_id", targetOpenID, "error", err)
		return nil, err
	}
	if targetUser == nil {
		return nil, errors.New("用户不存在")
	}

	profile := &model.UserPublicProfile{
		OpenID:   targetUser.OpenID,
		Nickname: targetUser.Nickname,
		Avatar:   targetUser.Avatar,
		Gender:   targetUser.Gender,
		City:     targetUser.City,
		Province: targetUser.Province,
	}

	// check friendship status
	if currentUserID != targetUser.ID {
		friendship, err := s.friendRepo.GetFriendshipRecord(currentUserID, targetUser.ID)
		if err != nil {
			logger.Error("Get friendship record failed", "user_id", currentUserID, "target_id", targetUser.ID, "error", err)
			return nil, err
		}

		if friendship != nil {
			profile.FriendshipID = friendship.ID
			profile.IsRequester = friendship.UserID == currentUserID

			switch friendship.Status {
			case model.FriendStatusAccepted:
				profile.IsFriend = true
			case model.FriendStatusPending:
				profile.IsPending = true
			}
		}

		// if friends, show car info (except car number)
		if profile.IsFriend {
			profile.CarBrand = targetUser.CarBrand
			profile.CarModel = targetUser.CarModel
			profile.CarColor = targetUser.CarColor
		}
	} else {
		// viewing own profile - show all car info
		profile.CarBrand = targetUser.CarBrand
		profile.CarModel = targetUser.CarModel
		profile.CarColor = targetUser.CarColor
	}

	// get recent trips (last 3 days)
	trips, err := s.friendRepo.GetRecentTripsByUserID(targetUser.ID, 3)
	if err != nil {
		logger.Error("Get recent trips failed", "user_id", targetUser.ID, "error", err)
		return nil, err
	}
	profile.RecentTrips = trips

	return profile, nil
}

// GetFriendCount gets friend and request counts for current user
func (s *FriendService) GetFriendCount(userID uint64) (*model.FriendCountResp, error) {
	friendCount, err := s.friendRepo.GetFriendCount(userID)
	if err != nil {
		logger.Error("Get friend count failed", "user_id", userID, "error", err)
		return nil, err
	}

	requestCount, err := s.friendRepo.GetPendingRequestCount(userID)
	if err != nil {
		logger.Error("Get request count failed", "user_id", userID, "error", err)
		return nil, err
	}

	return &model.FriendCountResp{
		FriendCount:  friendCount,
		RequestCount: requestCount,
	}, nil
}

// CancelFriendRequest cancels a sent friend request
func (s *FriendService) CancelFriendRequest(userID uint64, requestID uint64) error {
	friend, err := s.friendRepo.GetByID(requestID)
	if err != nil {
		logger.Error("Get friend request failed", "id", requestID, "error", err)
		return err
	}
	if friend == nil {
		return errors.New("好友申请不存在")
	}

	// only the sender can cancel
	if friend.UserID != userID {
		return errors.New("无权操作此申请")
	}

	if friend.Status != model.FriendStatusPending {
		return errors.New("该申请已处理")
	}

	if err := s.friendRepo.Delete(requestID); err != nil {
		logger.Error("Cancel friend request failed", "id", requestID, "error", err)
		return err
	}

	logger.Info("Friend request cancelled", "id", requestID, "user", userID)
	return nil
}
