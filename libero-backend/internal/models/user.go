package models

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// User represents the user data model
type User struct {
	ID         uint   `gorm:"primaryKey" json:"id"`
	Email      string `gorm:"uniqueIndex;not null" json:"email"`
	Username   string `gorm:"uniqueIndex" json:"username"` // Removed not null for OAuth users
	Name       string `json:"name,omitempty"`              // Display name, potentially from OAuth
	Password   string `json:"-"`                           // Removed not null for OAuth users, Password is not exposed in JSON
	Role       string `gorm:"default:user" json:"role"`    // user, admin, etc.
	Active     bool   `gorm:"default:true" json:"active"`
	Provider   string `gorm:"index" json:"-"` // OAuth Provider (e.g., google), indexed
	ProviderID string `gorm:"index" json:"-"` // User ID from the OAuth provider, indexed

	// Password reset fields
	ResetToken          string    `gorm:"index" json:"-"` // Password reset token
	ResetTokenExpiresAt time.Time `json:"-"`              // When the reset token expires

	// Relationships for Preferences
	FollowedTeams        []*Team        `gorm:"many2many:user_followed_teams;" json:"followed_teams,omitempty"`
	FollowedPlayers      []*Player      `gorm:"many2many:user_followed_players;" json:"followed_players,omitempty"`
	FollowedCompetitions []*Competition `gorm:"many2many:user_followed_competitions;" json:"followed_competitions,omitempty"`

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
	Name      string    `json:"name,omitempty"`
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
		Name:      u.Name,
		Role:      u.Role,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

// --- DTOs for Profile Feature ---

// UpdatePreferencesRequest defines the structure for the PUT /api/users/preferences request body.
type UpdatePreferencesRequest struct {
	AddTeams           []uint `json:"add_teams"`           // IDs of teams to follow
	RemoveTeams        []uint `json:"remove_teams"`        // IDs of teams to unfollow
	AddPlayers         []uint `json:"add_players"`         // IDs of players to follow
	RemovePlayers      []uint `json:"remove_players"`      // IDs of players to unfollow
	AddCompetitions    []uint `json:"add_competitions"`    // IDs of competitions to follow
	RemoveCompetitions []uint `json:"remove_competitions"` // IDs of competitions to unfollow
}

// TeamPreferenceInfo defines the structure for team info within the profile response.
type TeamPreferenceInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}

// PlayerPreferenceInfo defines the structure for player info within the profile response.
type PlayerPreferenceInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Add TeamName or TeamID if needed in the future
}

// CompetitionPreferenceInfo defines the structure for competition info within the profile response.
type CompetitionPreferenceInfo struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
	// Add Sport or Region if needed
}

// UserPreferencesResponse defines the nested structure for preferences in the profile response.
type UserPreferencesResponse struct {
	FollowedTeams        []TeamPreferenceInfo        `json:"followed_teams"`
	FollowedPlayers      []PlayerPreferenceInfo      `json:"followed_players"`
	FollowedCompetitions []CompetitionPreferenceInfo `json:"followed_competitions"`
}

// UserProfileResponse defines the structure for the GET /api/users/profile response body.
type UserProfileResponse struct {
	ID          uint                    `json:"id"`
	Name        string                  `json:"name"`
	Email       string                  `json:"email"`
	Preferences UserPreferencesResponse `json:"preferences"`
}
