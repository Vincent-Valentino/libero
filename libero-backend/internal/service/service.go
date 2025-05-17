package service

import (
	"libero-backend/config"
	"libero-backend/internal/repository"
)

// Service provides access to all service operations
type Service struct {
	User  UserService
  // Keep as interface type
	Auth  AuthService   // Add AuthService field (interface type)
	OAuth OAuthService  // Add OAuthService field (interface type)
	ML    MLService     // Changed from pointer to interface type
	Fixtures FixturesService // Add FixturesService field (interface type)
}

// New creates a new service instance with all services
func New(repo *repository.Repository) *Service {
	// Load configuration (including ThirdPartyAPIKey & BaseURL)
	cfg := config.New()
	
	// Initialize services in dependency order
	userService := NewUserService(repo.User, cfg)
	authService := NewAuthService(userService, cfg.JWT) // AuthService depends on UserService
	oauthService := NewOAuthService(cfg, authService) // OAuthService depends on Config and AuthService
	mlService := NewMLService(cfg)         // MLService depends on Config
	fixturesService := NewFixturesService(cfg.ThirdPartyAPIKey, cfg.ThirdPartyBaseURL, repo.Cache)

	return &Service{
		User:     userService,
		Auth:     authService,
		OAuth:    oauthService,
		ML:       mlService, // Initialize MLService
		Fixtures: fixturesService, // Initialize FixturesService
	}
}