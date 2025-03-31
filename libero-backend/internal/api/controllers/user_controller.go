package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"

	"libero-backend/internal/middleware"
	"libero-backend/internal/models"
	"libero-backend/internal/service"
)

// UserController handles HTTP requests for user-related operations
type UserController struct {
	service service.UserService
 // Use interface value, not pointer
}

// NewUserController creates a new user controller instance
func NewUserController(service service.UserService) *UserController {
 // Use interface value
	return &UserController{
		service: service,
	}
}

// Register handles user registration requests
func (c *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var user models.User

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	// i think this part is wrong, because we haven't defined the fields in the struct yet
	if user.Email == "" || user.Username == "" || user.Password == "" {
		http.Error(w, "Email, username and password are required", http.StatusBadRequest)
		return
	}

	// Register user
	if err := c.service.RegisterUser(&user); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Return user data (without password)
	response := user.ToResponse()
	respondWithJSON(w, http.StatusCreated, response)
}

// Login handles user login requests
func (c *UserController) Login(w http.ResponseWriter, r *http.Request) {
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
	token, err := c.service.LoginUser(credentials.Email, credentials.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	// Return the JWT token
	respondWithJSON(w, http.StatusOK, map[string]string{"token": token})
}

// GetProfile handles requests to get the current user's profile
func (c *UserController) GetProfile(w http.ResponseWriter, r *http.Request) {
	// Get user from context (set by the auth middleware)
	claims, ok := middleware.GetUserFromContext(r.Context())
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	// Get user profile
	user, err := c.service.GetUserByID(claims.UserID)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, user.ToResponse()) // Use ToResponse
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
		// Improve error checking later (e.g., check for specific error type)
		if err.Error() == "user not found (placeholder)" {
			http.Error(w, "User not found", http.StatusNotFound)
		} else {
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
	if firstName, ok := updates["first_name"].(string); ok {
		user.FirstName = firstName
	}
	if lastName, ok := updates["last_name"].(string); ok {
		user.LastName = lastName
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
	respondWithJSON(w, http.StatusOK, user.ToResponse())
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

	respondWithJSON(w, http.StatusOK, user.ToResponse()) // Use ToResponse
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
	respondWithJSON(w, http.StatusOK, map[string]interface{}{
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

// Helper function to send JSON response
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}