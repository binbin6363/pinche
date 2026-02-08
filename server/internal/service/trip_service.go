package service

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"pinche/internal/cache"
	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/repository"
	"pinche/internal/websocket"
)

type TripService struct {
	repo         *repository.TripRepository
	userRepo     *repository.UserRepository
	notifyRepo   *repository.NotificationRepository
	matchService *MatchService
	wsHub        *websocket.Hub
	tripCache    *cache.TripCache
}

func NewTripService(matchService *MatchService, wsHub *websocket.Hub) *TripService {
	return &TripService{
		repo:         repository.NewTripRepository(),
		userRepo:     repository.NewUserRepository(),
		notifyRepo:   repository.NewNotificationRepository(),
		matchService: matchService,
		wsHub:        wsHub,
		tripCache:    cache.NewTripCache(),
	}
}

func (s *TripService) Create(userID uint64, req *model.TripCreateReq) (*model.Trip, error) {
	// check active trips limit (max 2)
	activeCount, err := s.repo.CountActiveByUserID(userID)
	if err != nil {
		logger.Error("Count active trips failed", "user_id", userID, "error", err)
		return nil, errors.New("查询行程数量失败")
	}
	if activeCount >= 2 {
		logger.Warn("User exceeded active trips limit", "user_id", userID, "active_count", activeCount)
		return nil, errors.New("您最多只能同时拥有2个有效行程，请先取消或完成现有行程")
	}

	// check daily publish limit (max 5)
	todayCount, err := s.repo.CountTodayByUserID(userID)
	if err != nil {
		logger.Error("Count today trips failed", "user_id", userID, "error", err)
		return nil, errors.New("查询今日发布数量失败")
	}
	if todayCount >= 5 {
		logger.Warn("User exceeded daily publish limit", "user_id", userID, "today_count", todayCount)
		return nil, errors.New("您今日发布行程次数已达上限（5次），请明日再试")
	}

	departureTime, err := time.ParseInLocation("2006-01-02 15:04", req.DepartureTime, time.Local)
	if err != nil {
		return nil, errors.New("出发时间格式错误")
	}

	if departureTime.Before(time.Now()) {
		return nil, errors.New("出发时间不能早于当前时间")
	}

	trip := &model.Trip{
		UserID:              userID,
		TripType:            req.TripType,
		DepartureCity:       req.DepartureCity,
		DepartureProvince:   req.DepartureProvince,
		DepartureAddress:    req.DepartureAddress,
		DepartureLat:        req.DepartureLat,
		DepartureLng:        req.DepartureLng,
		DestinationCity:     req.DestinationCity,
		DestinationProvince: req.DestinationProvince,
		DestinationAddress:  req.DestinationAddress,
		DestinationLat:      req.DestinationLat,
		DestinationLng:      req.DestinationLng,
		DepartureTime:       departureTime,
		Seats:               req.Seats,
		Price:               req.Price,
		Remark:              req.Remark,
		Images:              req.Images,
		Status:              model.TripStatusPending,
	}

	if err := s.repo.Create(trip); err != nil {
		logger.Error("Create trip failed", "user_id", userID, "error", err)
		return nil, err
	}

	logger.Info("Trip created",
		"trip_id", trip.ID,
		"user_id", userID,
		"trip_type", req.TripType,
		"from", req.DepartureCity,
		"to", req.DestinationCity)

	// invalidate trip list cache
	go s.tripCache.InvalidateTripLists()

	// async match
	go s.matchService.FindAndNotifyMatches(trip)

	return trip, nil
}

func (s *TripService) GetByID(id uint64) (*model.Trip, error) {
	// Cache Aside: try cache first
	if trip, err := s.tripCache.GetTrip(id); err == nil && trip != nil {
		return trip, nil
	}

	// cache miss, get from DB
	trip, err := s.repo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if trip == nil {
		return nil, nil
	}

	// store in cache
	go s.tripCache.SetTrip(trip)
	return trip, nil
}

