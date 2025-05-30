package controllers

import (
	"encoding/json"
	"errors" // Added for error checking
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"libero-backend/internal/middleware"
	"libero-backend/internal/models"
	"libero-backend/internal/service"
	"libero-backend/internal/utils"
)

// UserController handles HTTP requests for user-related operations
type UserController struct {
	service     service.UserService
	authService service.AuthService // Added AuthService dependency
}

// NewUserController creates a new user controller instance
func NewUserController(userService service.UserService, authService service.AuthService) *UserController {
	return &UserController{
		service:     userService,
		authService: authService, // Store injected AuthService
	}
}

// Register handles user registration requests
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get context

	// Define a separate struct for registration request
	var registrationRequest struct {
		Name     string `json:"name"`
		Username string `json:"username"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&registrationRequest); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if registrationRequest.Email == "" || registrationRequest.Username == "" || registrationRequest.Password == "" || registrationRequest.Name == "" {
		http.Error(w, "Name, email, username and password are required", http.StatusBadRequest)
		return
	}

	// Create user model from registration request
	user := &models.User{
		Name:     registrationRequest.Name,
		Email:    registrationRequest.Email,
		Username: registrationRequest.Username,
		Password: registrationRequest.Password, // This will be hashed by BeforeSave hook
	}

	// Register user using AuthService
	if err := c.authService.RegisterByPassword(ctx, user); err != nil {
		// Handle specific registration errors
		if errors.Is(err, service.ErrEmailAlreadyExists) {
			http.Error(w, "Email address already in use", http.StatusConflict)
		} else if errors.Is(err, service.ErrUsernameAlreadyExists) {
			http.Error(w, "Username already in use", http.StatusConflict)
		} else if errors.Is(err, service.ErrWeakPassword) {
			http.Error(w, "Password does not meet security requirements (minimum 8 characters)", http.StatusBadRequest)
		} else {
			// Log the actual error for debugging
			// TODO: Replace with structured logging
			http.Error(w, "Failed to create user account", http.StatusInternalServerError)
		}
		return
	}

	// Return user data (without password)
	response := user.ToResponse()
	utils.RespondWithJSON(w, http.StatusCreated, response)
}

// Login handles user login requests
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get context
	var credentials struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&credentials); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if credentials.Email == "" || credentials.Password == "" {
		http.Error(w, "Email and password are required", http.StatusBadRequest)
		return
	}

	// Authenticate user
	token, err := c.authService.LoginByPassword(ctx, credentials.Email, credentials.Password)
	if err != nil {
		// Handle specific authentication errors
		if errors.Is(err, service.ErrInvalidCredentials) || errors.Is(err, service.ErrAccountInactive) {
			http.Error(w, err.Error(), http.StatusUnauthorized)
		} else {
			// Log internal errors (replace with proper logging later)
			http.Error(w, "Internal server error during login", http.StatusInternalServerError)
		}
		return
	}

	// Return the JWT token
	utils.RespondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// RequestPasswordReset handles password reset requests
func (c *UserController) RequestPasswordReset(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		Email string `json:"email"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Email == "" {
		http.Error(w, "Email is required", http.StatusBadRequest)
		return
	}

	// Request password reset token
	token, err := c.authService.RequestPasswordReset(ctx, request.Email)
	if err != nil {
		// Log error but don't reveal to user for security
		// TODO: Add structured logging
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// For testing purposes, return the token
	// In production, this would send an email instead
	if token != "" {
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "Password reset token generated",
			"token":   token, // Remove this in production
		})
	} else {
		// Email not found, but don't reveal this for security
		utils.RespondWithJSON(w, http.StatusOK, map[string]string{
			"message": "If the email exists, a password reset token has been sent",
		})
	}
}

// ResetPassword handles password reset completion
func (c *UserController) ResetPassword(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var request struct {
		Token           string `json:"token"`
		NewPassword     string `json:"new_password"`
		ConfirmPassword string `json:"confirm_password"`
	}

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.Token == "" || request.NewPassword == "" || request.ConfirmPassword == "" {
		http.Error(w, "Token, new password, and confirm password are required", http.StatusBadRequest)
		return
	}

	// Reset password
	err := c.authService.ResetPassword(ctx, request.Token, request.NewPassword, request.ConfirmPassword)
	if err != nil {
		if errors.Is(err, service.ErrPasswordMismatch) {
			http.Error(w, "Passwords do not match", http.StatusBadRequest)
		} else if errors.Is(err, service.ErrWeakPassword) {
			http.Error(w, "Password does not meet security requirements (minimum 8 characters)", http.StatusBadRequest)
		} else if errors.Is(err, service.ErrResetTokenInvalid) {
			http.Error(w, "Invalid or expired reset token", http.StatusBadRequest)
		} else {
			// Log error
			// TODO: Add structured logging
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]string{
		"message": "Password reset successfully",
	})
}

