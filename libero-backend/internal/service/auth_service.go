package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"libero-backend/config"          // Added for JWT config
	"libero-backend/internal/models" // Added for User model
	"time"                           // Added for JWT expiration

	"github.com/golang-jwt/jwt/v5" // Added JWT library
	"gorm.io/gorm"                 // Added GORM for error checking
)

// Predefined errors for authentication
var (
	ErrInvalidCredentials    = errors.New("invalid email or password")
	ErrUserNotFound          = errors.New("user not found") // Potentially replace with specific repo error check
	ErrTokenInvalid          = errors.New("invalid or expired token")
	ErrAccountInactive       = errors.New("user account is inactive")
	ErrPasswordMismatch      = errors.New("new password doesn't match confirmation")
	ErrSamePassword          = errors.New("new password cannot be the same as the current password")
	ErrWeakPassword          = errors.New("password does not meet security requirements")
	ErrEmailAlreadyExists    = errors.New("email address already in use")
	ErrUsernameAlreadyExists = errors.New("username already in use")
	ErrResetTokenInvalid     = errors.New("reset token is invalid or has expired")
)

// AuthService defines the interface for authentication operations.
type AuthService interface {
	// LoginOrRegisterViaProvider handles user lookup or creation after successful OAuth.
	// It takes the UserInfo fetched by OAuthService and returns a JWT string.
	LoginOrRegisterViaProvider(ctx context.Context, userInfo *UserInfo) (string, error)

	// Added methods for password login and JWT validation
	LoginByPassword(ctx context.Context, email, password string) (string, error)
	ValidateJWTToken(tokenString string) (*JWTClaims, error)

	// Password registration
	RegisterByPassword(ctx context.Context, user *models.User) error

	// Password management
	ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword, confirmPassword string) error
	RequestPasswordReset(ctx context.Context, email string) (string, error)
	ResetPassword(ctx context.Context, token, newPassword, confirmPassword string) error
}

// JWTClaims defines the structure for custom JWT claims
// Exported for use by middleware
type JWTClaims struct {
	UserID uint   `json:"user_id"`
	Email  string `json:"email"`
	Role   string `json:"role"`
	jwt.RegisteredClaims
}

// authService implements the AuthService interface.
type authService struct {
	userService UserService      // Dependency on UserService
	jwtCfg      config.JWTConfig // Dependency on JWT configuration
}

// NewAuthService creates a new AuthService instance.
// Dependencies will be injected here.
func NewAuthService(userService UserService, jwtCfg config.JWTConfig) AuthService {
	return &authService{
		userService: userService,
		jwtCfg:      jwtCfg,
	}
}

// --- JWT Helper ---

// generateJWTToken creates a new JWT token for a given user.
func (s *authService) generateJWTToken(user *models.User) (string, error) {
	if user == nil {
		return "", errors.New("cannot generate token for nil user")
	}
	// Use a strong, configured secret key
	if s.jwtCfg.Secret == "" || s.jwtCfg.Secret == "your_secret_key" { // Check against default/empty
		// TODO: Add structured logging (Warning: JWT Secret is not configured or using default value.)
		return "", errors.New("JWT secret is not securely configured")
	}

	expirationTime := time.Now().Add(time.Second * time.Duration(s.jwtCfg.ExpiresIn))

	claims := &JWTClaims{
		UserID: user.ID,
		Email:  user.Email,
		Role:   user.Role, // Ensure Role is populated correctly in User model
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			// Consider adding Issuer, Subject, Audience for better token validation
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(s.jwtCfg.Secret))
	if err != nil {
		// TODO: Add structured logging (Error signing token: %v, err)
		return "", errors.New("could not generate token")
	}

	return tokenString, nil
}

// --- Interface Implementations ---

