package service

import (
	"context"
	"errors"
	"fmt"
	"libero-backend/config"               // Added config import
	"libero-backend/internal/models"      // Assuming user model is here
	"libero-backend/internal/repository" // Assuming repository interface is here
)

// UserService defines the interface for user data management.
type UserService interface {
	// Used by Auth/OAuth flows
	FindUserByEmail(ctx context.Context, email string) (*models.User, error)
	FindUserByProvider(ctx context.Context, provider string, providerID string) (*models.User, error)
	CreateUser(ctx context.Context, user *models.User) (*models.User, error)
	UpdateUser(ctx context.Context, user *models.User) error

	// Used by UserController
	RegisterUser(user *models.User) error                                     // Used in Register handler
	LoginUser(email, password string) (string, error)                         // Used in Login handler (returns token)
	GetUserByID(id uint) (*models.User, error)                                // Used in GetProfile, GetUser handlers
	// UpdateUser signature kept as is for now, controller call needs fixing
	ListUsers(page, limit int) ([]models.User, int64, error) // Used in ListUsers handler
	DeleteUser(id uint) error                                // Used in DeleteUser handler

	// Profile and Preferences
	GetUserProfile(ctx context.Context, userID uint) (*models.UserProfileResponse, error)
	UpdateUserPreferences(ctx context.Context, userID uint, req models.UpdatePreferencesRequest) error
}

// userService implements the UserService interface.
type userService struct {
	userRepo repository.UserRepository // Dependency on UserRepository
}

// NewUserService creates a new UserService instance.
func NewUserService(userRepo repository.UserRepository, cfg *config.Config) UserService { // Added cfg parameter
	return &userService{
		userRepo: userRepo,
	}
}

// FindUserByID finds a user by their internal ID.
// GetUserByID finds a user by their internal ID. (Renamed from FindUserByID)
func (s *userService) GetUserByID(id uint) (*models.User, error) {
	fmt.Printf("UserService: Finding user by ID: %d\n", id)
	// Placeholder: Replace with actual repository call
	// return s.userRepo.FindByID(ctx, id)
	if id == 1 { // Dummy user
		return &models.User{ID: 1, Email: "dummy@example.com", Name: "Dummy User"}, nil
	}
	return nil, errors.New("user not found (placeholder)")
}

// FindUserByEmail finds a user by their email address.
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	// Call the repository method
	return s.userRepo.FindByEmail(email)
}

// FindUserByProvider finds a user by their OAuth provider and provider-specific ID.
func (s *userService) FindUserByProvider(ctx context.Context, provider string, providerID string) (*models.User, error) {
	// Call the repository method
	return s.userRepo.FindByProvider(provider, providerID)
}

// CreateUser creates a new user record.
func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Call the repository method
	err := s.userRepo.Create(user)
	if err != nil {
		return nil, err // Return nil user on error
	}
	// Return the created user (which might now have an ID assigned by the DB)
	return user, nil
}

// UpdateUser updates an existing user record.
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	// Call the repository method
	return s.userRepo.Update(user)
}

// --- Placeholder implementations for new interface methods ---

func (s *userService) RegisterUser(user *models.User) error {
	fmt.Printf("UserService: Registering user (placeholder): Email=%s\n", user.Email)
	// Placeholder: Hash password, call CreateUser
	// _, err := s.CreateUser(context.Background(), user) // Example call
	// return err
	return nil
}

func (s *userService) LoginUser(email, password string) (string, error) {
	fmt.Printf("UserService: Logging in user (placeholder): Email=%s\n", email)
	// Placeholder: Find user by email, compare password, generate token
	// user, err := s.FindUserByEmail(context.Background(), email)
	// if err != nil { return "", err }
	// if !user.ComparePassword(password) { return "", errors.New("invalid credentials") }
	// token := "dummy-jwt-token-for-" + email // Replace with actual token generation
	// return token, nil
	if email == "existing@example.com" && password == "password" {
		return "dummy-jwt-token-for-" + email, nil
	}
	return "", errors.New("invalid credentials (placeholder)")
}

func (s *userService) ListUsers(page, limit int) ([]models.User, int64, error) {
	fmt.Printf("UserService: Listing users (placeholder): Page=%d, Limit=%d\n", page, limit)
	// Placeholder: Call repository List method
	// return s.userRepo.List(page, limit)
	dummyUsers := []models.User{
		{ID: 1, Email: "dummy@example.com", Name: "Dummy User"},
		{ID: 2, Email: "existing@example.com", Name: "Existing User"},
	}
	return dummyUsers, int64(len(dummyUsers)), nil
}

