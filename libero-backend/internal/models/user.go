package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user data model
type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Email     string    `gorm:"uniqueIndex;not null" json:"email"`
	Username  string    `gorm:"uniqueIndex;not null" json:"username"`
	// Name      string    `json:"name,omitempty"` // Comment out or remove old Name field
	FirstName string    `json:"first_name,omitempty" gorm:"column:firstname"` // Added FirstName, map to DB column
	LastName  string    `json:"last_name,omitempty"  gorm:"column:lastname"`  // Added LastName, map to DB column
	Password  string    `gorm:"not null" json:"-"`    // Password is not exposed in JSON
	Role      string    `gorm:"default:user" json:"role"` // user, admin, etc.
	Active    bool      `gorm:"default:true" json:"active"`
	Provider  string    `gorm:"index" json:"-"` // OAuth Provider (e.g., google), indexed
	ProviderID string   `gorm:"index" json:"-"` // User ID from the OAuth provider, indexed
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeSave is a GORM hook to hash the password before saving
func (u *User) BeforeSave(tx *gorm.DB) error {
	if u.Password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		u.Password = string(hashedPassword)
	}
	return nil
}

// ComparePassword checks if the provided password matches the stored hashed password
func (u *User) ComparePassword(password string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password))
	return err == nil
}

// UserResponse is used to return user data without sensitive information
type UserResponse struct {
	ID        uint      `json:"id"`
	Username  string    `json:"username"`
	Email     string    `json:"email"`
	FirstName string    `json:"first_name,omitempty"`
	LastName  string    `json:"last_name,omitempty"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ToResponse converts a User to UserResponse, omitting sensitive fields
func (u *User) ToResponse() UserResponse {
	return UserResponse{
		ID:        u.ID,
		Username:  u.Username,
		Email:     u.Email,
		FirstName: u.FirstName,
		LastName:  u.LastName,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}