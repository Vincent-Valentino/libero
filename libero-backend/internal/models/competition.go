package models

import "gorm.io/gorm"

// Competition represents a sports competition (e.g., Premier League, NBA).
type Competition struct {
	gorm.Model
	ID            uint   `gorm:"primaryKey"`
	Name          string `gorm:"uniqueIndex;not null"` // Name of the competition
	Sport         string // e.g., "Football", "Basketball"
	Region        string // e.g., "England", "USA", "International"
	ExternalApiID string `gorm:"index"` // Optional: ID from an external data source
	Users         []*User `gorm:"many2many:user_followed_competitions;"` // Users following this competition
}