// GetUserProfile handles requests to get the current user's profile including preferences
func (c *UserController) GetUserProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get user from context (set by the auth middleware)
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Get user profile with preferences from service
	profile, err := c.service.GetUserProfile(ctx, claims.UserID)
	if err != nil {
		// TODO: Add more specific error handling (e.g., check for gorm.ErrRecordNotFound)
		// For now, assume not found or internal error
		// Log the actual error server-side
		// log.Printf("Error fetching user profile for user %d: %v", claims.UserID, err)
		http.Error(w, "Failed to retrieve user profile", http.StatusNotFound) // Or InternalServerError depending on error type
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, profile)
}

// UpdateProfile handles requests to update the current user's profile
func (c *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context() // Get context from request
	// Get user from context
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Parse request body
	var updates map[string]interface{}
	if err := json.NewDecoder(r.Body).Decode(&updates); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// 1. Fetch the existing user
	user, err := c.service.GetUserByID(claims.UserID)
	if err != nil {
		// Handle user not found or other errors
		// Check for specific error type
		if errors.Is(err, service.ErrUserNotFound) { // Assuming UserService returns this error
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
			// TODO: Add structured logging (Failed to retrieve user: %v, err)
			http.Error(w, "Failed to retrieve user", http.StatusInternalServerError)
		}
		return
	}

	// 2. Apply updates from the request body (basic example)
	// Note: This is a simplified update logic. Production code might need reflection,
	// struct tags, or libraries like 'mergo' for more robust merging.
	// Also, consider validation and which fields are allowed to be updated.
	if name, ok := updates["name"].(string); ok {
		user.Name = name
	}
	// Add other updatable fields here (e.g., Username, but be careful with unique constraints)
	// Do NOT update Email, Password, Role, Provider, ProviderID, etc. here unless intended.

	// 3. Call the service to save the updated user
	err = c.service.UpdateUser(ctx, user) // Pass context and updated user object
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest) // Or InternalServerError depending on error
		return
	}

	// Return the updated user profile (without sensitive info)
	utils.RespondWithJSON(w, http.StatusOK, user.ToResponse())
}

// GetUser handles requests to get a specific user by ID
func (c *UserController) GetUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL parameters
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Get user
	user, err := c.service.GetUserByID(uint(id))
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, user.ToResponse()) // Use ToResponse
}

// ListUsers handles requests to list users with pagination
func (c *UserController) ListUsers(w http.ResponseWriter, r *http.Request) {
	// Get pagination parameters
	page, err := strconv.Atoi(r.URL.Query().Get("page"))
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || limit < 1 {
		limit = 10
	}

	// Get users
	users, total, err := c.service.ListUsers(page, limit)
	if err != nil {
		http.Error(w, "Failed to retrieve users", http.StatusInternalServerError)
		return
	}

	// Convert users to UserResponse
	userResponses := make([]models.UserResponse, len(users))
	for i, u := range users {
		userResponses[i] = u.ToResponse()
	}

	// Return users with pagination metadata
	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"users":       userResponses, // Return responses
		"total":       total,
		"page":        page,
		"limit":       limit,
		"total_pages": (total + int64(limit) - 1) / int64(limit),
	})
}

// DeleteUser handles requests to delete a user
func (c *UserController) DeleteUser(w http.ResponseWriter, r *http.Request) {
	// Extract user ID from URL parameters
	vars := mux.Vars(r)
	idStr := vars["id"]

	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	// Delete user
	if err := c.service.DeleteUser(uint(id)); err != nil {
		http.Error(w, "Failed to delete user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// UpdateUserPreferences handles requests to update the current user's preferences
func (c *UserController) UpdateUserPreferences(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Get user from context
	claims, ok := middleware.GetUserFromContext(ctx)
	if !ok {
		utils.RespondWithJSON(w, http.StatusUnauthorized, map[string]string{"error": "Unauthorized"})
		return
	}

	// Parse request body
	var req models.UpdatePreferencesRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		utils.RespondWithJSON(w, http.StatusBadRequest, map[string]string{"error": "Invalid request payload"})
		return
	}

	// Call the service to update preferences
	err := c.service.UpdateUserPreferences(ctx, claims.UserID, req)
	if err != nil {
		// The service currently logs warnings but returns nil.
		// If the service were to return errors (e.g., for invalid IDs), handle them here.
		// log.Printf("Error updating preferences for user %d: %v", claims.UserID, err)
		utils.RespondWithJSON(w, http.StatusInternalServerError, map[string]string{"error": "Failed to update preferences"})
		return
	}

	// Respond with No Content on success
	w.WriteHeader(http.StatusNoContent)
}
