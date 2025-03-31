package routes

import (
	"net/http"

	"github.com/gorilla/mux"

	// "libero-backend/config" // Removed if not needed directly
	"libero-backend/internal/api/controllers"
	"libero-backend/internal/middleware"
	"libero-backend/internal/service" // Import service package
)

// SetupRoutes configures all API routes
// It now accepts the main service struct to access all services
func SetupRoutes(router *mux.Router, service *service.Service) {
	// Extract controllers and services needed
	// Note: Controllers are instantiated here based on the service.
	// Alternatively, they could be instantiated in app.go and passed in.
	userController := controllers.NewUserController(service.User)
	oauthController := controllers.NewOAuthController(service.OAuth)
	authService := service.Auth // Get AuthService for middleware

	// API routes
	api := router.PathPrefix("/api").Subrouter()

	// Public routes (no authentication required)
	api.HandleFunc("/health", healthCheck).Methods(http.MethodGet)
	api.HandleFunc("/users/register", userController.Register).Methods(http.MethodPost)
	api.HandleFunc("/users/login", userController.Login).Methods(http.MethodPost) // Assumes Login handles password auth

	// OAuth routes (public) - Using root router for /auth path
	auth := router.PathPrefix("/auth").Subrouter()
	auth.HandleFunc("/google/login", oauthController.GoogleLogin).Methods(http.MethodGet)
	auth.HandleFunc("/google/callback", oauthController.GoogleCallback).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/login", oauthController.FacebookLogin).Methods(http.MethodGet)
	auth.HandleFunc("/facebook/callback", oauthController.FacebookCallback).Methods(http.MethodGet)
	auth.HandleFunc("/github/login", oauthController.GitHubLogin).Methods(http.MethodGet)
	auth.HandleFunc("/github/callback", oauthController.GitHubCallback).Methods(http.MethodGet)

	// Protected routes (authentication required)
	protected := api.PathPrefix("").Subrouter()
	protected.Use(middleware.AuthMiddleware(authService)) // Inject AuthService

	// User routes
	protected.HandleFunc("/users/profile", userController.GetProfile).Methods(http.MethodGet)
	protected.HandleFunc("/users/profile", userController.UpdateProfile).Methods(http.MethodPut)

	// Admin routes (requires admin role)
	admin := api.PathPrefix("/admin").Subrouter()
	admin.Use(middleware.AuthMiddleware(authService)) // Inject AuthService
	admin.Use(middleware.RoleMiddleware("admin")) // Example role check

	admin.HandleFunc("/users", userController.ListUsers).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id:[0-9]+}", userController.GetUser).Methods(http.MethodGet)
	admin.HandleFunc("/users/{id:[0-9]+}", userController.DeleteUser).Methods(http.MethodDelete)
}

// healthCheck is a simple health check endpoint
func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}