package models

import "time"

// Competition represents a sports competition or league
type Competition struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Name      string    `gorm:"not null" json:"name"`
	Code      string    `gorm:"uniqueIndex" json:"code,omitempty"`
	Country   string    `json:"country,omitempty"`
	LogoURL   string    `json:"logo_url,omitempty"`
	Sport     string    `json:"sport,omitempty"`
	Season    string    `json:"season,omitempty"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}