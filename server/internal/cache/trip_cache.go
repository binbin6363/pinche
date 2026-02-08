package cache

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"

	"github.com/redis/go-redis/v9"
	"pinche/internal/logger"
	"pinche/internal/model"
)

// TripCache handles trip caching with Cache Aside pattern
type TripCache struct{}

func NewTripCache() *TripCache {
	return &TripCache{}
}

// GetTrip gets trip from cache
// Returns nil, nil if not found (cache miss)
func (c *TripCache) GetTrip(id uint64) (*model.Trip, error) {
	key := TripKey(id)
	var trip model.Trip
	err := Get(key, &trip)
	if err == redis.Nil {
		logger.Debug("Cache miss for trip", "trip_id", id)
		return nil, nil
	}
	if err != nil {
		logger.Error("Cache get trip failed", "trip_id", id, "error", err)
		return nil, err
	}
	logger.Debug("Cache hit for trip", "trip_id", id)
	return &trip, nil
}

// SetTrip stores trip in cache
func (c *TripCache) SetTrip(trip *model.Trip) error {
	if trip == nil {
		return nil
	}
	key := TripKey(trip.ID)
	if err := Set(key, trip, TripDetailTTL); err != nil {
		logger.Error("Cache set trip failed", "trip_id", trip.ID, "error", err)
		return err
	}
	logger.Debug("Cache set trip", "trip_id", trip.ID, "ttl", TripDetailTTL)
	return nil
}

// InvalidateTrip removes trip from cache
func (c *TripCache) InvalidateTrip(id uint64) error {
	key := TripKey(id)
	if err := Delete(key); err != nil {
		logger.Error("Cache invalidate trip failed", "trip_id", id, "error", err)
		return err
	}
	logger.Debug("Cache invalidated trip", "trip_id", id)
	return nil
}

// TripListResult holds cached list result
type TripListResult struct {
	List  []*model.Trip `json:"list"`
	Total int64         `json:"total"`
}

// GetTripList gets trip list from cache
func (c *TripCache) GetTripList(req *model.TripListReq) (*TripListResult, error) {
	key := c.listCacheKey(req)
	var result TripListResult
	err := Get(key, &result)
	if err == redis.Nil {
		logger.Debug("Cache miss for trip list", "key", key)
		return nil, nil
	}
	if err != nil {
		logger.Error("Cache get trip list failed", "key", key, "error", err)
		return nil, err
	}
	logger.Debug("Cache hit for trip list", "key", key)
	return &result, nil
}

// SetTripList stores trip list in cache
func (c *TripCache) SetTripList(req *model.TripListReq, list []*model.Trip, total int64) error {
	key := c.listCacheKey(req)
	result := &TripListResult{
		List:  list,
		Total: total,
	}
	if err := Set(key, result, TripListTTL); err != nil {
		logger.Error("Cache set trip list failed", "key", key, "error", err)
		return err
	}
	logger.Debug("Cache set trip list", "key", key, "count", len(list), "ttl", TripListTTL)
	return nil
}

// InvalidateTripLists removes all trip list caches
func (c *TripCache) InvalidateTripLists() error {
	pattern := KeyPrefixTripList + "*"
	if err := DeleteByPattern(pattern); err != nil {
		logger.Error("Cache invalidate trip lists failed", "error", err)
		return err
	}
	logger.Debug("Cache invalidated all trip lists")
	return nil
}

// listCacheKey generates cache key for trip list request
func (c *TripCache) listCacheKey(req *model.TripListReq) string {
	data, _ := json.Marshal(req)
	hash := md5.Sum(data)
	return fmt.Sprintf("%s%s", KeyPrefixTripList, hex.EncodeToString(hash[:]))
}
