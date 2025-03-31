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
	// Add FindByProvider method if needed by UserService implementation
	FindByProvider(provider string, providerID string) (*models.User, error)
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