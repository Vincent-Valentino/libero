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

	// Add global middlewares
	router.Use(middleware.LoggingMiddleware)
	router.Use(middleware.CORSMiddleware)

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Public routes (no authentication required)
	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	api.HandleFunc("/auth/register", ctrl.User.Register).Methods(http.MethodPost)
	api.HandleFunc("/auth/login", ctrl.User.Login).Methods(http.MethodPost)

	// Sports data routes
	api.HandleFunc("/sports/fixtures/today", ctrl.SportsData.HandleGetTodaysFixtures).Methods(http.MethodGet)
	api.HandleFunc("/sports/fixtures/summary", ctrl.SportsData.HandleGetFixturesSummary).Methods(http.MethodGet)
	api.HandleFunc("/standings", ctrl.SportsData.HandleGetStandings).Methods(http.MethodGet)
	api.HandleFunc("/topscorers", ctrl.SportsData.HandleGetTopScorers).Methods(http.MethodGet)
	api.HandleFunc("/matches/upcoming", ctrl.SportsData.HandleGetUpcomingMatches).Methods(http.MethodGet)
	api.HandleFunc("/matches/results", ctrl.SportsData.HandleGetResults).Methods(http.MethodGet)
	api.HandleFunc("/players/{player_id}/stats", ctrl.SportsData.HandleGetPlayerStats).Methods(http.MethodGet)

	// OAuth routes
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/google/login", ctrl.Oauth.GoogleLogin).Methods(http.MethodGet)
	auth.HandleFunc("/google/callback", ctrl.Oauth.GoogleCallback).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/login", ctrl.Oauth.FacebookLogin).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/callback", ctrl.Oauth.FacebookCallback).Methods(http.MethodGet)
	auth.HandleFunc("/github/login", ctrl.Oauth.GitHubLogin).Methods(http.MethodGet)
	auth.HandleFunc("/github/callback", ctrl.Oauth.GitHubCallback).Methods(http.MethodGet)

	// Protected routes (require authentication)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService))

	// User routes
	protected.HandleFunc("/users/profile", ctrl.User.GetUserProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/preferences", ctrl.User.UpdateUserPreferences).Methods(http.MethodPut)
}
