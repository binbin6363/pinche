package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	COS      COSConfig
	Log      LogConfig
	Admin    AdminConfig
}

type AdminConfig struct {
	Username string
	Password string
}

type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type LogConfig struct {
	Level      string // debug, info, warn, error
	Filename   string // log file path
	MaxSize    int    // max size in MB before rotation
	MaxBackups int    // max number of old log files
	MaxAge     int    // max days to retain old log files
	Compress   bool   // compress rotated files
	Console    bool   // also output to console
}

type ServerConfig struct {
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
}

type JWTConfig struct {
	Secret     string
	ExpireHour int
}

type COSConfig struct {
	SecretID  string
	SecretKey string
	Bucket    string
	Region    string
}

func Load() *Config {
	// load .env file if exists
	godotenv.Load()

	return &Config{
		Server: ServerConfig{
			Port: getEnv("SERVER_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "127.0.0.1"),
			Port:     getEnv("DB_PORT", "3306"),
			User:     getEnv("DB_USER", "root"),
			Password: getEnv("DB_PASSWORD", ""),
			DBName:   getEnv("DB_NAME", "pinche"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST", "127.0.0.1"),
			Port:     getEnv("REDIS_PORT", "6379"),
			Password: getEnv("REDIS_PASSWORD", ""),
			DB:       getEnvInt("REDIS_DB", 0),
		},
		JWT: JWTConfig{
			Secret:     getEnv("JWT_SECRET", "pinche-secret-key-2024"),
			ExpireHour: getEnvInt("JWT_EXPIRE_HOUR", 168),
		},
		COS: COSConfig{
			SecretID:  getEnv("COS_SECRET_ID", ""),
			SecretKey: getEnv("COS_SECRET_KEY", ""),
			Bucket:    getEnv("COS_BUCKET", ""),
			Region:    getEnv("COS_REGION", "ap-shanghai"),
		},
		Log: LogConfig{
			Level:      getEnv("LOG_LEVEL", "info"),
			Filename:   getEnv("LOG_FILENAME", "logs/app.log"),
			MaxSize:    getEnvInt("LOG_MAX_SIZE", 200),
			MaxBackups: getEnvInt("LOG_MAX_BACKUPS", 7),
			MaxAge:     getEnvInt("LOG_MAX_AGE", 30),
			Compress:   getEnvBool("LOG_COMPRESS", true),
			Console:    getEnvBool("LOG_CONSOLE", true),
		},
		Admin: AdminConfig{
			Username: getEnv("ADMIN_USERNAME", "admin"),
			Password: getEnv("ADMIN_PASSWORD", ""),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intVal, err := strconv.Atoi(value); err == nil {
			return intVal
		}
	}
	return defaultValue
}

func getEnvBool(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if boolVal, err := strconv.ParseBool(value); err == nil {
			return boolVal
		}
	}
	return defaultValue
}
