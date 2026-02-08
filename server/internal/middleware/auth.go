package middleware

import (
	"crypto/md5"
	"crypto/subtle"
	"encoding/hex"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"pinche/config"
	"pinche/internal/logger"
	"pinche/internal/model"
	"pinche/internal/service"
)

func AuthMiddleware(userService *service.UserService) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Auth failed: missing authorization header",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, model.Error(model.ErrCodeUnauthorized, "请先登录"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		userID, err := userService.ParseToken(tokenString)
		if err != nil {
			logger.Warn("Auth failed: invalid token",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
				"token", logger.MaskToken(tokenString),
				"error", err)
			c.JSON(http.StatusUnauthorized, model.Error(model.ErrCodeUnauthorized, "登录已过期，请重新登录"))
			c.Abort()
			return
		}

		c.Set("user_id", userID)
		c.Next()
	}
}

func GetUserID(c *gin.Context) uint64 {
	userID, exists := c.Get("user_id")
	if !exists || userID == nil {
		return 0
	}
	if id, ok := userID.(uint64); ok {
		return id
	}
	return 0
}

// AdminAuthMiddleware validates admin JWT token
func AdminAuthMiddleware(cfg *config.Config) gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			logger.Warn("Admin auth failed: missing authorization header",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP())
			c.JSON(http.StatusUnauthorized, model.Error(model.ErrCodeUnauthorized, "请先登录"))
			c.Abort()
			return
		}

		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims, err := ParseAdminToken(tokenString, cfg.JWT.Secret)
		if err != nil {
			logger.Warn("Admin auth failed: invalid token",
				"path", c.Request.URL.Path,
				"client_ip", c.ClientIP(),
				"error", err)
			c.JSON(http.StatusUnauthorized, model.Error(model.ErrCodeUnauthorized, "登录已过期，请重新登录"))
			c.Abort()
			return
		}

		c.Set("admin_username", claims.Username)
		c.Next()
	}
}

// AdminClaims represents admin JWT claims
type AdminClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

// GenerateAdminToken creates a JWT token for admin
func GenerateAdminToken(username string, secret string, expireHours int) (string, error) {
	claims := AdminClaims{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(expireHours) * time.Hour)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "pinche-admin",
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// ParseAdminToken validates and parses admin JWT token
func ParseAdminToken(tokenString string, secret string) (*AdminClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AdminClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})
	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*AdminClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, jwt.ErrSignatureInvalid
}

// md5Hash calculates MD5 hash of password
func md5Hash(password string) string {
	hash := md5.Sum([]byte(password))
	return hex.EncodeToString(hash[:])
}

// ValidateAdminCredentials validates admin username and password
// Password from request is already MD5 hashed by frontend
func ValidateAdminCredentials(cfg *config.Config, username, hashedPassword string) bool {
	if cfg.Admin.Password == "" {
		logger.Warn("Admin login attempted but ADMIN_PASSWORD is not configured")
		return false
	}

	// MD5 hash the configured password to compare with frontend-hashed password
	expectedHash := md5Hash(cfg.Admin.Password)

	usernameMatch := subtle.ConstantTimeCompare([]byte(username), []byte(cfg.Admin.Username)) == 1
	passwordMatch := subtle.ConstantTimeCompare([]byte(hashedPassword), []byte(expectedHash)) == 1

	return usernameMatch && passwordMatch
}