// GetByIDAndIncrementView gets trip by ID and increments view count (for non-owner views)
func (s *TripService) GetByIDAndIncrementView(id uint64, viewerID uint64) (*model.Trip, error) {
	// Cache Aside: try cache first
	trip, err := s.tripCache.GetTrip(id)
	if err != nil || trip == nil {
		// cache miss, get from DB
		trip, err = s.repo.GetByID(id)
		if err != nil {
			return nil, err
		}
		if trip == nil {
			return nil, nil
		}
		// store in cache
		go s.tripCache.SetTrip(trip)
	}

	// only increment view count if viewer is not the owner
	if trip.UserID != viewerID {
		go s.repo.IncrementViewCount(id)
	}
	return trip, nil
}

// GetMyTripDetail gets trip detail with grabbers list (only for trip owner)
func (s *TripService) GetMyTripDetail(tripID uint64, userID uint64) (*model.Trip, error) {
	// Cache Aside: try cache first for basic trip info
	trip, err := s.tripCache.GetTrip(tripID)
	if err != nil || trip == nil {
		trip, err = s.repo.GetByID(tripID)
		if err != nil {
			return nil, err
		}
		if trip == nil {
			return nil, errors.New("行程不存在")
		}
		// store in cache
		go s.tripCache.SetTrip(trip)
	}

	if trip.UserID != userID {
		return nil, errors.New("无权查看此行程")
	}

	// get grabbers list (not cached as it changes frequently)
	grabbers, err := s.repo.GetGrabsByTripID(tripID)
	if err == nil {
		trip.Grabbers = grabbers
	}

	return trip, nil
}

func (s *TripService) List(req *model.TripListReq) (*model.TripListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	// Cache Aside: try cache first
	if cached, err := s.tripCache.GetTripList(req); err == nil && cached != nil {
		return &model.TripListResp{
			List:  cached.List,
			Total: cached.Total,
		}, nil
	}

	// cache miss, get from DB
	trips, total, err := s.repo.List(req)
	if err != nil {
		return nil, err
	}

	// store in cache
	go s.tripCache.SetTripList(req, trips, total)

	return &model.TripListResp{
		List:  trips,
		Total: total,
	}, nil
}

func (s *TripService) GetMyTrips(userID uint64) ([]*model.Trip, error) {
	return s.repo.GetByUserID(userID)
}

func (s *TripService) Cancel(id uint64, userID uint64) error {
	trip, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if trip == nil {
		return errors.New("行程不存在")
	}
	if trip.UserID != userID {
		return errors.New("无权操作此行程")
	}
	if trip.Status != model.TripStatusPending {
		return errors.New("只能取消待匹配的行程")
	}
	if err := s.repo.UpdateStatus(id, model.TripStatusCancelled); err != nil {
		return err
	}
	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(id)
		s.tripCache.InvalidateTripLists()
	}()
	return nil
}

func (s *TripService) Complete(id uint64, userID uint64) error {
	trip, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if trip == nil {
		return errors.New("行程不存在")
	}
	if trip.UserID != userID {
		return errors.New("无权操作此行程")
	}
	if trip.Status != model.TripStatusPending && trip.Status != model.TripStatusMatched {
		return errors.New("只能标记待匹配或已匹配的行程为已成行")
	}
	if err := s.repo.UpdateStatus(id, model.TripStatusCompleted); err != nil {
		return err
	}
	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(id)
		s.tripCache.InvalidateTripLists()
	}()
	return nil
}

func (s *TripService) Delete(id uint64, userID uint64) error {
	trip, err := s.repo.GetByID(id)
	if err != nil {
		return err
	}
	if trip == nil {
		return errors.New("行程不存在")
	}
	if trip.UserID != userID {
		return errors.New("无权操作此行程")
	}
	if err := s.repo.Delete(id, userID); err != nil {
		return err
	}
	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(id)
		s.tripCache.InvalidateTripLists()
	}()
	return nil
}

// Admin functions

func (s *TripService) AdminListTrips(req *model.AdminTripListReq) (*model.TripListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	trips, total, err := s.repo.AdminListAll(req)
	if err != nil {
		logger.Error("Admin list trips failed", "error", err)
		return nil, err
	}

	return &model.TripListResp{
		List:  trips,
		Total: total,
	}, nil
}

func (s *TripService) AdminBanTrip(id uint64) error {
	logger.Info("Admin banning trip", "trip_id", id)
	if err := s.repo.UpdateStatus(id, model.TripStatusBanned); err != nil {
		return err
	}
	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(id)
		s.tripCache.InvalidateTripLists()
	}()
	return nil
}

