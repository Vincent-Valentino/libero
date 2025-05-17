package repository

import (
	"libero-backend/internal/models"
	"time"

	"gorm.io/gorm"
)

// CacheRepository defines the interface for cache data operations
type CacheRepository interface {
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
	// Try to find existing record first
	var existing models.CachedFixtures
	result := r.db.Where("competition_code = ? AND data_type = ?", competitionCode, dataType).First(&existing)
	
	expiresAt := time.Now().Add(ttl)
	
	// Update existing or create new
	if result.Error == nil {
		// Update existing
		existing.Data = data
		existing.ExpiresAt = expiresAt
		return r.db.Save(&existing).Error
	} else {
		// Create new
		cache := models.CachedFixtures{
			CompetitionCode: competitionCode,
			DataType:        dataType,
			Data:            data,
			CachedData: models.CachedData{
				ExpiresAt: expiresAt,
			},
		}
		return r.db.Create(&cache).Error
	}
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

// StoreCachedTodayFixtures stores today's fixtures data with an expiration time
func (r *cacheRepository) StoreCachedTodayFixtures(data models.JSONB, ttl time.Duration) error {
	// Try to find existing record first
	var existing models.CachedTodayFixtures
	result := r.db.First(&existing)
	
	expiresAt := time.Now().Add(ttl)
	
	// Update existing or create new
	if result.Error == nil {
		// Update existing
		existing.Data = data
		existing.ExpiresAt = expiresAt
		return r.db.Save(&existing).Error
	} else {
		// Create new
		cache := models.CachedTodayFixtures{
			Data: data,
			CachedData: models.CachedData{
				ExpiresAt: expiresAt,
			},
		}
		return r.db.Create(&cache).Error
	}
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