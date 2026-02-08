package router

import (
	"github.com/gin-gonic/gin"
	"pinche/config"
	"pinche/internal/handler"
	"pinche/internal/middleware"
	"pinche/internal/service"
	"pinche/internal/websocket"
)

func Setup(cfg *config.Config, wsHub *websocket.Hub) *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CorsMiddleware())
	r.Use(middleware.LoggingMiddleware())

	// serve uploaded files as static
	r.Static("/uploads", "./uploads")

	// services
	userService := service.NewUserService(cfg)
	matchService := service.NewMatchService(wsHub)
	tripService := service.NewTripService(matchService, wsHub)
	notificationService := service.NewNotificationService()
	messageService := service.NewMessageService()
	announcementService := service.NewAnnouncementService()
	uploadService := service.NewUploadService(cfg)

	// handlers
	userHandler := handler.NewUserHandler(userService)
	tripHandler := handler.NewTripHandler(tripService)
	matchHandler := handler.NewMatchHandler(matchService)
	notificationHandler := handler.NewNotificationHandler(notificationService)
	messageHandler := handler.NewMessageHandler(messageService, userService, wsHub)
	wsHandler := handler.NewWebSocketHandler(wsHub, userService)
	announcementHandler := handler.NewAnnouncementHandler(announcementService)
	uploadHandler := handler.NewUploadHandler(uploadService, userService)

	// public routes
	r.POST("/api/user/register", userHandler.Register)
	r.POST("/api/user/login", userHandler.Login)
	r.GET("/api/trips", tripHandler.List)
	r.GET("/api/trips/:id", tripHandler.GetByID)
	r.GET("/api/announcements", announcementHandler.GetActiveAnnouncements)

	// websocket
	r.GET("/ws", wsHandler.HandleConnection)

	// protected routes
	auth := r.Group("/api")
	auth.Use(middleware.AuthMiddleware(userService))
	{
		// user
		auth.GET("/user/profile", userHandler.GetProfile)
		auth.PUT("/user/profile", userHandler.UpdateProfile)

		// trips
		auth.POST("/trips", tripHandler.Create)
		auth.GET("/trips/my", tripHandler.GetMyTrips)
		auth.GET("/trips/my/:id", tripHandler.GetMyTripDetail)
		auth.PUT("/trips/:id", tripHandler.UpdateTrip)
		auth.PUT("/trips/:id/cancel", tripHandler.Cancel)
		auth.PUT("/trips/:id/complete", tripHandler.Complete)
		auth.DELETE("/trips/:id", tripHandler.Delete)
		auth.POST("/trips/:id/grab", tripHandler.GrabTrip)

		// matches
		auth.GET("/matches", matchHandler.GetMyMatches)
		auth.GET("/matches/:id", matchHandler.GetByID)
		auth.POST("/matches/:id/confirm", matchHandler.Confirm)
		auth.GET("/matches/:id/contact", matchHandler.GetContactInfo)

		// notifications
		auth.GET("/notifications", notificationHandler.GetList)
		auth.PUT("/notifications/:id/read", notificationHandler.MarkAsRead)
		auth.PUT("/notifications/read-all", notificationHandler.MarkAllAsRead)

		// messages
		auth.POST("/messages", messageHandler.SendMessage)
		auth.GET("/messages", messageHandler.GetConversationMessages)
		auth.GET("/conversations", messageHandler.GetConversations)
		auth.PUT("/messages/read", messageHandler.MarkAsRead)
		auth.GET("/messages/unread-count", messageHandler.GetUnreadCount)

		// upload
		auth.POST("/upload/image", uploadHandler.UploadImage)
		auth.GET("/resource/url", uploadHandler.GetSignedURL)
	}

	// admin routes (should add admin auth middleware in production)
	admin := r.Group("/api/admin")
	{
		admin.GET("/announcements", announcementHandler.ListAll)
		admin.POST("/announcements", announcementHandler.Create)
		admin.PUT("/announcements/:id", announcementHandler.Update)
		admin.DELETE("/announcements/:id", announcementHandler.Delete)

		admin.GET("/users", userHandler.AdminListUsers)
		admin.POST("/users/:id/ban", userHandler.AdminBanUser)
		admin.POST("/users/:id/unban", userHandler.AdminUnbanUser)

		admin.GET("/trips", tripHandler.AdminListTrips)
		admin.POST("/trips/:id/ban", tripHandler.AdminBanTrip)
		admin.POST("/trips/:id/unban", tripHandler.AdminUnbanTrip)

		admin.GET("/stats", userHandler.AdminGetStats)
	}

	return r
}
