package service

import (
	"libero-backend/config"
	"libero-backend/internal/repository"
)

// Service provides access to all service operations
type Service struct {
	User     UserService
	Auth     AuthService
	OAuth    OAuthService
	ML       MLService
	Fixtures FixturesService
	Football *FootballService // Add Football service
}

// New creates a new service instance with all services
func New(repo *repository.Repository) *Service {
	// Load configuration (including ThirdPartyAPIKey & BaseURL)
	cfg := config.New()

	// Initialize services in dependency order
	userService := NewUserService(repo.User, cfg)
	authService := NewAuthService(userService, cfg.JWT) // AuthService depends on UserService
	oauthService := NewOAuthService(cfg, authService)   // OAuthService depends on Config and AuthService
	mlService := NewMLService(cfg)                      // MLService depends on Config
	fixturesService := NewFixturesService(cfg.ThirdPartyAPIKey, cfg.ThirdPartyBaseURL, repo.Cache)
	footballService := NewFootballService(cfg.ThirdPartyBaseURL, cfg.ThirdPartyAPIKey) // Initialize with API config

	return &Service{
		User:     userService,
		Auth:     authService,
		OAuth:    oauthService,
		ML:       mlService,
		Fixtures: fixturesService,
		Football: footballService, // Add to returned service
	}
}
