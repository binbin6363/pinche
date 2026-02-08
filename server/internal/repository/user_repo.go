package repository

import (
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"fmt"
	"strings"

	"pinche/internal/database"
	"pinche/internal/model"
)

type UserRepository struct{}

func NewUserRepository() *UserRepository {
	return &UserRepository{}
}

// generateOpenID generates a unique 24-character open_id
func generateOpenID() string {
	bytes := make([]byte, 12)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}

func (r *UserRepository) Create(user *model.User) error {
	// generate unique open_id
	user.OpenID = generateOpenID()

	query := `INSERT INTO users (open_id, phone, password, nickname, avatar, gender, status, city, province) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`
	result, err := database.DB.Exec(query, user.OpenID, user.Phone, user.Password, user.Nickname, user.Avatar, user.Gender, 0, user.City, user.Province)
	if err != nil {
		return err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return err
	}
	user.ID = uint64(id)
	return nil
}

func (r *UserRepository) GetByID(id uint64) (*model.User, error) {
	query := `SELECT id, open_id, phone, password, nickname, avatar, gender, COALESCE(status, 0), COALESCE(city, ''), COALESCE(province, ''), created_at, updated_at FROM users WHERE id = ?`
	user := &model.User{}
	err := database.DB.QueryRow(query, id).Scan(
		&user.ID, &user.OpenID, &user.Phone, &user.Password, &user.Nickname,
		&user.Avatar, &user.Gender, &user.Status, &user.City, &user.Province, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByOpenID(openID string) (*model.User, error) {
	query := `SELECT id, open_id, phone, password, nickname, avatar, gender, COALESCE(status, 0), COALESCE(city, ''), COALESCE(province, ''), created_at, updated_at FROM users WHERE open_id = ?`
	user := &model.User{}
	err := database.DB.QueryRow(query, openID).Scan(
		&user.ID, &user.OpenID, &user.Phone, &user.Password, &user.Nickname,
		&user.Avatar, &user.Gender, &user.Status, &user.City, &user.Province, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) GetByPhone(phone string) (*model.User, error) {
	query := `SELECT id, open_id, phone, password, nickname, avatar, gender, COALESCE(status, 0), COALESCE(city, ''), COALESCE(province, ''), created_at, updated_at FROM users WHERE phone = ?`
	user := &model.User{}
	err := database.DB.QueryRow(query, phone).Scan(
		&user.ID, &user.OpenID, &user.Phone, &user.Password, &user.Nickname,
		&user.Avatar, &user.Gender, &user.Status, &user.City, &user.Province, &user.CreatedAt, &user.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (r *UserRepository) Update(user *model.User) error {
	query := `UPDATE users SET nickname = ?, avatar = ?, gender = ?, city = ?, province = ? WHERE id = ?`
	_, err := database.DB.Exec(query, user.Nickname, user.Avatar, user.Gender, user.City, user.Province, user.ID)
	return err
}

func (r *UserRepository) UpdateStatus(id uint64, status int8) error {
	query := `UPDATE users SET status = ? WHERE id = ?`
	_, err := database.DB.Exec(query, status, id)
	return err
}

func (r *UserRepository) UpdateStatusByOpenID(openID string, status int8) error {
	query := `UPDATE users SET status = ? WHERE open_id = ?`
	_, err := database.DB.Exec(query, status, openID)
	return err
}

// ListAll returns all users for admin panel
func (r *UserRepository) ListAll(req *model.AdminUserListReq) ([]*model.User, int64, error) {
	var conditions []string
	var args []interface{}

	if req.Search != "" {
		conditions = append(conditions, "(phone LIKE ? OR nickname LIKE ?)")
		args = append(args, "%"+req.Search+"%", "%"+req.Search+"%")
	}
	if req.Phone != "" {
		conditions = append(conditions, "phone LIKE ?")
		args = append(args, "%"+req.Phone+"%")
	}
	if req.Nickname != "" {
		conditions = append(conditions, "nickname LIKE ?")
		args = append(args, "%"+req.Nickname+"%")
	}
	if req.Status != nil {
		conditions = append(conditions, "COALESCE(status, 0) = ?")
		args = append(args, *req.Status)
	}

	whereClause := ""
	if len(conditions) > 0 {
		whereClause = "WHERE " + strings.Join(conditions, " AND ")
	}

	// count
	countQuery := fmt.Sprintf("SELECT COUNT(*) FROM users %s", whereClause)
	var total int64
	if err := database.DB.QueryRow(countQuery, args...).Scan(&total); err != nil {
		return nil, 0, err
	}

	// list
	offset := (req.Page - 1) * req.PageSize
	listQuery := fmt.Sprintf(`SELECT id, open_id, phone, password, nickname, avatar, gender, COALESCE(status, 0), COALESCE(city, ''), COALESCE(province, ''), created_at, updated_at 
		FROM users %s ORDER BY created_at DESC LIMIT ? OFFSET ?`, whereClause)
	args = append(args, req.PageSize, offset)

	rows, err := database.DB.Query(listQuery, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []*model.User
	for rows.Next() {
		user := &model.User{}
		err := rows.Scan(&user.ID, &user.OpenID, &user.Phone, &user.Password, &user.Nickname, &user.Avatar, &user.Gender, &user.Status, &user.City, &user.Province, &user.CreatedAt, &user.UpdatedAt)
		if err != nil {
			return nil, 0, err
		}
		user.Password = ""
		users = append(users, user)
	}
	return users, total, nil
}

// GetStats returns statistics for admin dashboard
func (r *UserRepository) GetStats() (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// total users
	var totalUsers int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM users").Scan(&totalUsers); err != nil {
		return nil, err
	}
	stats["total_users"] = totalUsers

	// banned users
	var bannedUsers int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE status = 1").Scan(&bannedUsers); err != nil {
		return nil, err
	}
	stats["banned_users"] = bannedUsers

	// total trips
	var totalTrips int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM trips").Scan(&totalTrips); err != nil {
		return nil, err
	}
	stats["total_trips"] = totalTrips

	// active trips (pending)
	var activeTrips int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM trips WHERE status = 1").Scan(&activeTrips); err != nil {
		return nil, err
	}
	stats["active_trips"] = activeTrips

	// banned trips
	var bannedTrips int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM trips WHERE status = 5").Scan(&bannedTrips); err != nil {
		return nil, err
	}
	stats["banned_trips"] = bannedTrips

	// total matches
	var totalMatches int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM matches").Scan(&totalMatches); err != nil {
		return nil, err
	}
	stats["total_matches"] = totalMatches

	// successful matches
	var successMatches int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM matches WHERE status = 1").Scan(&successMatches); err != nil {
		return nil, err
	}
	stats["success_matches"] = successMatches

	// today new users
	var todayUsers int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM users WHERE DATE(created_at) = CURDATE()").Scan(&todayUsers); err != nil {
		return nil, err
	}
	stats["today_users"] = todayUsers

	// today new trips
	var todayTrips int64
	if err := database.DB.QueryRow("SELECT COUNT(*) FROM trips WHERE DATE(created_at) = CURDATE()").Scan(&todayTrips); err != nil {
		return nil, err
	}
	stats["today_trips"] = todayTrips

	return stats, nil
}
