package repository

import (
	"database/sql"
	"pinche/internal/database"
	"pinche/internal/model"
)

type MatchRepository struct{}

func NewMatchRepository() *MatchRepository {
	return &MatchRepository{}
}

func (r *MatchRepository) Create(match *model.Match) error {
	query := `INSERT INTO matches (driver_trip_id, passenger_trip_id, driver_id, passenger_id, match_score, driver_status, passenger_status, status) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query,
		match.DriverTripID, match.PassengerTripID, match.DriverID, match.PassengerID,
		match.MatchScore, match.DriverStatus, match.PassengerStatus, match.Status,
	)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	match.ID = uint64(id)
	return nil
}

func (r *MatchRepository) GetByID(id uint64) (*model.Match, error) {
	query := `SELECT m.id, m.driver_trip_id, m.passenger_trip_id, m.driver_id, m.passenger_id, m.match_score, 
		m.driver_status, m.passenger_status, m.status, m.created_at, m.updated_at,
		COALESCE(d.open_id, ''), COALESCE(p.open_id, '')
		FROM matches m
		LEFT JOIN users d ON m.driver_id = d.id
		LEFT JOIN users p ON m.passenger_id = p.id
		WHERE m.id = ?`
	match := &model.Match{}
	err := database.DB.QueryRow(query, id).Scan(
		&match.ID, &match.DriverTripID, &match.PassengerTripID, &match.DriverID, &match.PassengerID,
		&match.MatchScore, &match.DriverStatus, &match.PassengerStatus, &match.Status, &match.CreatedAt, &match.UpdatedAt,
		&match.DriverOpenID, &match.PassengerOpenID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (r *MatchRepository) GetByTrips(driverTripID, passengerTripID uint64) (*model.Match, error) {
	query := `SELECT m.id, m.driver_trip_id, m.passenger_trip_id, m.driver_id, m.passenger_id, m.match_score, 
		m.driver_status, m.passenger_status, m.status, m.created_at, m.updated_at,
		COALESCE(d.open_id, ''), COALESCE(p.open_id, '')
		FROM matches m
		LEFT JOIN users d ON m.driver_id = d.id
		LEFT JOIN users p ON m.passenger_id = p.id
		WHERE m.driver_trip_id = ? AND m.passenger_trip_id = ?`
	match := &model.Match{}
	err := database.DB.QueryRow(query, driverTripID, passengerTripID).Scan(
		&match.ID, &match.DriverTripID, &match.PassengerTripID, &match.DriverID, &match.PassengerID,
		&match.MatchScore, &match.DriverStatus, &match.PassengerStatus, &match.Status, &match.CreatedAt, &match.UpdatedAt,
		&match.DriverOpenID, &match.PassengerOpenID,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return match, nil
}

func (r *MatchRepository) GetByUserID(userID uint64) ([]*model.Match, error) {
	query := `
		SELECT m.id, m.driver_trip_id, m.passenger_trip_id, m.driver_id, m.passenger_id, m.match_score,
			m.driver_status, m.passenger_status, m.status, m.created_at, m.updated_at,
			dt.departure_city, dt.departure_address, dt.destination_city, dt.destination_address, dt.departure_time, dt.seats, dt.price,
			pt.departure_city, pt.departure_address, pt.destination_city, pt.destination_address, pt.departure_time, pt.seats,
			d.id, COALESCE(d.open_id, ''), d.nickname, d.avatar,
			p.id, COALESCE(p.open_id, ''), p.nickname, p.avatar
		FROM matches m
		LEFT JOIN trips dt ON m.driver_trip_id = dt.id
		LEFT JOIN trips pt ON m.passenger_trip_id = pt.id
		LEFT JOIN users d ON m.driver_id = d.id
		LEFT JOIN users p ON m.passenger_id = p.id
		WHERE m.driver_id = ? OR m.passenger_id = ?
		ORDER BY m.created_at DESC
	`
	rows, err := database.DB.Query(query, userID, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var matches []*model.Match
	for rows.Next() {
		m := &model.Match{
			DriverTrip:    &model.Trip{},
			PassengerTrip: &model.Trip{},
			Driver:        &model.User{},
			Passenger:     &model.User{},
		}
		err := rows.Scan(
			&m.ID, &m.DriverTripID, &m.PassengerTripID, &m.DriverID, &m.PassengerID, &m.MatchScore,
			&m.DriverStatus, &m.PassengerStatus, &m.Status, &m.CreatedAt, &m.UpdatedAt,
			&m.DriverTrip.DepartureCity, &m.DriverTrip.DepartureAddress, &m.DriverTrip.DestinationCity,
			&m.DriverTrip.DestinationAddress, &m.DriverTrip.DepartureTime, &m.DriverTrip.Seats, &m.DriverTrip.Price,
			&m.PassengerTrip.DepartureCity, &m.PassengerTrip.DepartureAddress, &m.PassengerTrip.DestinationCity,
			&m.PassengerTrip.DestinationAddress, &m.PassengerTrip.DepartureTime, &m.PassengerTrip.Seats,
			&m.Driver.ID, &m.Driver.OpenID, &m.Driver.Nickname, &m.Driver.Avatar,
			&m.Passenger.ID, &m.Passenger.OpenID, &m.Passenger.Nickname, &m.Passenger.Avatar,
		)
		if err != nil {
			return nil, err
		}
		// set open_id fields
		m.DriverOpenID = m.Driver.OpenID
		m.PassengerOpenID = m.Passenger.OpenID
		matches = append(matches, m)
	}
	return matches, nil
}

func (r *MatchRepository) UpdateDriverStatus(id uint64, status int8) error {
	query := `UPDATE matches SET driver_status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *MatchRepository) UpdatePassengerStatus(id uint64, status int8) error {
	query := `UPDATE matches SET passenger_status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *MatchRepository) UpdateStatus(id uint64, status int8) error {
	query := `UPDATE matches SET status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *MatchRepository) GetContactInfo(matchID uint64) (*model.ContactInfo, error) {
	query := `
		SELECT d.phone, d.nickname, p.phone, p.nickname
		FROM matches m
		LEFT JOIN users d ON m.driver_id = d.id
		LEFT JOIN users p ON m.passenger_id = p.id
		WHERE m.id = ? AND m.status = ?
	`
	info := &model.ContactInfo{}
	err := database.DB.QueryRow(query, matchID, model.MatchStatusSuccess).Scan(
		&info.DriverPhone, &info.DriverNickname, &info.PassengerPhone, &info.PassengerNickname,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return info, nil
}
