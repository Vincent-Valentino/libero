package models

import "time"

// Team represents a sports team.
type Team struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"unique;not null"`
	// TODO: Add other relevant fields like sport, league, external_api_id etc. based on future needs.
	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	FollowedByUsers []*User `gorm:"many2many:user_followed_teams;"` // Users who follow this team
}