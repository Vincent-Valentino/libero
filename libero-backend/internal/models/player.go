package models

import "time"

// Player represents a sports player.
type Player struct {
	ID   uint   `gorm:"primaryKey"`
	Name string `gorm:"not null"`
	// TODO: Add other relevant fields like position, external_api_id etc. based on future needs.

	TeamID *uint // Optional foreign key to Team
	Team   Team  `gorm:"foreignKey:TeamID"` // Belongs To relationship

	CreatedAt time.Time
	UpdatedAt time.Time

	// Relationships
	FollowedByUsers []*User `gorm:"many2many:user_followed_players;"` // Users who follow this player
}