package database

import (
	"database/sql"
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"pinche/config"
	"pinche/internal/logger"
)

var DB *sql.DB

func Init(cfg *config.DatabaseConfig) error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.DBName)

	var err error
	DB, err = sql.Open("mysql", dsn)
	if err != nil {
		logger.Error("Failed to open database", "error", err)
		return fmt.Errorf("failed to open database: %w", err)
	}

	DB.SetMaxOpenConns(100)
	DB.SetMaxIdleConns(10)
	DB.SetConnMaxLifetime(time.Hour)

	if err = DB.Ping(); err != nil {
		logger.Error("Failed to ping database", "error", err)
		return fmt.Errorf("failed to ping database: %w", err)
	}

	logger.Debug("Database connection pool configured",
		"max_open_conns", 100,
		"max_idle_conns", 10,
		"conn_max_lifetime", "1h")

	return nil
}

func Close() {
	if DB != nil {
		DB.Close()
		logger.Info("Database connection closed")
	}
}
