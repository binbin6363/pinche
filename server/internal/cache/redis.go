package cache

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"pinche/config"
	"pinche/internal/logger"
)

var (
	Client *redis.Client
	ctx    = context.Background()
)

// cache key prefixes
const (
	KeyPrefixTrip     = "trip:"
	KeyPrefixTripList = "trip_list:"
)

// default TTL
const (
	TripDetailTTL = 10 * time.Minute
	TripListTTL   = 5 * time.Minute
)

func Init(cfg *config.RedisConfig) error {
	Client = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", cfg.Host, cfg.Port),
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if err := Client.Ping(ctx).Err(); err != nil {
		logger.Error("Failed to connect to Redis", "error", err)
		return fmt.Errorf("failed to connect to redis: %w", err)
	}

	logger.Info("Redis connected", "host", cfg.Host, "port", cfg.Port, "db", cfg.DB)
	return nil
}

func Close() {
	if Client != nil {
		Client.Close()
		logger.Info("Redis connection closed")
	}
}

// Get retrieves value from cache and unmarshal to target
func Get(key string, target interface{}) error {
	val, err := Client.Get(ctx, key).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(val), target)
}

// Set stores value in cache with TTL
func Set(key string, value interface{}, ttl time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return err
	}
	return Client.Set(ctx, key, data, ttl).Err()
}

// Delete removes key from cache
func Delete(keys ...string) error {
	if len(keys) == 0 {
		return nil
	}
	return Client.Del(ctx, keys...).Err()
}

// DeleteByPattern removes all keys matching pattern
func DeleteByPattern(pattern string) error {
	iter := Client.Scan(ctx, 0, pattern, 0).Iterator()
	var keys []string
	for iter.Next(ctx) {
		keys = append(keys, iter.Val())
	}
	if err := iter.Err(); err != nil {
		return err
	}
	if len(keys) > 0 {
		return Client.Del(ctx, keys...).Err()
	}
	return nil
}

// Exists checks if key exists
func Exists(key string) (bool, error) {
	n, err := Client.Exists(ctx, key).Result()
	return n > 0, err
}

// trip cache keys
func TripKey(id uint64) string {
	return fmt.Sprintf("%s%d", KeyPrefixTrip, id)
}

func TripListKey(params string) string {
	return fmt.Sprintf("%s%s", KeyPrefixTripList, params)
}
