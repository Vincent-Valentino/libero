package models

import (
    "time"
)

// CacheItem represents a generic cache entry
type CacheItem struct {
    Key       string    `gorm:"primaryKey"`
    Value     []byte    `gorm:"type:bytea"`
    ETag      string    `gorm:"index"` // For tracking data version
    ExpiresAt time.Time `gorm:"index"`
    CreatedAt time.Time `gorm:"autoCreateTime"`
    LastModified time.Time // Track when the actual data was last modified
}
