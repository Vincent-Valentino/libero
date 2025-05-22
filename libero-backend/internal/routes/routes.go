package routes

import (
	"net/http"

	"libero-backend/config" // Added config dependency

	"github.com/gorilla/mux"

	// "libero-backend/config" // Removed if not needed directly
	"libero-backend/internal/controllers"
	"libero-backend/internal/middleware"
	"libero-backend/internal/service" // Import service package
)

// SetupRoutes configures all API routes
// It now accepts the main service struct to access all services
func SetupRoutes(router *mux.Router, service *service.Service, cfg *config.Config) { // Added cfg
	// Extract controllers and services needed
	// Instantiate the main controller which holds sub-controllers
	ctrl := controllers.New(service, cfg)
	authService := service.Auth // Get AuthService for middleware

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Auth routes
	api.HandleFunc("/auth/register", ctrl.User.Register).Methods(http.MethodPost) // Changed path prefix
	api.HandleFunc("/auth/login", ctrl.User.Login).Methods(http.MethodPost)       // Changed path prefix
	api.HandleFunc("/auth/password/reset-request", ctrl.User.RequestPasswordReset).Methods(http.MethodPost)
	api.HandleFunc("/auth/password/reset", ctrl.User.ResetPassword).Methods(http.MethodPost)

	// NEW: Public Sports Data routes
	api.HandleFunc("/matches/upcoming", ctrl.SportsData.HandleGetUpcomingMatches).Methods(http.MethodGet)
	api.HandleFunc("/matches/results", ctrl.SportsData.HandleGetResults).Methods(http.MethodGet)
	api.HandleFunc("/players/{player_id}/stats", ctrl.SportsData.HandleGetPlayerStats).Methods(http.MethodGet)
	api.HandleFunc("/sports/fixtures/today", ctrl.SportsData.HandleGetTodaysFixtures).Methods(http.MethodGet)    // Today's fixtures
	api.HandleFunc("/sports/fixtures/summary", ctrl.SportsData.HandleGetFixturesSummary).Methods(http.MethodGet) // Fixtures summary per competition

	// OAuth routes (public) - Using root router for /auth path
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/google/login", ctrl.Oauth.GoogleLogin).Methods(http.MethodGet)
	auth.HandleFunc("/google/callback", ctrl.Oauth.GoogleCallback).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/login", ctrl.Oauth.FacebookLogin).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/callback", ctrl.Oauth.FacebookCallback).Methods(http.MethodGet)
	auth.HandleFunc("/github/login", ctrl.Oauth.GitHubLogin).Methods(http.MethodGet)
	auth.HandleFunc("/github/callback", ctrl.Oauth.GitHubCallback).Methods(http.MethodGet)

	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService)) // Inject AuthService

	// User routes (Profile & Preferences)
	protected.HandleFunc("/users/profile", ctrl.User.GetUserProfile).Methods(http.MethodGet)            // Get profile with preferences
	protected.HandleFunc("/users/preferences", ctrl.User.UpdateUserPreferences).Methods(http.MethodPut) // Update preferences
	// Password change route (requires authentication)
	protected.HandleFunc("/auth/password/change", ctrl.User.ChangePassword).Methods(http.MethodPost)
}
