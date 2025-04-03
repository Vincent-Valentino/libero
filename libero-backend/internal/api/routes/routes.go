package routes

import (
	"net/http"

	"github.com/gorilla/mux"
	"libero-backend/config" // Added config dependency

	// "libero-backend/config" // Removed if not needed directly
	"libero-backend/internal/api/controllers"
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

	// Public routes (no authentication required)
	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	api.HandleFunc("/auth/register", ctrl.User.Register).Methods(http.MethodPost)
 // Changed path prefix
	api.HandleFunc("/auth/login", ctrl.User.Login).Methods(http.MethodPost) // Changed path prefix

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

	// User routes
	protected.HandleFunc("/users/profile", ctrl.User.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/profile", ctrl.User.UpdateProfile).Methods(http.MethodPut)

	// Admin routes (requires admin role)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware(authService)) // Inject AuthService
	admin.Use(middleware.RoleMiddleware("admin")) // Example role check

	admin.HandleFunc("/users", ctrl.User.ListUsers).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id:[0-9]+}", ctrl.User.GetUser).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id:[0-9]+}", ctrl.User.DeleteUser).Methods(http.MethodDelete)
}

// healthCheck is a simple health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}