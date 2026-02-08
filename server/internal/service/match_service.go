package service

import (
	"errors"
	"fmt"
	"math"

	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/repository"
	"pinche/internal/websocket"
)

type MatchService struct {
	repo         *repository.MatchRepository
	tripRepo     *repository.TripRepository
	userRepo     *repository.UserRepository
	notifyRepo   *repository.NotificationRepository
	wsHub        *websocket.Hub
}

func NewMatchService(wsHub *websocket.Hub) *MatchService {
	return &MatchService{
		repo:       repository.NewMatchRepository(),
		tripRepo:   repository.NewTripRepository(),
		userRepo:   repository.NewUserRepository(),
		notifyRepo: repository.NewNotificationRepository(),
		wsHub:      wsHub,
	}
}

// FindAndNotifyMatches finds matching trips and sends notifications
func (s *MatchService) FindAndNotifyMatches(trip *model.Trip) {
	matchingTrips, err := s.tripRepo.FindMatchingTrips(trip)
	if err != nil {
		logger.Error("Find matching trips failed", "trip_id", trip.ID, "error", err)
		return
	}

	logger.Debug("Finding matches for trip", "trip_id", trip.ID, "candidates", len(matchingTrips))

	for _, matchTrip := range matchingTrips {
		score := s.calculateMatchScore(trip, matchTrip)
		if score < 50 {
			continue
		}

		var driverTripID, passengerTripID uint64
		var driverID, passengerID uint64
		if trip.TripType == model.TripTypeDriver {
			driverTripID = trip.ID
			driverID = trip.UserID
			passengerTripID = matchTrip.ID
			passengerID = matchTrip.UserID
		} else {
			driverTripID = matchTrip.ID
			driverID = matchTrip.UserID
			passengerTripID = trip.ID
			passengerID = trip.UserID
		}

		// check if match already exists
		existing, _ := s.repo.GetByTrips(driverTripID, passengerTripID)
		if existing != nil {
			continue
		}

		match := &model.Match{
			DriverTripID:    driverTripID,
			PassengerTripID: passengerTripID,
			DriverID:        driverID,
			PassengerID:     passengerID,
			MatchScore:      score,
			DriverStatus:    model.ConfirmStatusPending,
			PassengerStatus: model.ConfirmStatusPending,
			Status:          model.MatchStatusPending,
		}

		if err := s.repo.Create(match); err != nil {
			logger.Error("Create match failed",
				"driver_trip_id", driverTripID,
				"passenger_trip_id", passengerTripID,
				"error", err)
			continue
		}

		logger.Info("Match created",
			"match_id", match.ID,
			"driver_id", driverID,
			"passenger_id", passengerID,
			"score", score)

		// send notifications
		s.sendMatchNotification(match, trip, matchTrip)
	}
}

func (s *MatchService) calculateMatchScore(trip1, trip2 *model.Trip) float64 {
	var score float64 = 100

	// time difference (max 12 hours = 720 minutes)
	timeDiff := math.Abs(trip1.DepartureTime.Sub(trip2.DepartureTime).Minutes())
	timeScore := math.Max(0, 100-timeDiff/7.2)
	score = score * (timeScore / 100)

	// location distance (Haversine formula)
	departDist := haversineDistance(
		trip1.DepartureLat, trip1.DepartureLng,
		trip2.DepartureLat, trip2.DepartureLng,
	)
	destDist := haversineDistance(
		trip1.DestinationLat, trip1.DestinationLng,
		trip2.DestinationLat, trip2.DestinationLng,
	)

	// distance penalty (max 50km)
	departScore := math.Max(0, 100-departDist*2)
	destScore := math.Max(0, 100-destDist*2)
	score = score * ((departScore + destScore) / 200)

	return math.Round(score*100) / 100
}

func haversineDistance(lat1, lng1, lat2, lng2 float64) float64 {
	const R = 6371 // Earth radius in km
	dLat := (lat2 - lat1) * math.Pi / 180
	dLng := (lng2 - lng1) * math.Pi / 180
	a := math.Sin(dLat/2)*math.Sin(dLat/2) +
		math.Cos(lat1*math.Pi/180)*math.Cos(lat2*math.Pi/180)*
			math.Sin(dLng/2)*math.Sin(dLng/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))
	return R * c
}

func (s *MatchService) sendMatchNotification(match *model.Match, trip1, trip2 *model.Trip) {
	// notify driver
	driverNotify := &model.Notification{
		UserID:  match.DriverID,
		MatchID: match.ID,
		Title:   "发现匹配乘客",
		Content: fmt.Sprintf("有乘客想从%s到%s，匹配度%.0f%%，请确认是否接受", trip2.DepartureCity, trip2.DestinationCity, match.MatchScore),
	}
	if err := s.notifyRepo.Create(driverNotify); err == nil {
		s.wsHub.SendToUser(match.DriverID, websocket.Message{
			Type: "match_found",
			Data: map[string]interface{}{
				"match_id":     match.ID,
				"notification": driverNotify,
			},
		})
	}

	// notify passenger
	passengerNotify := &model.Notification{
		UserID:  match.PassengerID,
		MatchID: match.ID,
		Title:   "发现匹配司机",
		Content: fmt.Sprintf("有司机从%s到%s，匹配度%.0f%%，请确认是否接受", trip1.DepartureCity, trip1.DestinationCity, match.MatchScore),
	}
	if err := s.notifyRepo.Create(passengerNotify); err == nil {
		s.wsHub.SendToUser(match.PassengerID, websocket.Message{
			Type: "match_found",
			Data: map[string]interface{}{
				"match_id":     match.ID,
				"notification": passengerNotify,
			},
		})
	}
}

