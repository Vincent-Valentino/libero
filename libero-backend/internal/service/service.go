package service

import (
	"libero-backend/config"
	"libero-backend/internal/repository"
)

// Service provides access to all service operations
type Service struct {
	User  UserService
  // Keep as interface type
	Auth  AuthService  // Add AuthService field (interface type)
	OAuth OAuthService // Add OAuthService field (interface type)
	// Add other services here
}

// New creates a new service instance with all services
func New(repo *repository.Repository) *Service {
	// Get configuration
	cfg := config.New()
 // Config needed by multiple services

	// Initialize services in dependency order
	userService := NewUserService(repo.User, cfg)
	authService := NewAuthService(userService, cfg.JWT) // AuthService depends on UserService
	oauthService := NewOAuthService(cfg, authService) // OAuthService depends on Config and AuthService

	return &Service{
		User:  userService,
		Auth:  authService,
		OAuth: oauthService,
		// Initialize other services
	}
}