func (s *TripService) AdminUnbanTrip(id uint64) error {
	logger.Info("Admin unbanning trip", "trip_id", id)
	if err := s.repo.UpdateStatus(id, model.TripStatusPending); err != nil {
		return err
	}
	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(id)
		s.tripCache.InvalidateTripLists()
	}()
	return nil
}

// GrabTrip handles when a user wants to grab/accept a trip
// For driver trips (trip_type=1): passenger grabs the trip
// For passenger trips (trip_type=2): driver grabs the trip
func (s *TripService) GrabTrip(tripID uint64, grabberID uint64, message string) (*model.GrabTripResp, error) {
	// get trip info
	trip, err := s.repo.GetByID(tripID)
	if err != nil {
		logger.Error("Get trip failed for grab", "trip_id", tripID, "error", err)
		return nil, errors.New("获取行程信息失败")
	}
	if trip == nil {
		return nil, errors.New("行程不存在")
	}

	// cannot grab own trip
	if trip.UserID == grabberID {
		return nil, errors.New("不能抢自己的行程")
	}

	// only pending trips can be grabbed
	if trip.Status != model.TripStatusPending {
		return nil, errors.New("该行程已不可抢单")
	}

	// get grabber info
	grabber, err := s.userRepo.GetByID(grabberID)
	if err != nil || grabber == nil {
		logger.Error("Get grabber info failed", "grabber_id", grabberID, "error", err)
		return nil, errors.New("获取用户信息失败")
	}

	// save grab record
	grab := &model.TripGrab{
		TripID:  tripID,
		UserID:  grabberID,
		Message: message,
	}
	if err := s.repo.CreateGrab(grab); err != nil {
		// ignore duplicate errors, continue to notify
		logger.Debug("Create grab record failed (may be duplicate)", "trip_id", tripID, "grabber_id", grabberID, "error", err)
	}

	logger.Info("Trip grabbed",
		"trip_id", tripID,
		"grabber_id", grabberID,
		"owner_id", trip.UserID,
		"trip_type", trip.TripType)

	// build notification content
	var title, content string
	if trip.TripType == model.TripTypeDriver {
		// driver's trip, grabber is a passenger
		title = "有乘客想搭您的车"
		content = fmt.Sprintf("乘客「%s」想搭您的车（%s→%s），请查看详情并联系对方",
			s.maskNickname(grabber.Nickname),
			trip.DepartureCity,
			trip.DestinationCity)
		if message != "" {
			content += fmt.Sprintf("。留言：%s", message)
		}
	} else {
		// passenger's trip, grabber is a driver
		title = "有司机愿意接单"
		content = fmt.Sprintf("司机「%s」愿意接您的行程（%s→%s），请查看详情并联系对方",
			s.maskNickname(grabber.Nickname),
			trip.DepartureCity,
			trip.DestinationCity)
		if message != "" {
			content += fmt.Sprintf("。留言：%s", message)
		}
	}

	// create notification for trip owner
	notification := &model.Notification{
		UserID:  trip.UserID,
		TripID:  tripID,
		Title:   title,
		Content: content,
	}
	if err := s.notifyRepo.Create(notification); err != nil {
		logger.Error("Create grab notification failed", "trip_id", tripID, "owner_id", trip.UserID, "error", err)
		return nil, errors.New("发送通知失败")
	}

	// send real-time notification via websocket
	if s.wsHub != nil {
		s.wsHub.SendToUser(trip.UserID, websocket.Message{
			Type: "trip_grabbed",
			Data: map[string]interface{}{
				"trip_id":      tripID,
				"grabber_id":   grabber.OpenID,
				"grabber_name": s.maskNickname(grabber.Nickname),
				"notification": notification,
			},
		})
	}

	return &model.GrabTripResp{
		Success: true,
		Message: "抢单成功，已通知车主",
	}, nil
}

func (s *TripService) maskNickname(nickname string) string {
	if nickname == "" {
		return "用户**"
	}
	if len(nickname) <= 2 {
		return nickname[:1] + "**"
	}
	runes := []rune(nickname)
	if len(runes) <= 2 {
		return string(runes[:1]) + "**"
	}
	return string(runes[:len(runes)-2]) + "**"
}

