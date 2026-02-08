package repository

import (
	"database/sql"
	"fmt"
	"strings"
	"time"

	"pinche/internal/database"
	"pinche/internal/model"
)

type TripRepository struct{}

func NewTripRepository() *TripRepository {
	return &TripRepository{}
}

// CountActiveByUserID returns count of active trips (pending/matched) for a user
func (r *TripRepository) CountActiveByUserID(userID uint64) (int, error) {
	query := `SELECT COUNT(*) FROM trips WHERE user_id = ? AND status IN (?, ?)`
	var count int
	err := database.DB.QueryRow(query, userID, model.TripStatusPending, model.TripStatusMatched).Scan(&count)
	return count, err
}

// CountTodayByUserID returns count of trips created today by a user
func (r *TripRepository) CountTodayByUserID(userID uint64) (int, error) {
	query := `SELECT COUNT(*) FROM trips WHERE user_id = ? AND DATE(created_at) = CURDATE()`
	var count int
	err := database.DB.QueryRow(query, userID).Scan(&count)
	return count, err
}

func (r *TripRepository) Create(trip *model.Trip) error {
	query := `INSERT INTO trips (user_id, trip_type, departure_city, departure_province, departure_address, departure_lat, departure_lng, 
		destination_city, destination_province, destination_address, destination_lat, destination_lng, departure_time, seats, price, remark, images, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query,
		trip.UserID, trip.TripType, trip.DepartureCity, trip.DepartureProvince, trip.DepartureAddress, trip.DepartureLat, trip.DepartureLng,
		trip.DestinationCity, trip.DestinationProvince, trip.DestinationAddress, trip.DestinationLat, trip.DestinationLng,
		trip.DepartureTime, trip.Seats, trip.Price, trip.Remark, trip.Images, trip.Status,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	trip.ID = uint64(id)
	return nil
}

func (r *TripRepository) GetByID(id uint64) (*model.Trip, error) {
	query := `SELECT t.id, t.user_id, t.trip_type, t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, t.departure_lat, t.departure_lng,
		t.destination_city, COALESCE(t.destination_province, ''), t.destination_address, t.destination_lat, t.destination_lng, t.departure_time, 
		t.seats, t.price, t.remark, COALESCE(t.images, ''), t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at,
		COALESCE(u.id, 0), COALESCE(u.open_id, ''), COALESCE(u.phone, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), COALESCE(u.gender, 0)
		FROM trips t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.id = ?`
	trip := &model.Trip{User: &model.User{}}
	err := database.DB.QueryRow(query, id).Scan(
		&trip.ID, &trip.UserID, &trip.TripType, &trip.DepartureCity, &trip.DepartureProvince, &trip.DepartureAddress, &trip.DepartureLat, &trip.DepartureLng,
		&trip.DestinationCity, &trip.DestinationProvince, &trip.DestinationAddress, &trip.DestinationLat, &trip.DestinationLng,
		&trip.DepartureTime, &trip.Seats, &trip.Price, &trip.Remark, &trip.Images, &trip.Status, &trip.ViewCount, &trip.CreatedAt, &trip.UpdatedAt,
		&trip.User.ID, &trip.User.OpenID, &trip.User.Phone, &trip.User.Nickname, &trip.User.Avatar, &trip.User.Gender,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	// set user open_id
	trip.UserOpenID = trip.User.OpenID
	// hide phone
	trip.User.Phone = ""
	return trip, nil
}

func (r *TripRepository) List(req *model.TripListReq) ([]*model.Trip, int64, error) {
	var conditions []string
	var args []interface{}

	// only show pending trips (not banned)
	conditions = append(conditions, "t.status = ?")
	args = append(args, model.TripStatusPending)

	if req.TripType > 0 {
		conditions = append(conditions, "t.trip_type = ?")
		args = append(args, req.TripType)
	}
	if req.DepartureCity != "" {
		conditions = append(conditions, "t.departure_city LIKE ?")
		args = append(args, "%"+req.DepartureCity+"%")
	}
	if req.DestinationCity != "" {
		conditions = append(conditions, "t.destination_city LIKE ?")
		args = append(args, "%"+req.DestinationCity+"%")
	}
	if req.Date != "" {
		conditions = append(conditions, "DATE(t.departure_time) = ?")
		args = append(args, req.Date)
	}

	// relevance filtering: only show trips related to user's city or province
	if req.UserCity != "" || req.UserProvince != "" {
		var relevanceConditions []string
		if req.UserCity != "" {
			relevanceConditions = append(relevanceConditions, "t.departure_city = ?")
			args = append(args, req.UserCity)
			relevanceConditions = append(relevanceConditions, "t.destination_city = ?")
			args = append(args, req.UserCity)
		}
		if req.UserProvince != "" {
			relevanceConditions = append(relevanceConditions, "COALESCE(t.departure_province, '') = ?")
			args = append(args, req.UserProvince)
			relevanceConditions = append(relevanceConditions, "COALESCE(t.destination_province, '') = ?")
			args = append(args, req.UserProvince)
		}
		if len(relevanceConditions) > 0 {
			conditions = append(conditions, "("+strings.Join(relevanceConditions, " OR ")+")")
		}
	}

	whereClause := strings.Join(conditions, " AND ")

	// count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM trips t WHERE %s", whereClause)
	var total int64
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list with sorting: city match first, then province match
	offset := (req.Page - 1) * req.PageSize
	orderClause := "t.departure_time ASC"
	if req.UserCity != "" {
		// prioritize city matches
		orderClause = fmt.Sprintf(`
			CASE 
				WHEN t.departure_city = '%s' OR t.destination_city = '%s' THEN 0
				ELSE 1
			END,
			t.departure_time ASC
		`, req.UserCity, req.UserCity)
	}

	listQuery := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.trip_type, t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, t.departure_lat, t.departure_lng,
			t.destination_city, COALESCE(t.destination_province, ''), t.destination_address, t.destination_lat, t.destination_lng, t.departure_time, 
			t.seats, t.price, t.remark, COALESCE(t.images, ''), t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at,
			COALESCE(u.id, 0), COALESCE(u.open_id, ''), COALESCE(u.phone, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), COALESCE(u.gender, 0)
		FROM trips t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE %s
		ORDER BY %s
		LIMIT ? OFFSET ?
	`, whereClause, orderClause)
	args = append(args, req.PageSize, offset)

	rows, err := database.DB.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var trips []*model.Trip
	for rows.Next() {
		trip := &model.Trip{User: &model.User{}}
		err := rows.Scan(
			&trip.ID, &trip.UserID, &trip.TripType, &trip.DepartureCity, &trip.DepartureProvince, &trip.DepartureAddress, &trip.DepartureLat, &trip.DepartureLng,
			&trip.DestinationCity, &trip.DestinationProvince, &trip.DestinationAddress, &trip.DestinationLat, &trip.DestinationLng,
			&trip.DepartureTime, &trip.Seats, &trip.Price, &trip.Remark, &trip.Images, &trip.Status, &trip.ViewCount, &trip.CreatedAt, &trip.UpdatedAt,
			&trip.User.ID, &trip.User.OpenID, &trip.User.Phone, &trip.User.Nickname, &trip.User.Avatar, &trip.User.Gender,
		)
		if err != nil {
			return nil, 0, err
		}
		// set user open_id
		trip.UserOpenID = trip.User.OpenID
		// hide phone
		trip.User.Phone = ""
		trips = append(trips, trip)
	}

	return trips, total, nil
}