// LoginOrRegisterViaProvider implements the logic to find an existing user
// based on provider info or create a new one, returning a JWT.
func (s *authService) LoginOrRegisterViaProvider(ctx context.Context, userInfo *UserInfo) (string, error) {
	// 1. Check if user exists by ProviderID
	user, err := s.userService.FindUserByProvider(ctx, userInfo.Provider, userInfo.ProviderID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) { // Use errors.Is
		// TODO: Add structured logging (Error finding user by provider %s (%s): %v, userInfo.Provider, userInfo.ProviderID, err)
		return "", fmt.Errorf("database error checking provider identity: %w", err)
	}

	// 2. If user exists by ProviderID, generate token.
	if user != nil {
		// TODO: Add structured logging (AuthService: Found existing user by provider: ID=%d, user.ID)
		// TODO: Potentially update user details (e.g., Name) if they differ from userInfo
		// if user.Name != userInfo.Name && userInfo.Name != "" { user.Name = userInfo.Name; /* call update */ }
		// updateErr := s.userService.UpdateUser(ctx, user) ... (handle error)
		return s.generateJWTToken(user) // Generate JWT
	}

	// 3. If user does not exist by ProviderID, check by Email (if available)
	if userInfo.Email != "" {
		existingUserByEmail, emailErr := s.userService.FindUserByEmail(ctx, userInfo.Email)
		if emailErr != nil && !errors.Is(emailErr, gorm.ErrRecordNotFound) { // Use errors.Is
			// TODO: Add structured logging (Error finding user by email %s: %v, userInfo.Email, emailErr)
			return "", fmt.Errorf("database error checking email: %w", emailErr)
		}

		if existingUserByEmail != nil {
			// User exists with this email but different/no provider link. Link them.
			// TODO: Add structured logging (AuthService: Found existing user by email, linking provider: ID=%d, existingUserByEmail.ID)
			existingUserByEmail.Provider = userInfo.Provider
			existingUserByEmail.ProviderID = userInfo.ProviderID
			// Potentially update Name if empty or different?
			needsUpdate := false
			if existingUserByEmail.Name == "" && userInfo.Name != "" {
				existingUserByEmail.Name = userInfo.Name
				needsUpdate = true
			}
			// Ensure the user has a username if they don't have one (for OAuth linking)
			if existingUserByEmail.Username == "" {
				existingUserByEmail.Username = fmt.Sprintf("%s_%s", userInfo.Provider, userInfo.ProviderID)
				needsUpdate = true
			}
			// Add other fields to potentially update here

			if needsUpdate {
				updateErr := s.userService.UpdateUser(ctx, existingUserByEmail)
				if updateErr != nil {
					// Log the error but proceed with login
					// TODO: Add structured logging (AuthService: Warning - failed to link provider %s to existing user %d: %v, userInfo.Provider, existingUserByEmail.ID, updateErr)
				}
			}
			return s.generateJWTToken(existingUserByEmail) // Generate JWT
		}
	}

	// 4. If user still not found, create a new user.
	// TODO: Add structured logging (AuthService: Creating new user for provider %s, email %s, userInfo.Provider, userInfo.Email)

	// For GitHub users without email, generate a placeholder email
	email := userInfo.Email
	if email == "" && userInfo.Provider == "github" {
		// Create a placeholder email for GitHub users who don't provide email
		email = fmt.Sprintf("%s+github@noemail.local", userInfo.ProviderID)
	}

	// Validate we have required fields for user creation
	if email == "" {
		return "", fmt.Errorf("cannot create user without email address for provider %s", userInfo.Provider)
	}

	if userInfo.Name == "" {
		// Use provider ID as fallback name
		userInfo.Name = fmt.Sprintf("%s_user_%s", userInfo.Provider, userInfo.ProviderID)
	}

	// Generate a unique username for OAuth users to avoid database constraint violations
	// Format: provider_providerID (e.g., google_123456789, github_987654321)
	generatedUsername := fmt.Sprintf("%s_%s", userInfo.Provider, userInfo.ProviderID)

	newUser := &models.User{
		Email:      email,             // Use processed email (may be placeholder for GitHub)
		Username:   generatedUsername, // Generate unique username for OAuth users
		Name:       userInfo.Name,     // Ensure name is not empty
		Provider:   userInfo.Provider,
		ProviderID: userInfo.ProviderID,
		Role:       "user", // Default role
		Active:     true,   // Default active
		// Password will be empty for OAuth users (now nullable in model)
	}
	// Use CreateUser which should handle basic creation.
	createdUser, createErr := s.userService.CreateUser(ctx, newUser)
	if createErr != nil {
		// TODO: Add structured logging (Error creating new user for provider %s: %v, userInfo.Provider, createErr)
		return "", fmt.Errorf("failed to create new user: %w", createErr)
	}
	// TODO: Add structured logging (AuthService: Created new user with ID: %d, createdUser.ID)

	// Assuming CreateUser returns the user with ID set.
	return s.generateJWTToken(createdUser) // Generate JWT
}

// LoginByPassword handles standard email/password authentication.
func (s *authService) LoginByPassword(ctx context.Context, email, password string) (string, error) {
	user, err := s.userService.FindUserByEmail(ctx, email)
	if err != nil {
		// Use specific error check if available from repository/service
		// Check if the error is specifically ErrUserNotFound
		if errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials // Don't reveal if user exists
		}
		// TODO: Replace with structured logging
		// (Error finding user by email %s during login: %v, email, err)
		return "", fmt.Errorf("error during authentication") // Generic error
	}

	if user == nil {
		return "", ErrInvalidCredentials
	}

	// Compare the provided password with the stored hash
	if !user.ComparePassword(password) {
		return "", ErrInvalidCredentials
	}

	// Check if user is active
	if !user.Active {
		return "", ErrAccountInactive
	}

	// Generate JWT token
	tokenString, err := s.generateJWTToken(user)
	if err != nil {
		// Error already logged in generateJWTToken
		return "", errors.New("could not process login") // Generic error
	}

	return tokenString, nil
}