// UpdateTrip updates a trip
// Direct updates: images, remark, seats, price
// Review required: departure/destination city/address, departure_time
func (s *TripService) UpdateTrip(tripID uint64, userID uint64, req *model.TripUpdateReq) (bool, string, error) {
	trip, err := s.repo.GetByID(tripID)
	if err != nil {
		return false, "", errors.New("获取行程失败")
	}
	if trip == nil {
		return false, "", errors.New("行程不存在")
	}
	if trip.UserID != userID {
		return false, "", errors.New("无权修改此行程")
	}
	if trip.Status != model.TripStatusPending && trip.Status != model.TripStatusMatched {
		return false, "", errors.New("只能修改待匹配或已匹配的行程")
	}

	needsReview := false
	reviewMessage := ""

	// check if location change requires review
	locationChanged := (req.DepartureCity != "" && req.DepartureCity != trip.DepartureCity) ||
		(req.DepartureAddress != "" && req.DepartureAddress != trip.DepartureAddress) ||
		(req.DestinationCity != "" && req.DestinationCity != trip.DestinationCity) ||
		(req.DestinationAddress != "" && req.DestinationAddress != trip.DestinationAddress)

	if locationChanged {
		needsReview = true
		oldValue := map[string]string{
			"departure_city":      trip.DepartureCity,
			"departure_address":   trip.DepartureAddress,
			"destination_city":    trip.DestinationCity,
			"destination_address": trip.DestinationAddress,
		}
		newValue := map[string]string{
			"departure_city":      ifEmpty(req.DepartureCity, trip.DepartureCity),
			"departure_address":   ifEmpty(req.DepartureAddress, trip.DepartureAddress),
			"destination_city":    ifEmpty(req.DestinationCity, trip.DestinationCity),
			"destination_address": ifEmpty(req.DestinationAddress, trip.DestinationAddress),
		}
		oldJSON, _ := json.Marshal(oldValue)
		newJSON, _ := json.Marshal(newValue)

		update := &model.TripUpdate{
			TripID:     tripID,
			UserID:     userID,
			UpdateType: model.TripUpdateTypeLocation,
			OldValue:   string(oldJSON),
			NewValue:   string(newJSON),
			Status:     model.TripUpdateStatusPending,
		}
		if err := s.repo.CreateTripUpdate(update); err != nil {
			return false, "", errors.New("提交修改审核失败")
		}
		reviewMessage += "起终点修改已提交审核；"
	}

	// check if time change requires review
	if req.DepartureTime != "" {
		newTime, err := time.ParseInLocation("2006-01-02 15:04", req.DepartureTime, time.Local)
		if err == nil && !newTime.Equal(trip.DepartureTime) {
			needsReview = true
			oldValue := map[string]string{"departure_time": trip.DepartureTime.Format("2006-01-02 15:04")}
			newValue := map[string]string{"departure_time": req.DepartureTime}
			oldJSON, _ := json.Marshal(oldValue)
			newJSON, _ := json.Marshal(newValue)

			update := &model.TripUpdate{
				TripID:     tripID,
				UserID:     userID,
				UpdateType: model.TripUpdateTypeTime,
				OldValue:   string(oldJSON),
				NewValue:   string(newJSON),
				Status:     model.TripUpdateStatusPending,
			}
			if err := s.repo.CreateTripUpdate(update); err != nil {
				return false, "", errors.New("提交时间修改审核失败")
			}
			reviewMessage += "出发时间修改已提交审核；"
		}
	}

	// direct updates (images, remark, seats, price)
	images := req.Images
	if images == "" {
		images = trip.Images
	}
	remark := req.Remark
	if remark == "" {
		remark = trip.Remark
	}

	if err := s.repo.UpdateTrip(tripID, userID, images, remark, req.Seats, req.Price); err != nil {
		return false, "", errors.New("更新行程失败")
	}

	// invalidate cache
	go func() {
		s.tripCache.InvalidateTrip(tripID)
		s.tripCache.InvalidateTripLists()
	}()

	if needsReview {
		return true, reviewMessage, nil
	}
	return false, "行程更新成功", nil
}

func ifEmpty(val, defaultVal string) string {
	if val == "" {
		return defaultVal
	}
	return val
}