func (r *TripRepository) GetByUserID(userID uint64) ([]*model.Trip, error) {
	query := `SELECT t.id, t.user_id, t.trip_type, t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, t.departure_lat, t.departure_lng,
		t.destination_city, COALESCE(t.destination_province, ''), t.destination_address, t.destination_lat, t.destination_lng, t.departure_time, t.seats, t.price, t.remark, COALESCE(t.images, ''), t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at,
		COALESCE(u.open_id, '')
		FROM trips t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.user_id = ? ORDER BY t.created_at DESC`

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*model.Trip
	for rows.Next() {
		trip := &model.Trip{}
		err := rows.Scan(
			&trip.ID, &trip.UserID, &trip.TripType, &trip.DepartureCity, &trip.DepartureProvince, &trip.DepartureAddress, &trip.DepartureLat, &trip.DepartureLng,
			&trip.DestinationCity, &trip.DestinationProvince, &trip.DestinationAddress, &trip.DestinationLat, &trip.DestinationLng,
			&trip.DepartureTime, &trip.Seats, &trip.Price, &trip.Remark, &trip.Images, &trip.Status, &trip.ViewCount, &trip.CreatedAt, &trip.UpdatedAt,
			&trip.UserOpenID,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}

func (r *TripRepository) UpdateStatus(id uint64, status int8) error {
	query := `UPDATE trips SET status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *TripRepository) FindMatchingTrips(trip *model.Trip) ([]*model.Trip, error) {
	// find opposite type trips with similar route and time
	oppositeType := model.TripTypePassenger
	if trip.TripType == model.TripTypePassenger {
		oppositeType = model.TripTypeDriver
	}

	// time window: -12h to +12h
	startTime := trip.DepartureTime.Add(-12 * time.Hour)
	endTime := trip.DepartureTime.Add(12 * time.Hour)

	query := `
		SELECT t.id, t.user_id, t.trip_type, t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, t.departure_lat, t.departure_lng,
			t.destination_city, COALESCE(t.destination_province, ''), t.destination_address, t.destination_lat, t.destination_lng, t.departure_time, 
			t.seats, t.price, t.remark, COALESCE(t.images, ''), t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at,
			COALESCE(u.id, 0), COALESCE(u.open_id, ''), COALESCE(u.phone, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), COALESCE(u.gender, 0)
		FROM trips t
		LEFT JOIN users u ON t.user_id = u.id
		WHERE t.trip_type = ?
		AND t.status = ?
		AND t.departure_city = ?
		AND t.destination_city = ?
		AND t.departure_time BETWEEN ? AND ?
		AND t.user_id != ?
		ORDER BY ABS(TIMESTAMPDIFF(MINUTE, t.departure_time, ?)) ASC
		LIMIT 20
	`

	rows, err := database.DB.Query(query, oppositeType, model.TripStatusPending,
		trip.DepartureCity, trip.DestinationCity, startTime, endTime, trip.UserID, trip.DepartureTime)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*model.Trip
	for rows.Next() {
		t := &model.Trip{User: &model.User{}}
		err := rows.Scan(
			&t.ID, &t.UserID, &t.TripType, &t.DepartureCity, &t.DepartureProvince, &t.DepartureAddress, &t.DepartureLat, &t.DepartureLng,
			&t.DestinationCity, &t.DestinationProvince, &t.DestinationAddress, &t.DestinationLat, &t.DestinationLng,
			&t.DepartureTime, &t.Seats, &t.Price, &t.Remark, &t.Images, &t.Status, &t.ViewCount, &t.CreatedAt, &t.UpdatedAt,
			&t.User.ID, &t.User.OpenID, &t.User.Phone, &t.User.Nickname, &t.User.Avatar, &t.User.Gender,
		)
		if err != nil {
			return nil, err
		}
		t.UserOpenID = t.User.OpenID
		t.User.Phone = ""
		trips = append(trips, t)
	}
	return trips, nil
}

func (r *TripRepository) Delete(id uint64, userID uint64) error {
	query := `DELETE FROM trips WHERE id = ? AND user_id = ?`
	_, err := database.DB.Exec(query, id, userID)
	return err
}

// AdminListAll returns all trips for admin panel
func (r *TripRepository) AdminListAll(req *model.AdminTripListReq) ([]*model.Trip, int64, error) {
	var conditions []string
	var args []interface{}

	if req.Search != "" {
		conditions = append(conditions, "(t.departure_city LIKE ? OR t.destination_city LIKE ?)")
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.TripType > 0 {
		conditions = append(conditions, "t.trip_type = ?")
		args = append(args, req.TripType)
	}
	if req.DepartureCity != "" {
		conditions = append(conditions, "t.departure_city LIKE ?")
		args = append(args, "%"+req.DepartureCity+"%")
	}
	if req.DestinationCity != "" {
		conditions = append(conditions, "t.destination_city LIKE ?")
		args = append(args, "%"+req.DestinationCity+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "t.status = ?")
		args = append(args, *req.Status)
	}
	if req.UserOpenID != "" {
		conditions = append(conditions, "u.open_id = ?")
		args = append(args, req.UserOpenID)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM trips t LEFT JOIN users u ON t.user_id = u.id %s", whereClause)
	var total int64
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list
	offset := (req.Page - 1) * req.PageSize
	listQuery := fmt.Sprintf(`
		SELECT t.id, t.user_id, t.trip_type, t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, t.departure_lat, t.departure_lng,
			t.destination_city, COALESCE(t.destination_province, ''), t.destination_address, t.destination_lat, t.destination_lng, t.departure_time, 
			t.seats, t.price, t.remark, COALESCE(t.images, ''), t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at,
			COALESCE(u.id, 0), COALESCE(u.open_id, ''), COALESCE(u.phone, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), COALESCE(u.gender, 0)
		FROM trips t
		LEFT JOIN users u ON t.user_id = u.id
		%s
		ORDER BY t.created_at DESC
		LIMIT ? OFFSET ?
	`, whereClause)
	args = append(args, req.PageSize, offset)

	rows, err := database.DB.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var trips []*model.Trip
	for rows.Next() {
		trip := &model.Trip{User: &model.User{}}
		err := rows.Scan(
			&trip.ID, &trip.UserID, &trip.TripType, &trip.DepartureCity, &trip.DepartureProvince, &trip.DepartureAddress, &trip.DepartureLat, &trip.DepartureLng,
			&trip.DestinationCity, &trip.DestinationProvince, &trip.DestinationAddress, &trip.DestinationLat, &trip.DestinationLng,
			&trip.DepartureTime, &trip.Seats, &trip.Price, &trip.Remark, &trip.Images, &trip.Status, &trip.ViewCount, &trip.CreatedAt, &trip.UpdatedAt,
			&trip.User.ID, &trip.User.OpenID, &trip.User.Phone, &trip.User.Nickname, &trip.User.Avatar, &trip.User.Gender,
		)
		if err != nil {
			return nil, 0, err
		}
		trip.UserOpenID = trip.User.OpenID
		trips = append(trips, trip)
	}

	return trips, total, nil
}

// IncrementViewCount increases the view count of a trip
func (r *TripRepository) IncrementViewCount(id uint64) error {
	query := `UPDATE trips SET view_count = COALESCE(view_count, 0) + 1 WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}

// CreateGrab creates a grab record for a trip
func (r *TripRepository) CreateGrab(grab *model.TripGrab) error {
	query := `INSERT INTO trip_grabs (trip_id, user_id, message) VALUES (?, ?, ?)
		ON DUPLICATE KEY UPDATE message = VALUES(message)`
	result, err := database.DB.Exec(query, grab.TripID, grab.UserID, grab.Message)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	grab.ID = uint64(id)
	return nil
}

// GetGrabsByTripID returns all grab records for a trip
func (r *TripRepository) GetGrabsByTripID(tripID uint64) ([]*model.TripGrab, error) {
	query := `SELECT g.id, g.trip_id, g.user_id, g.message, g.created_at,
		COALESCE(u.id, 0), COALESCE(u.open_id, ''), COALESCE(u.nickname, ''), COALESCE(u.avatar, ''), COALESCE(u.gender, 0)
		FROM trip_grabs g
		LEFT JOIN users u ON g.user_id = u.id
		WHERE g.trip_id = ?
		ORDER BY g.created_at DESC`
	rows, err := database.DB.Query(query, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var grabs []*model.TripGrab
	for rows.Next() {
		g := &model.TripGrab{User: &model.User{}}
		err := rows.Scan(
			&g.ID, &g.TripID, &g.UserID, &g.Message, &g.CreatedAt,
			&g.User.ID, &g.User.OpenID, &g.User.Nickname, &g.User.Avatar, &g.User.Gender,
		)
		if err != nil {
			return nil, err
		}
		grabs = append(grabs, g)
	}
	return grabs, nil
}

// UpdateTrip updates trip fields that can be changed directly (images, remark, seats, price)
func (r *TripRepository) UpdateTrip(id uint64, userID uint64, images string, remark string, seats *int, price *float64) error {
	query := `UPDATE trips SET images = ?, remark = ?`
	args := []interface{}{images, remark}

	if seats != nil {
		query += `, seats = ?`
		args = append(args, *seats)
	}
	if price != nil {
		query += `, price = ?`
		args = append(args, *price)
	}
	query += ` WHERE id = ? AND user_id = ?`
	args = append(args, id, userID)

	result, err := database.DB.Exec(query, args...)
	if err != nil {
		return err
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return fmt.Errorf("trip not found or permission denied")
	}
	return nil
}

// CreateTripUpdate creates a pending update request
func (r *TripRepository) CreateTripUpdate(update *model.TripUpdate) error {
	query := `INSERT INTO trip_updates (trip_id, user_id, update_type, old_value, new_value, status) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, update.TripID, update.UserID, update.UpdateType, update.OldValue, update.NewValue, update.Status)
	if err != nil {
		return err
	}
	id, _ := result.LastInsertId()
	update.ID = uint64(id)
	return nil
}

// GetPendingUpdatesByTripID returns pending update requests for a trip
func (r *TripRepository) GetPendingUpdatesByTripID(tripID uint64) ([]*model.TripUpdate, error) {
	query := `SELECT id, trip_id, user_id, update_type, old_value, new_value, status, reject_reason, created_at, updated_at
		FROM trip_updates WHERE trip_id = ? AND status = 0 ORDER BY created_at DESC`
	rows, err := database.DB.Query(query, tripID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var updates []*model.TripUpdate
	for rows.Next() {
		u := &model.TripUpdate{}
		err := rows.Scan(&u.ID, &u.TripID, &u.UserID, &u.UpdateType, &u.OldValue, &u.NewValue, &u.Status, &u.RejectReason, &u.CreatedAt, &u.UpdatedAt)
		if err != nil {
			return nil, err
		}
		updates = append(updates, u)
	}
	return updates, nil
}
