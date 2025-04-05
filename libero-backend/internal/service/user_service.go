package service

import (
	"context"
	"errors"
	"fmt"
	"libero-backend/config"              // Added config import
	"libero-backend/internal/models"     // Assuming user model is here
	"libero-backend/internal/repository" // Assuming repository interface is here
	"gorm.io/gorm"                       // <-- Added GORM import
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
	// LoginUser method removed from interface (handled by AuthService)
	GetUserByID(id uint) (*models.User, error)                                // Used in GetProfile, GetUser handlers
	// UpdateUser signature kept as is for now, controller call needs fixing
	ListUsers(page, limit int) ([]models.User, int64, error)                  // Used in ListUsers handler
	DeleteUser(id uint) error                                                 // Used in DeleteUser handler
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
	// Call the repository to find the user by ID
	// Note: Repository FindByID doesn't currently use context
	user, err := s.userRepo.FindByID(id)
	if err != nil {
		// Check if the error is specifically GORM's record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Return our canonical not found error
		}
		// TODO: Add structured logging for other DB errors
		fmt.Printf("ERROR: Database error finding user by ID '%d': %v\n", id, err)
		return nil, fmt.Errorf("database error finding user: %w", err)
	}
	return user, nil
}

// FindUserByEmail finds a user by their email address.
func (s *userService) FindUserByEmail(ctx context.Context, email string) (*models.User, error) {
	fmt.Printf("UserService: Finding user by Email: %s\n", email)
	// Call the repository to find the user by email
	// Note: The repository FindByEmail method doesn't use context currently.
	user, err := s.userRepo.FindByEmail(email)
	if err != nil {
		// Check if the error is specifically GORM's record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Return our canonical not found error
		}
		// TODO: Add structured logging for other DB errors
		fmt.Printf("ERROR: Database error finding user by email '%s': %v\n", email, err)
		return nil, fmt.Errorf("database error finding user: %w", err)
	}
	return user, nil
}

// FindUserByProvider finds a user by their OAuth provider and provider-specific ID.
func (s *userService) FindUserByProvider(ctx context.Context, provider string, providerID string) (*models.User, error) {
	fmt.Printf("UserService: Finding user by Provider: %s, ProviderID: %s\n", provider, providerID)
	// Call the repository to find the user by provider details
	user, err := s.userRepo.FindByProvider(provider, providerID)
	if err != nil {
		// Check if the error is specifically GORM's record not found
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrUserNotFound // Return our canonical not found error
		}
		// TODO: Add structured logging for other DB errors
		fmt.Printf("ERROR: Database error finding user by provider '%s'/'%s': %v\n", provider, providerID, err)
		return nil, fmt.Errorf("database error finding user by provider: %w", err)
	}
	return user, nil
}

// CreateUser creates a new user record.
func (s *userService) CreateUser(ctx context.Context, user *models.User) (*models.User, error) {
	// Note: The repository Create method doesn't use context currently.
	// Also, Username might be relevant for logging if available.
	fmt.Printf("UserService: Creating user: Email=%s, Username=%s, FirstName=%s, LastName=%s, Provider=%s\n", user.Email, user.Username, user.FirstName, user.LastName, user.Provider) // Updated logging
	// Call the repository to create the user
	err := s.userRepo.Create(user) // Call repo Create with just user
	if err != nil {
		return nil, err // Return error if creation failed
	}
	// If successful, GORM might have updated the user object with ID. Return the input user.
	return user, nil
}

// UpdateUser updates an existing user record.
func (s *userService) UpdateUser(ctx context.Context, user *models.User) error {
	fmt.Printf("UserService: Updating user: ID=%d\n", user.ID)
	// Call the repository to update the user
	err := s.userRepo.Update(user)
	if err != nil {
		// TODO: Add structured logging
		fmt.Printf("ERROR: Failed to update user ID '%d': %v\n", user.ID, err)
		return fmt.Errorf("failed to update user: %w", err)
	}
	return nil
}

// RegisterUser handles the registration process for a new user.
// It's called by the UserController's Register handler.
func (s *userService) RegisterUser(user *models.User) error {
	fmt.Printf("UserService: Registering user: Email=%s\n", user.Email)
	// Note: Password hashing happens in BeforeSave hook in models/user.go

	// Call CreateUser (which now calls the repository's Create method)
	// We need a context here. Use context.Background() for now,
	// but ideally, context should be passed down from the controller if needed by repo.
	_, err := s.CreateUser(context.Background(), user)
	if err != nil {
		// TODO: Add structured logging
		fmt.Printf("ERROR: Failed to register user (CreateUser failed): %v\n", err)
		// Consider returning more specific errors (e.g., ErrEmailExists if repo returns constraint violation)
		return fmt.Errorf("registration failed: %w", err)
	}
	return nil
}

// --- Other Placeholder implementations ---

// LoginUser function removed as login logic is handled by AuthService.LoginByPassword

func (s *userService) ListUsers(page, limit int) ([]models.User, int64, error) {
	fmt.Printf("UserService: Listing users: Page=%d, Limit=%d\n", page, limit)
	// Call the repository to list users with pagination
	return s.userRepo.List(page, limit) // Directly return the repository result
	// Error handling should ideally happen in the repository or be more specific here if needed.
	// For now, assuming the repository List handles errors appropriately or returns them.
}

func (s *userService) DeleteUser(id uint) error {
	fmt.Printf("UserService: Deleting user: ID=%d\n", id)
	// Call the repository to delete the user
	return s.userRepo.Delete(id) // Directly return the repository result
	// Error handling should ideally happen in the repository or be more specific here if needed.
}

// Add other UserService method implementations here...