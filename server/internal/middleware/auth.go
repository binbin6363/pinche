package middleware

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
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
