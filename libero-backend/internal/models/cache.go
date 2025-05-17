package models

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
	"time"

	"gorm.io/gorm"
)

// CachedData represents base fields for all cached data
type CachedData struct {
	gorm.Model
	ExpiresAt time.Time `gorm:"index"` // When this cache entry expires
}

// JSONB is a custom type for handling JSON data in PostgreSQL
type JSONB map[string]interface{}

// Value make the JSONB struct implement the driver.Valuer interface
func (j JSONB) Value() (driver.Value, error) {
	return json.Marshal(j)
}

// Scan make the JSONB struct implement the sql.Scanner interface
func (j *JSONB) Scan(value interface{}) error {
	bytes, ok := value.([]byte)
	if !ok {
		return errors.New("type assertion to []byte failed")
	}
	
	return json.Unmarshal(bytes, &j)
}

// CachedFixtures represents cached fixture data for a competition
type CachedFixtures struct {
	CachedData
	CompetitionCode string `gorm:"uniqueIndex"` // Code like "PL", "EL", etc.
	DataType        string `gorm:"not null"`    // "today", "fixtures_summary"
	Data            JSONB  `gorm:"type:jsonb"`  // The actual JSON data
}

// CachedTodayFixtures represents cached data for today's fixtures across all competitions
type CachedTodayFixtures struct {
	CachedData
	Data JSONB `gorm:"type:jsonb"` // The actual JSON data
} 