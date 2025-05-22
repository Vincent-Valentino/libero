package models

import "time"

// Team represents a sports team
type Team struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	LogoURL   string    `json:"logo_url,omitempty"`
	Country   string    `json:"country,omitempty"`
	Sport     string    `json:"sport,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	// Relationships
	FollowedByUsers []*User `gorm:"many2many:user_followed_teams;"` // Users who follow this team
}