// ValidateJWTToken parses and validates a JWT string.
func (s *authService) ValidateJWTToken(tokenString string) (*JWTClaims, error) {
	if s.jwtCfg.Secret == "" || s.jwtCfg.Secret == "your_secret_key" { // Check against default/empty
		// TODO: Add structured logging (Warning: JWT Secret is not configured or using default value.)
		return nil, errors.New("JWT secret is not securely configured")
	}

	claims := &JWTClaims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		// Validate the alg is what we expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(s.jwtCfg.Secret), nil
	})

	if err != nil {
		// Check for specific errors
		if errors.Is(err, jwt.ErrTokenMalformed) {
			// TODO: Add structured logging (Malformed token: %v, err)
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			// TODO: Add structured logging (Token is expired: %v, err)
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			// TODO: Add structured logging (Token not active yet: %v, err)
			return nil, ErrTokenInvalid
		} else {
			// Log other parsing errors
			// TODO: Add structured logging (Couldn't handle this token: %v, err)
			return nil, ErrTokenInvalid // Generic invalid token error
		}
	}

	if !token.Valid {
		// TODO: Add structured logging (Token validation failed.)
		return nil, ErrTokenInvalid
	}

	// We can trust claims now
	return claims, nil
}

// RegisterByPassword registers a new user with email and password
func (s *authService) RegisterByPassword(ctx context.Context, user *models.User) error {
	// Check if email already exists
	existingUser, err := s.userService.FindUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		// If there's a database error (not just "not found"), return it
		return fmt.Errorf("error checking email availability: %w", err)
	}
	if existingUser != nil {
		return ErrEmailAlreadyExists
	}

	// Check if username already exists
	if user.Username != "" {
		existingUser, err = s.userService.FindUserByUsername(ctx, user.Username)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			// If there's a database error (not just "not found"), return it
			return fmt.Errorf("error checking username availability: %w", err)
		}
		if existingUser != nil {
			return ErrUsernameAlreadyExists
		}
	}

	// Validate password strength (implement your own criteria)
	if len(user.Password) < 8 {
		return ErrWeakPassword
	}

	// Set default values
	user.Role = "user"
	user.Active = true

	// Create user (password will be hashed by BeforeSave hook in User model)
	_, err = s.userService.CreateUser(ctx, user)
	if err != nil {
		return fmt.Errorf("failed to create user: %w", err)
	}

	return nil
}

// ChangePassword changes the password for an authenticated user
func (s *authService) ChangePassword(ctx context.Context, userID uint, currentPassword, newPassword, confirmPassword string) error {
	// Check if new password and confirmation match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}

	// Check password strength
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}

	// Get current user
	user, err := s.userService.GetUserByID(userID)
	if err != nil {
		return ErrUserNotFound
	}

	// Verify current password
	if !user.ComparePassword(currentPassword) {
		return ErrInvalidCredentials
	}

	// Check if new password is different from old password
	if user.ComparePassword(newPassword) {
		return ErrSamePassword
	}

	// Update password
	user.Password = newPassword
	return s.userService.UpdateUser(ctx, user)
}

// generateResetToken creates a secure random token
func generateResetToken() (string, error) {
	b := make([]byte, 32)
	_, err := rand.Read(b)
	if err != nil {
		return "", err
	}

	return base64.URLEncoding.EncodeToString(b), nil
}

// RequestPasswordReset sends a password reset token for a user
func (s *authService) RequestPasswordReset(ctx context.Context, email string) (string, error) {
	// Find user by email
	user, err := s.userService.FindUserByEmail(ctx, email)
	if err != nil {
		// Don't reveal if email exists or not for security reasons
		return "", nil
	}

	// Generate reset token
	resetToken, err := generateResetToken()
	if err != nil {
		return "", fmt.Errorf("failed to generate reset token: %w", err)
	}

	// Set reset token and expiration time (1 hour)
	expiresAt := time.Now().Add(1 * time.Hour)

	// Store token in user record
	user.ResetToken = resetToken
	user.ResetTokenExpiresAt = expiresAt
	err = s.userService.UpdateUser(ctx, user)
	if err != nil {
		return "", fmt.Errorf("failed to store reset token: %w", err)
	}

	// In production, you would send this token via email instead of returning it
	// For testing purposes, we're returning the token
	return resetToken, nil
}

// ResetPassword resets a user's password using a valid token
func (s *authService) ResetPassword(ctx context.Context, token, newPassword, confirmPassword string) error {
	// Check if new password and confirmation match
	if newPassword != confirmPassword {
		return ErrPasswordMismatch
	}

	// Check password strength
	if len(newPassword) < 8 {
		return ErrWeakPassword
	}

	// Find user by reset token
	user, err := s.userService.FindUserByResetToken(ctx, token)
	if err != nil {
		return ErrResetTokenInvalid
	}

	// Check if token is expired
	if user.ResetTokenExpiresAt.Before(time.Now()) {
		return ErrResetTokenInvalid
	}

	// Update password and clear reset token fields
	user.Password = newPassword
	user.ResetToken = ""
	user.ResetTokenExpiresAt = time.Time{}
	return s.userService.UpdateUser(ctx, user)
}

// Add other AuthService method implementations here...