func (s *MatchService) GetMyMatches(userID uint64) ([]*model.Match, error) {
	return s.repo.GetByUserID(userID)
}

func (s *MatchService) Confirm(matchID uint64, userID uint64, accept bool) error {
	match, err := s.repo.GetByID(matchID)
	if err != nil {
		return err
	}
	if match == nil {
		return errors.New("匹配记录不存在")
	}
	if match.Status != model.MatchStatusPending {
		return errors.New("匹配已结束")
	}

	isDriver := match.DriverID == userID
	isPassenger := match.PassengerID == userID
	if !isDriver && !isPassenger {
		return errors.New("无权操作此匹配")
	}

	status := model.ConfirmStatusAccepted
	if !accept {
		status = model.ConfirmStatusRejected
	}

	if isDriver {
		if match.DriverStatus != model.ConfirmStatusPending {
			return errors.New("您已确认过")
		}
		if err := s.repo.UpdateDriverStatus(matchID, int8(status)); err != nil {
			return err
		}
		match.DriverStatus = int8(status)
	} else {
		if match.PassengerStatus != model.ConfirmStatusPending {
			return errors.New("您已确认过")
		}
		if err := s.repo.UpdatePassengerStatus(matchID, int8(status)); err != nil {
			return err
		}
		match.PassengerStatus = int8(status)
	}

	// check if both confirmed
	s.checkMatchComplete(match)

	return nil
}

func (s *MatchService) checkMatchComplete(match *model.Match) {
	// if either rejected, match failed
	if match.DriverStatus == model.ConfirmStatusRejected || match.PassengerStatus == model.ConfirmStatusRejected {
		s.repo.UpdateStatus(match.ID, model.MatchStatusFailed)

		// notify the other party
		var notifyUserID uint64
		if match.DriverStatus == model.ConfirmStatusRejected {
			notifyUserID = match.PassengerID
		} else {
			notifyUserID = match.DriverID
		}

		notify := &model.Notification{
			UserID:  notifyUserID,
			MatchID: match.ID,
			Title:   "匹配未成功",
			Content: "对方已拒绝本次匹配，您可以继续寻找其他匹配",
		}
		if err := s.notifyRepo.Create(notify); err == nil {
			s.wsHub.SendToUser(notifyUserID, websocket.Message{
				Type: "match_rejected",
				Data: map[string]interface{}{
					"match_id":     match.ID,
					"notification": notify,
				},
			})
		}
		return
	}

	// if both accepted, match success
	if match.DriverStatus == model.ConfirmStatusAccepted && match.PassengerStatus == model.ConfirmStatusAccepted {
		s.repo.UpdateStatus(match.ID, model.MatchStatusSuccess)

		// update trip status
		s.tripRepo.UpdateStatus(match.DriverTripID, model.TripStatusMatched)
		s.tripRepo.UpdateStatus(match.PassengerTripID, model.TripStatusMatched)

		// get contact info
		contactInfo, _ := s.repo.GetContactInfo(match.ID)

		// notify both parties
		driverNotify := &model.Notification{
			UserID:  match.DriverID,
			MatchID: match.ID,
			Title:   "拼车成功",
			Content: fmt.Sprintf("恭喜！拼车成功，乘客：%s，联系电话：%s。请自行联系对方确认出行细节。", contactInfo.PassengerNickname, contactInfo.PassengerPhone),
		}
		passengerNotify := &model.Notification{
			UserID:  match.PassengerID,
			MatchID: match.ID,
			Title:   "拼车成功",
			Content: fmt.Sprintf("恭喜！拼车成功，司机：%s，联系电话：%s。请自行联系对方确认出行细节。", contactInfo.DriverNickname, contactInfo.DriverPhone),
		}

		s.notifyRepo.Create(driverNotify)
		s.notifyRepo.Create(passengerNotify)

		s.wsHub.SendToUser(match.DriverID, websocket.Message{
			Type: "match_success",
			Data: map[string]interface{}{
				"match_id":     match.ID,
				"contact":      contactInfo,
				"notification": driverNotify,
			},
		})
		s.wsHub.SendToUser(match.PassengerID, websocket.Message{
			Type: "match_success",
			Data: map[string]interface{}{
				"match_id":     match.ID,
				"contact":      contactInfo,
				"notification": passengerNotify,
			},
		})
	}
}

func (s *MatchService) GetContactInfo(matchID uint64, userID uint64) (*model.ContactInfo, error) {
	match, err := s.repo.GetByID(matchID)
	if err != nil {
		return nil, err
	}
	if match == nil {
		return nil, errors.New("匹配记录不存在")
	}
	if match.DriverID != userID && match.PassengerID != userID {
		return nil, errors.New("无权查看联系信息")
	}
	if match.Status != model.MatchStatusSuccess {
		return nil, errors.New("匹配未成功，无法查看联系信息")
	}
	return s.repo.GetContactInfo(matchID)
}

func (s *MatchService) GetByID(id uint64) (*model.Match, error) {
	return s.repo.GetByID(id)
}
