package models

import "time"

// Player represents a sports player
type Player struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Position  string    `json:"position,omitempty"`
	TeamID    uint      `json:"team_id,omitempty"`
	Country   string    `json:"country,omitempty"`
	PhotoURL  string    `json:"photo_url,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}