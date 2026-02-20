package repository

import (
	"database/sql"
	"fmt"

	"pinche/internal/database"
	"pinche/internal/model"
)

type FriendRepository struct{}

func NewFriendRepository() *FriendRepository {
	return &FriendRepository{}
}

// Create creates a new friend request
func (r *FriendRepository) Create(friend *model.Friend) error {
	query := `INSERT INTO friends (user_id, friend_id, status, message) VALUES (?, ?, ?, ?)`
	result, err := database.DB.Exec(query, friend.UserID, friend.FriendID, friend.Status, friend.Message)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	friend.ID = uint64(id)
	return nil
}

// GetByID gets a friend record by ID
func (r *FriendRepository) GetByID(id uint64) (*model.Friend, error) {
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id, u2.open_id
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE f.id = ?`
	friend := &model.Friend{}
	err := database.DB.QueryRow(query, id).Scan(
		&friend.ID, &friend.UserID, &friend.FriendID, &friend.Status, &friend.Message,
		&friend.CreatedAt, &friend.UpdatedAt, &friend.UserOpenID, &friend.FriendOpenID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return friend, nil
}

// GetByUserAndFriend gets a friend record by user_id and friend_id
func (r *FriendRepository) GetByUserAndFriend(userID, friendID uint64) (*model.Friend, error) {
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id, u2.open_id
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE f.user_id = ? AND f.friend_id = ?`
	friend := &model.Friend{}
	err := database.DB.QueryRow(query, userID, friendID).Scan(
		&friend.ID, &friend.UserID, &friend.FriendID, &friend.Status, &friend.Message,
		&friend.CreatedAt, &friend.UpdatedAt, &friend.UserOpenID, &friend.FriendOpenID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return friend, nil
}

// UpdateStatus updates the status of a friend record
func (r *FriendRepository) UpdateStatus(id uint64, status int8) error {
	query := `UPDATE friends SET status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

// Delete deletes a friend record
func (r *FriendRepository) Delete(id uint64) error {
	query := `DELETE FROM friends WHERE id = ?`
	_, err := database.DB.Exec(query, id)
	return err
}

// DeleteByUserAndFriend deletes friend records between two users (both directions)
func (r *FriendRepository) DeleteByUserAndFriend(userID, friendID uint64) error {
	query := `DELETE FROM friends WHERE (user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)`
	_, err := database.DB.Exec(query, userID, friendID, friendID, userID)
	return err
}

// GetPendingRequests gets pending friend requests received by a user
func (r *FriendRepository) GetPendingRequests(userID uint64) ([]*model.Friend, int64, error) {
	// count
	var total int64
	countQuery := `SELECT COUNT(*) FROM friends WHERE friend_id = ? AND status = ?`
	if err := database.DB.QueryRow(countQuery, userID, model.FriendStatusPending).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list with applicant user info
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id, u2.open_id,
		u1.id, u1.open_id, u1.nickname, u1.avatar, u1.gender, COALESCE(u1.city, ''), COALESCE(u1.province, '')
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE f.friend_id = ? AND f.status = ?
		ORDER BY f.created_at DESC`

	rows, err := database.DB.Query(query, userID, model.FriendStatusPending)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var friends []*model.Friend
	for rows.Next() {
		f := &model.Friend{User: &model.User{}}
		err := rows.Scan(
			&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.Message, &f.CreatedAt, &f.UpdatedAt,
			&f.UserOpenID, &f.FriendOpenID,
			&f.User.ID, &f.User.OpenID, &f.User.Nickname, &f.User.Avatar, &f.User.Gender,
			&f.User.City, &f.User.Province,
		)
		if err != nil {
			return nil, 0, err
		}
		friends = append(friends, f)
	}
	return friends, total, nil
}

// GetFriends gets all friends of a user (accepted status, both directions)
func (r *FriendRepository) GetFriends(userID uint64) ([]*model.Friend, int64, error) {
	// count: user sent request and accepted OR user received request and accepted
	countQuery := `SELECT COUNT(*) FROM friends 
		WHERE ((user_id = ? OR friend_id = ?) AND status = ?)`
	var total int64
	if err := database.DB.QueryRow(countQuery, userID, userID, model.FriendStatusAccepted).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list with friend user info
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id AS user_open_id, u2.open_id AS friend_open_id,
		CASE WHEN f.user_id = ? THEN u2.id ELSE u1.id END AS other_id,
		CASE WHEN f.user_id = ? THEN u2.open_id ELSE u1.open_id END AS other_open_id,
		CASE WHEN f.user_id = ? THEN u2.nickname ELSE u1.nickname END AS other_nickname,
		CASE WHEN f.user_id = ? THEN u2.avatar ELSE u1.avatar END AS other_avatar,
		CASE WHEN f.user_id = ? THEN u2.gender ELSE u1.gender END AS other_gender,
		CASE WHEN f.user_id = ? THEN COALESCE(u2.city, '') ELSE COALESCE(u1.city, '') END AS other_city,
		CASE WHEN f.user_id = ? THEN COALESCE(u2.province, '') ELSE COALESCE(u1.province, '') END AS other_province
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE (f.user_id = ? OR f.friend_id = ?) AND f.status = ?
		ORDER BY f.updated_at DESC`

	rows, err := database.DB.Query(query, userID, userID, userID, userID, userID, userID, userID, userID, userID, model.FriendStatusAccepted)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var friends []*model.Friend
	for rows.Next() {
		f := &model.Friend{Friend: &model.User{}}
		err := rows.Scan(
			&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.Message, &f.CreatedAt, &f.UpdatedAt,
			&f.UserOpenID, &f.FriendOpenID,
			&f.Friend.ID, &f.Friend.OpenID, &f.Friend.Nickname, &f.Friend.Avatar, &f.Friend.Gender,
			&f.Friend.City, &f.Friend.Province,
		)
		if err != nil {
			return nil, 0, err
		}
		friends = append(friends, f)
	}
	return friends, total, nil
}

// CheckFriendship checks if two users are friends (accepted status)
func (r *FriendRepository) CheckFriendship(userID, friendID uint64) (bool, error) {
	query := `SELECT COUNT(*) FROM friends 
		WHERE ((user_id = ? AND friend_id = ?) OR (user_id = ? AND friend_id = ?)) AND status = ?`
	var count int64
	err := database.DB.QueryRow(query, userID, friendID, friendID, userID, model.FriendStatusAccepted).Scan(&count)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// GetFriendshipRecord gets the friendship record between two users (any status)
func (r *FriendRepository) GetFriendshipRecord(userID, friendID uint64) (*model.Friend, error) {
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id, u2.open_id
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE (f.user_id = ? AND f.friend_id = ?) OR (f.user_id = ? AND f.friend_id = ?)`
	friend := &model.Friend{}
	err := database.DB.QueryRow(query, userID, friendID, friendID, userID).Scan(
		&friend.ID, &friend.UserID, &friend.FriendID, &friend.Status, &friend.Message,
		&friend.CreatedAt, &friend.UpdatedAt, &friend.UserOpenID, &friend.FriendOpenID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return friend, nil
}

// GetPendingRequestCount gets the count of pending friend requests received by a user
func (r *FriendRepository) GetPendingRequestCount(userID uint64) (int64, error) {
	query := `SELECT COUNT(*) FROM friends WHERE friend_id = ? AND status = ?`
	var count int64
	err := database.DB.QueryRow(query, userID, model.FriendStatusPending).Scan(&count)
	return count, err
}

// GetFriendCount gets the count of friends for a user
func (r *FriendRepository) GetFriendCount(userID uint64) (int64, error) {
	query := `SELECT COUNT(*) FROM friends WHERE (user_id = ? OR friend_id = ?) AND status = ?`
	var count int64
	err := database.DB.QueryRow(query, userID, userID, model.FriendStatusAccepted).Scan(&count)
	return count, err
}

// ResetRejectedRequest resets a rejected request to pending (for re-sending)
func (r *FriendRepository) ResetRejectedRequest(id uint64, message string) error {
	query := `UPDATE friends SET status = ?, message = ? WHERE id = ?`
	_, err := database.DB.Exec(query, model.FriendStatusPending, message, id)
	return err
}

// GetSentPendingRequests gets pending friend requests sent by a user
func (r *FriendRepository) GetSentPendingRequests(userID uint64) ([]*model.Friend, error) {
	query := `SELECT f.id, f.user_id, f.friend_id, f.status, f.message, f.created_at, f.updated_at,
		u1.open_id, u2.open_id,
		u2.id, u2.open_id, u2.nickname, u2.avatar, u2.gender, COALESCE(u2.city, ''), COALESCE(u2.province, '')
		FROM friends f
		JOIN users u1 ON f.user_id = u1.id
		JOIN users u2 ON f.friend_id = u2.id
		WHERE f.user_id = ? AND f.status = ?
		ORDER BY f.created_at DESC`

	rows, err := database.DB.Query(query, userID, model.FriendStatusPending)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var friends []*model.Friend
	for rows.Next() {
		f := &model.Friend{Friend: &model.User{}}
		err := rows.Scan(
			&f.ID, &f.UserID, &f.FriendID, &f.Status, &f.Message, &f.CreatedAt, &f.UpdatedAt,
			&f.UserOpenID, &f.FriendOpenID,
			&f.Friend.ID, &f.Friend.OpenID, &f.Friend.Nickname, &f.Friend.Avatar, &f.Friend.Gender,
			&f.Friend.City, &f.Friend.Province,
		)
		if err != nil {
			return nil, err
		}
		friends = append(friends, f)
	}
	return friends, nil
}

// GetRecentTripsByUserID gets trips from the last 3 days for a user
func (r *FriendRepository) GetRecentTripsByUserID(userID uint64, days int) ([]*model.Trip, error) {
	query := fmt.Sprintf(`SELECT t.id, t.user_id, u.open_id, t.trip_type, 
		t.departure_city, COALESCE(t.departure_province, ''), t.departure_address, 
		COALESCE(t.departure_lat, 0), COALESCE(t.departure_lng, 0),
		t.destination_city, COALESCE(t.destination_province, ''), t.destination_address,
		COALESCE(t.destination_lat, 0), COALESCE(t.destination_lng, 0),
		t.departure_time, t.seats, t.price, COALESCE(t.remark, ''), COALESCE(t.images, ''),
		t.status, COALESCE(t.view_count, 0), t.created_at, t.updated_at
		FROM trips t
		JOIN users u ON t.user_id = u.id
		WHERE t.user_id = ? AND t.departure_time >= NOW() - INTERVAL %d DAY AND t.status IN (1, 2)
		ORDER BY t.departure_time ASC`, days)

	rows, err := database.DB.Query(query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var trips []*model.Trip
	for rows.Next() {
		trip := &model.Trip{}
		err := rows.Scan(
			&trip.ID, &trip.UserID, &trip.UserOpenID, &trip.TripType,
			&trip.DepartureCity, &trip.DepartureProvince, &trip.DepartureAddress,
			&trip.DepartureLat, &trip.DepartureLng,
			&trip.DestinationCity, &trip.DestinationProvince, &trip.DestinationAddress,
			&trip.DestinationLat, &trip.DestinationLng,
			&trip.DepartureTime, &trip.Seats, &trip.Price, &trip.Remark, &trip.Images,
			&trip.Status, &trip.ViewCount, &trip.CreatedAt, &trip.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		trips = append(trips, trip)
	}
	return trips, nil
}
