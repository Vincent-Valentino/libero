package repository

import (
	"libero-backend/internal/models"
	"gorm.io/gorm"
)

// UserRepository defines the interface for user data operations.
type UserRepository interface {
	Create(user *models.User) error
	FindByID(id uint) (*models.User, error)
	FindByEmail(email string) (*models.User, error)
	FindByUsername(username string) (*models.User, error)
	Update(user *models.User) error
	Delete(id uint) error
	List(page, limit int) ([]models.User, int64, error)
	FindByProvider(provider string, providerID string) (*models.User, error)
	FindByResetToken(token string) (*models.User, error)

	// Preference related methods
	FindByIDWithPreferences(id uint) (*models.User, error)
	AddFollowedTeam(userID uint, teamID uint) error
	RemoveFollowedTeam(userID uint, teamID uint) error
	AddFollowedPlayer(userID uint, playerID uint) error
	RemoveFollowedPlayer(userID uint, playerID uint) error
	AddFollowedCompetition(userID uint, competitionID uint) error
	RemoveFollowedCompetition(userID uint, competitionID uint) error
}

// userRepository implements the UserRepository interface.
type userRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a new user repository instance
func NewUserRepository(db *gorm.DB) UserRepository {
 // Return interface type
	return &userRepository{db: db}
 // Return pointer to concrete struct
}

// Create adds a new user to the database
func (r *userRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

// FindByID retrieves a user by ID
func (r *userRepository) FindByID(id uint) (*models.User, error) {
	var user models.User
	err := r.db.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByEmail retrieves a user by email
func (r *userRepository) FindByEmail(email string) (*models.User, error) {
	var user models.User
	err := r.db.Where("email = ?", email).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// FindByUsername retrieves a user by username
func (r *userRepository) FindByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.db.Where("username = ?", username).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// Update updates user information
func (r *userRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// Delete removes a user from the database
func (r *userRepository) Delete(id uint) error {
	return r.db.Delete(&models.User{}, id).Error
}

// List retrieves all users with pagination
func (r *userRepository) List(page, limit int) ([]models.User, int64, error) {
	var users []models.User
	var count int64

	offset := (page - 1) * limit

	// Get total count
	if err := r.db.Model(&models.User{}).Count(&count).Error; err != nil {
		return nil, 0, err
	}

	// Get users for the current page
	if err := r.db.Offset(offset).Limit(limit).Find(&users).Error; err != nil {
		return nil, 0, err
	}

	return users, count, nil
}

// FindByProvider retrieves a user by OAuth provider and provider ID
// Note: This needs to be implemented based on how provider info is stored.
// Assuming 'provider' and 'provider_id' columns exist in the 'users' table.
func (r *userRepository) FindByProvider(provider string, providerID string) (*models.User, error) {
	var user models.User
	err := r.db.Where("provider = ? AND provider_id = ?", provider, providerID).First(&user).Error
	if err != nil {
		// Consider returning a specific "not found" error type here
		return nil, err
	}
	return &user, nil
}

// FindByResetToken retrieves a user by password reset token
func (r *userRepository) FindByResetToken(token string) (*models.User, error) {
	var user models.User
	err := r.db.Where("reset_token = ?", token).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// --- Preference Methods ---

// FindByIDWithPreferences retrieves a user by ID, preloading followed teams and players.
func (r *userRepository) FindByIDWithPreferences(id uint) (*models.User, error) {
	var user models.User
	err := r.db.Preload("FollowedTeams").Preload("FollowedPlayers").Preload("FollowedCompetitions").First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

// findTeamByID is an internal helper to find a team by ID.
func (r *userRepository) findTeamByID(teamID uint) (*models.Team, error) {
	var team models.Team
	if err := r.db.First(&team, teamID).Error; err != nil {
		return nil, err // Return error to handle not found case upstream
	}
	return &team, nil
}

// findPlayerByID is an internal helper to find a player by ID.
func (r *userRepository) findPlayerByID(playerID uint) (*models.Player, error) {
	var player models.Player
	if err := r.db.First(&player, playerID).Error; err != nil {
		return nil, err // Return error to handle not found case upstream
	}
	return &player, nil
}

// findCompetitionByID is an internal helper to find a competition by ID.
func (r *userRepository) findCompetitionByID(competitionID uint) (*models.Competition, error) {
	var competition models.Competition
	if err := r.db.First(&competition, competitionID).Error; err != nil {
		return nil, err // Return error to handle not found case upstream
	}
	return &competition, nil
}

// AddFollowedTeam adds a team to a user's followed list.
func (r *userRepository) AddFollowedTeam(userID uint, teamID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	team, err := r.findTeamByID(teamID)
	if err != nil {
		return err // Team not found
	}

	// Use GORM's Association Append method
	return r.db.Model(user).Association("FollowedTeams").Append(team)
}

// RemoveFollowedTeam removes a team from a user's followed list.
func (r *userRepository) RemoveFollowedTeam(userID uint, teamID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	team, err := r.findTeamByID(teamID)
	if err != nil {
		// If team doesn't exist, we can't remove it anyway, maybe log this?
		// For now, return nil as the association likely doesn't exist.
		// Or return err if strict checking is needed. Let's return nil for idempotency.
		return nil
	}

	// Use GORM's Association Delete method
	return r.db.Model(user).Association("FollowedTeams").Delete(team)
}

// AddFollowedPlayer adds a player to a user's followed list.
func (r *userRepository) AddFollowedPlayer(userID uint, playerID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	player, err := r.findPlayerByID(playerID)
	if err != nil {
		return err // Player not found
	}

	// Use GORM's Association Append method
	return r.db.Model(user).Association("FollowedPlayers").Append(player)
}

// RemoveFollowedPlayer removes a player from a user's followed list.
func (r *userRepository) RemoveFollowedPlayer(userID uint, playerID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	player, err := r.findPlayerByID(playerID)
	if err != nil {
		// Similar to RemoveFollowedTeam, return nil for idempotency.
		return nil
	}

	// Use GORM's Association Delete method
	return r.db.Model(user).Association("FollowedPlayers").Delete(player)
}

// AddFollowedCompetition adds a competition to a user's followed list.
func (r *userRepository) AddFollowedCompetition(userID uint, competitionID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	competition, err := r.findCompetitionByID(competitionID)
	if err != nil {
		return err // Competition not found
	}

	// Use GORM's Association Append method
	return r.db.Model(user).Association("FollowedCompetitions").Append(competition)
}

// RemoveFollowedCompetition removes a competition from a user's followed list.
func (r *userRepository) RemoveFollowedCompetition(userID uint, competitionID uint) error {
	user, err := r.FindByID(userID)
	if err != nil {
		return err // User not found
	}
	competition, err := r.findCompetitionByID(competitionID)
	if err != nil {
		// Similar to RemoveFollowedTeam, return nil for idempotency.
		return nil
	}

	// Use GORM's Association Delete method
	return r.db.Model(user).Association("FollowedCompetitions").Delete(competition)
}