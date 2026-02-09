package service

import (
	"errors"
	"time"

	"pinche/config"
	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/repository"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo   *repository.UserRepository
	config *config.Config
}

func NewUserService(cfg *config.Config) *UserService {
	return &UserService{
		repo:   repository.NewUserRepository(),
		config: cfg,
	}
}

func (s *UserService) Register(req *model.UserRegisterReq) (*model.User, error) {
	// check phone exists
	existing, err := s.repo.GetByPhone(req.Phone)
	if err != nil {
		logger.Error("Check phone exists failed", "phone", logger.MaskPhone(req.Phone), "error", err)
		return nil, err
	}
	if existing != nil {
		logger.Warn("Phone already registered", "phone", logger.MaskPhone(req.Phone))
		return nil, errors.New("手机号已注册")
	}

	// password is already MD5 hashed from frontend, bcrypt it for storage
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("Hash password failed", "error", err)
		return nil, err
	}

	user := &model.User{
		Phone:    req.Phone,
		Password: string(hashedPassword),
		Nickname: req.Nickname,
	}

	if err := s.repo.Create(user); err != nil {
		logger.Error("Create user failed", "phone", logger.MaskPhone(req.Phone), "error", err)
		return nil, err
	}

	logger.Info("User registered", "user_id", user.ID, "phone", logger.MaskPhone(req.Phone))
	return user, nil
}

func (s *UserService) Login(req *model.UserLoginReq) (*model.UserLoginResp, error) {
	user, err := s.repo.GetByPhone(req.Phone)
	if err != nil {
		logger.Error("Get user by phone failed", "phone", logger.MaskPhone(req.Phone), "error", err)
		return nil, err
	}
	if user == nil {
		logger.Warn("Login failed: user not found", "phone", logger.MaskPhone(req.Phone))
		return nil, errors.New("用户不存在，请先注册")
	}

	// check if user is banned
	if user.Status == 1 {
		logger.Warn("Login failed: user banned", "user_id", user.ID, "phone", logger.MaskPhone(req.Phone))
		return nil, errors.New("账号已被封禁，请联系客服")
	}

	// validate stored password is a valid bcrypt hash (must start with $2a$ or $2b$ and be 60 chars)
	storedPwd := user.Password
	if len(storedPwd) != 60 || (storedPwd[:4] != "$2a$" && storedPwd[:4] != "$2b$") {
		logger.Error("Invalid bcrypt hash in database", "user_id", user.ID, "phone", logger.MaskPhone(req.Phone), "pwd_len", len(storedPwd))
		return nil, errors.New("密码数据异常，请联系客服重置")
	}

	// password from frontend is MD5 hashed, compare with bcrypt stored password
	if err := bcrypt.CompareHashAndPassword([]byte(storedPwd), []byte(req.Password)); err != nil {
		logger.Warn("Login failed: wrong password", "user_id", user.ID, "phone", logger.MaskPhone(req.Phone))
		return nil, errors.New("用户名或密码错误")
	}

	// generate token with internal ID
	token, err := s.generateToken(user.ID)
	if err != nil {
		logger.Error("Generate token failed", "user_id", user.ID, "error", err)
		return nil, err
	}

	logger.Info("User logged in", "user_id", user.ID, "phone", logger.MaskPhone(req.Phone))

	return &model.UserLoginResp{
		Token: token,
		User:  user,
	}, nil
}

func (s *UserService) GetByID(id uint64) (*model.User, error) {
	return s.repo.GetByID(id)
}

func (s *UserService) GetByOpenID(openID string) (*model.User, error) {
	return s.repo.GetByOpenID(openID)
}

func (s *UserService) Update(userID uint64, req *model.UserUpdateReq) (*model.User, error) {
	user, err := s.repo.GetByID(userID)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("用户不存在")
	}

	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if req.Gender != nil {
		user.Gender = *req.Gender
	}
	if req.City != "" {
		user.City = req.City
	}
	if req.Province != "" {
		user.Province = req.Province
	}

	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *UserService) generateToken(userID uint64) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(time.Duration(s.config.JWT.ExpireHour) * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.config.JWT.Secret))
}

func (s *UserService) ParseToken(tokenString string) (uint64, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.config.JWT.Secret), nil
	})
	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint64(claims["user_id"].(float64))
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

// Admin functions

func (s *UserService) AdminListUsers(req *model.AdminUserListReq) (*model.AdminUserListResp, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 || req.PageSize > 100 {
		req.PageSize = 20
	}

	users, total, err := s.repo.ListAll(req)
	if err != nil {
		return nil, err
	}

	return &model.AdminUserListResp{
		List:  users,
		Total: total,
	}, nil
}

func (s *UserService) AdminBanUser(openID string) error {
	logger.Info("Admin banning user", "open_id", openID)
	return s.repo.UpdateStatusByOpenID(openID, 1)
}

func (s *UserService) AdminUnbanUser(openID string) error {
	logger.Info("Admin unbanning user", "open_id", openID)
	return s.repo.UpdateStatusByOpenID(openID, 0)
}

func (s *UserService) AdminGetStats() (map[string]interface{}, error) {
	return s.repo.GetStats()
}
