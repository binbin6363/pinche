package main

import (
	"pinche/config"
	"pinche/internal/cache"
	"pinche/internal/database"
	"pinche/internal/logger"
	"pinche/internal/router"
	"pinche/internal/websocket"
)

func main() {
	cfg := config.Load()

	// init logger
	if err := logger.Init(&logger.Config{
		Level:      cfg.Log.Level,
		Filename:   cfg.Log.Filename,
		MaxSize:    cfg.Log.MaxSize,
		MaxBackups: cfg.Log.MaxBackups,
		MaxAge:     cfg.Log.MaxAge,
		Compress:   cfg.Log.Compress,
		Console:    cfg.Log.Console,
	}); err != nil {
		panic("Failed to init logger: " + err.Error())
	}
	defer logger.Sync()
	logger.Info("Logger initialized", "level", cfg.Log.Level)

	// init database
	if err := database.Init(&cfg.Database); err != nil {
		logger.Fatal("Failed to init database", "error", err)
	}
	defer database.Close()
	logger.Info("Database connected", "host", cfg.Database.Host, "db", cfg.Database.DBName)

	// init redis cache
	if err := cache.Init(&cfg.Redis); err != nil {
		logger.Fatal("Failed to init redis", "error", err)
	}
	defer cache.Close()

	// init websocket hub
	wsHub := websocket.NewHub()
	go wsHub.Run()
	logger.Info("WebSocket hub started")

	// setup router
	r := router.Setup(cfg, wsHub)

	// start server
	logger.Info("Server starting", "port", cfg.Server.Port)
	if err := r.Run(":" + cfg.Server.Port); err != nil {
		logger.Fatal("Failed to start server", "error", err)
	}
}
