package repository

import "gorm.io/gorm"

// Repository provides access to all repository operations
type Repository struct {
	User  UserRepository
	Cache CacheRepository
	// Add more repositories here as needed
}

// New creates a new repository instance with all repositories
func New(db *gorm.DB) *Repository {
	return &Repository{
		User:  NewUserRepository(db),
		Cache: NewCacheRepository(db),
		// Initialize other repositories here
	}
}