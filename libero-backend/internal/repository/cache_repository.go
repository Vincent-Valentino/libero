package repository

import (
	"libero-backend/internal/models"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CacheRepository defines the interface for cache data operations
type CacheRepository interface {
	// Generic cache operations
	Get(key string) (*models.CacheItem, error)
	GetWithVersion(key string, version string) (*models.CacheItem, error)
	Set(key string, value []byte, ttl time.Duration) error
	SetWithMetadata(key string, item models.CacheItem) error
	UpdateVersion(key string, version string) error

	// Fixtures operations
	GetCachedFixtures(competitionCode, dataType string) (*models.CachedFixtures, error)
	StoreCachedFixtures(competitionCode, dataType string, data models.JSONB, ttl time.Duration) error
	GetCachedFixturesIgnoringExpiry(competitionCode, dataType string) (*models.CachedFixtures, error)

	// Today's fixtures operations
	GetCachedTodayFixtures() (*models.CachedTodayFixtures, error)
	StoreCachedTodayFixtures(data models.JSONB, ttl time.Duration) error
	GetCachedTodayFixturesIgnoringExpiry() (*models.CachedTodayFixtures, error)

	// Common operations
	CleanExpiredCache() error
}

// cacheRepository implements the CacheRepository interface.
type cacheRepository struct {
	db *gorm.DB
}

// NewCacheRepository creates a new cache repository instance
func NewCacheRepository(db *gorm.DB) CacheRepository {
	return &cacheRepository{db: db}
}

// GetCachedFixtures retrieves cached fixtures data for a competition if not expired
func (r *cacheRepository) GetCachedFixtures(competitionCode, dataType string) (*models.CachedFixtures, error) {
	var cache models.CachedFixtures

	err := r.db.Where("competition_code = ? AND data_type = ? AND expires_at > ?",
		competitionCode, dataType, time.Now()).First(&cache).Error

	if err != nil {
		return nil, err // Return nil if not found or error
	}

	return &cache, nil
}

// GetCachedFixturesIgnoringExpiry retrieves cached fixtures data even if expired
func (r *cacheRepository) GetCachedFixturesIgnoringExpiry(competitionCode, dataType string) (*models.CachedFixtures, error) {
	var cache models.CachedFixtures

	err := r.db.Where("competition_code = ? AND data_type = ?",
		competitionCode, dataType).First(&cache).Error

	if err != nil {
		return nil, err // Return nil if not found or error
	}

	return &cache, nil
}

// StoreCachedFixtures stores fixtures data with an expiration time
func (r *cacheRepository) StoreCachedFixtures(competitionCode, dataType string, data models.JSONB, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	cache := models.CachedFixtures{
		CompetitionCode: competitionCode,
		DataType:        dataType,
		Data:            data,
		CachedData: models.CachedData{
			ExpiresAt: expiresAt,
		},
	}
	// Use upsert to avoid duplicate key errors (only one row per competitionCode+dataType)
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "competition_code"}, {Name: "data_type"}},
		DoUpdates: clause.AssignmentColumns([]string{"data", "expires_at", "updated_at"}),
	}).Create(&cache).Error
}

// GetCachedTodayFixtures retrieves cached today's fixtures data if not expired
func (r *cacheRepository) GetCachedTodayFixtures() (*models.CachedTodayFixtures, error) {
	var cache models.CachedTodayFixtures

	err := r.db.Where("expires_at > ?", time.Now()).First(&cache).Error

	if err != nil {
		return nil, err // Return nil if not found or error
	}

	return &cache, nil
}

// GetCachedTodayFixturesIgnoringExpiry retrieves cached today's fixtures data even if expired
func (r *cacheRepository) GetCachedTodayFixturesIgnoringExpiry() (*models.CachedTodayFixtures, error) {
	var cache models.CachedTodayFixtures

	err := r.db.First(&cache).Error

	if err != nil {
		return nil, err // Return nil if not found or error
	}

	return &cache, nil
}

// StoreCachedFixtures stores fixtures data with an expiration time
func (r *cacheRepository) StoreCachedTodayFixtures(data models.JSONB, ttl time.Duration) error {
	expiresAt := time.Now().Add(ttl)
	cache := models.CachedTodayFixtures{
		Data: data,
		CachedData: models.CachedData{
			ExpiresAt: expiresAt,
		},
	}
	// Use upsert to avoid duplicate key errors (only one row should exist)
	return r.db.Clauses(clause.OnConflict{
		DoUpdates: clause.AssignmentColumns([]string{"data", "expires_at", "updated_at"}),
	}).Create(&cache).Error
}

// CleanExpiredCache removes all expired cache entries
func (r *cacheRepository) CleanExpiredCache() error {
	// Clean fixtures cache
	if err := r.db.Where("expires_at < ?", time.Now()).Delete(&models.CachedFixtures{}).Error; err != nil {
		return err
	}

	// Clean today fixtures cache
	if err := r.db.Where("expires_at < ?", time.Now()).Delete(&models.CachedTodayFixtures{}).Error; err != nil {
		return err
	}

	return nil
}

// Get retrieves a value from cache by key
func (r *cacheRepository) Get(key string) (*models.CacheItem, error) {
	var cacheItem models.CacheItem

	err := r.db.Where("key = ? AND expires_at > ?", key, time.Now()).First(&cacheItem).Error
	if err != nil {
		return nil, err
	}

	return &cacheItem, nil
}

// GetWithVersion gets a cache item only if its version matches
func (r *cacheRepository) GetWithVersion(key string, version string) (*models.CacheItem, error) {
	var item models.CacheItem
	err := r.db.Where("key = ? AND e_tag = ? AND expires_at > ?", key, version, time.Now()).First(&item).Error
	if err != nil {
		return nil, err
	}
	return &item, nil
}

// Set stores a value in cache with an expiration time
func (r *cacheRepository) Set(key string, value []byte, ttl time.Duration) error {
	cacheItem := models.CacheItem{
		Key:          key,
		Value:        value,
		ExpiresAt:    time.Now().Add(ttl),
		LastModified: time.Now(),
	}
	// Upsert the cache item (avoid duplicate key errors)
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "expires_at", "last_modified"}),
	}).Create(&cacheItem).Error
}

// SetWithMetadata stores a cache item with all metadata
func (r *cacheRepository) SetWithMetadata(key string, item models.CacheItem) error {
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "key"}},
		DoUpdates: clause.AssignmentColumns([]string{"value", "expires_at", "last_modified", "e_tag"}),
	}).Create(&item).Error
}

// UpdateVersion updates the version (ETag) of a cached item
func (r *cacheRepository) UpdateVersion(key string, version string) error {
	now := time.Now()
	result := r.db.Model(&models.CacheItem{}).
		Where("key = ?", key).
		Updates(map[string]interface{}{
			"e_tag":         version,
			"last_modified": now,
		})
	return result.Error
}