func (s *userService) DeleteUser(id uint) error {
	fmt.Printf("UserService: Deleting user (placeholder): ID=%d\n", id)
	// Placeholder: Call repository Delete method
	// return s.userRepo.Delete(id)
	return nil
}

// --- Profile and Preference Methods ---

// GetUserProfile retrieves the user's profile including followed teams and players.
func (s *userService) GetUserProfile(ctx context.Context, userID uint) (*models.UserProfileResponse, error) {
	user, err := s.userRepo.FindByIDWithPreferences(userID)
	if err != nil {
		// Handle specific errors like gorm.ErrRecordNotFound if needed
		return nil, fmt.Errorf("failed to find user profile: %w", err)
	}

	// Map User model to UserProfileResponse DTO
	profileResponse := &models.UserProfileResponse{
		ID:    user.ID,
		Name:  user.Name, // Assuming Name is populated (e.g., from OAuth or profile settings)
		Email: user.Email,
		Preferences: models.UserPreferencesResponse{
			FollowedTeams:        make([]models.TeamPreferenceInfo, 0, len(user.FollowedTeams)),
			FollowedPlayers:      make([]models.PlayerPreferenceInfo, 0, len(user.FollowedPlayers)),
			FollowedCompetitions: make([]models.CompetitionPreferenceInfo, 0, len(user.FollowedCompetitions)),
		},
	}

	for _, team := range user.FollowedTeams {
		profileResponse.Preferences.FollowedTeams = append(profileResponse.Preferences.FollowedTeams, models.TeamPreferenceInfo{
			ID:   team.ID,
			Name: team.Name,
		})
	}

	for _, player := range user.FollowedPlayers {
		profileResponse.Preferences.FollowedPlayers = append(profileResponse.Preferences.FollowedPlayers, models.PlayerPreferenceInfo{
			ID:   player.ID,
			Name: player.Name,
		})
	}

	for _, competition := range user.FollowedCompetitions {
		profileResponse.Preferences.FollowedCompetitions = append(profileResponse.Preferences.FollowedCompetitions, models.CompetitionPreferenceInfo{
			ID:   competition.ID,
			Name: competition.Name,
		})
	}

	return profileResponse, nil
}

// UpdateUserPreferences updates the user's followed teams and players based on the request.
func (s *userService) UpdateUserPreferences(ctx context.Context, userID uint, req models.UpdatePreferencesRequest) error {
	// Process removals first
	for _, teamID := range req.RemoveTeams {
		err := s.userRepo.RemoveFollowedTeam(userID, teamID)
		if err != nil {
			// Log error but potentially continue? Or return immediately?
			// For now, let's log and continue to make it more robust against partial failures.
			fmt.Printf("WARN: Failed to remove followed team %d for user %d: %v\n", teamID, userID, err)
			// Consider collecting errors and returning them all at the end.
		}
	}
	for _, playerID := range req.RemovePlayers {
		err := s.userRepo.RemoveFollowedPlayer(userID, playerID)
		if err != nil {
			fmt.Printf("WARN: Failed to remove followed player %d for user %d: %v\n", playerID, userID, err)
		}
	}
	for _, competitionID := range req.RemoveCompetitions {
		err := s.userRepo.RemoveFollowedCompetition(userID, competitionID)
		if err != nil {
			fmt.Printf("WARN: Failed to remove followed competition %d for user %d: %v\n", competitionID, userID, err)
		}
	}

	// Process additions
	for _, teamID := range req.AddTeams {
		err := s.userRepo.AddFollowedTeam(userID, teamID)
		if err != nil {
			// Handle errors like team not found, user not found, or DB constraint errors
			fmt.Printf("WARN: Failed to add followed team %d for user %d: %v\n", teamID, userID, err)
			// If a team/player doesn't exist, should we create it here?
			// The plan suggests FindTeamByNameOrID in repo, implying they should exist.
			// Let's assume they must exist for now.
		}
	}
	for _, playerID := range req.AddPlayers {
		err := s.userRepo.AddFollowedPlayer(userID, playerID)
		if err != nil {
			fmt.Printf("WARN: Failed to add followed player %d for user %d: %v\n", playerID, userID, err)
		}
	}
	for _, competitionID := range req.AddCompetitions {
		err := s.userRepo.AddFollowedCompetition(userID, competitionID)
		if err != nil {
			fmt.Printf("WARN: Failed to add followed competition %d for user %d: %v\n", competitionID, userID, err)
		}
	}

	// Currently returns nil even if some operations failed (logged as warnings).
	// Consider returning a multi-error if stricter error handling is needed.
	return nil
}