package routes

import (
	"encoding/json"
	"libero-backend/config"
	"libero-backend/internal/controllers"
	"libero-backend/internal/middleware"
	"libero-backend/internal/repository"
	"libero-backend/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

// healthCheck is a simple health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]bool{"ok": true})
}

// SetupRoutes configures all API routes
func SetupRoutes(router *mux.Router, service *service.Service, cfg *config.Config, repo *repository.Repository) {
	// Extract controllers and services needed
	ctrl := controllers.New(service, cfg, repo)
	authService := service.Auth // Get AuthService for middleware

	// Add global middlewares in correct order
	// CORS must be first to handle preflight requests properly
	router.Use(middleware.CORSMiddleware)
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.JSONMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()
	// Ensure subrouter inherits middleware from parent
	api.Use(middleware.CORSMiddleware)

	// Public routes (no authentication required)
	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/auth/register", ctrl.User.Register).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/auth/login", ctrl.User.Login).Methods(http.MethodPost, http.MethodOptions)

	// Sports data routes
	api.HandleFunc("/sports/fixtures/today", ctrl.SportsData.HandleGetTodaysFixtures).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/sports/fixtures/summary", ctrl.SportsData.HandleGetFixturesSummary).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/standings", ctrl.SportsData.HandleGetStandings).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/topscorers", ctrl.SportsData.HandleGetTopScorers).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/matches/upcoming", ctrl.SportsData.HandleGetUpcomingMatches).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/matches/results", ctrl.SportsData.HandleGetResults).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/players/{player_id}/stats", ctrl.SportsData.HandleGetPlayerStats).Methods(http.MethodGet, http.MethodOptions)

	// Match prediction routes
	api.HandleFunc("/predict/match", ctrl.Prediction.PredictMatch).Methods(http.MethodPost, http.MethodOptions)
	api.HandleFunc("/predict/teams", ctrl.Prediction.GetAvailableTeams).Methods(http.MethodGet, http.MethodOptions)
	api.HandleFunc("/predict/leagues", ctrl.Prediction.GetAvailableLeagues).Methods(http.MethodGet, http.MethodOptions)

	// OAuth routes - create subrouter and explicitly apply CORS middleware
	auth := router.PathPrefix("/auth").Subrouter()
	auth.Use(middleware.CORSMiddleware) // Explicitly apply CORS to OAuth subrouter

	// OAuth routes with explicit OPTIONS handling for preflight requests
	auth.HandleFunc("/google/login", ctrl.Oauth.GoogleLogin).Methods(http.MethodGet, http.MethodOptions)
	auth.HandleFunc("/google/callback", ctrl.Oauth.GoogleCallback).Methods(http.MethodGet, http.MethodOptions)
	auth.HandleFunc("/facebook/login", ctrl.Oauth.FacebookLogin).Methods(http.MethodGet, http.MethodOptions)
	auth.HandleFunc("/facebook/callback", ctrl.Oauth.FacebookCallback).Methods(http.MethodGet, http.MethodOptions)
	auth.HandleFunc("/github/login", ctrl.Oauth.GitHubLogin).Methods(http.MethodGet, http.MethodOptions)
	auth.HandleFunc("/github/callback", ctrl.Oauth.GitHubCallback).Methods(http.MethodGet, http.MethodOptions)

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))
	// Ensure protected routes also have CORS (inherited from api subrouter)

	// User routes
	protected.HandleFunc("/users/profile", ctrl.User.GetUserProfile).Methods(http.MethodGet, http.MethodOptions)
	protected.HandleFunc("/users/preferences", ctrl.User.UpdateUserPreferences).Methods(http.MethodPut, http.MethodOptions)
}
