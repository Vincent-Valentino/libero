package service

import (
	"context"
	"errors"
	"fmt"
	"libero-backend/config" // Added for JWT config
	"libero-backend/internal/models" // Added for User model
	"time"                           // Added for JWT expiration

	"github.com/golang-jwt/jwt/v5" // Added JWT library
)

// Predefined errors for authentication
var (
	ErrInvalidCredentials = errors.New("invalid email or password")
	ErrUserNotFound       = errors.New("user not found") // Potentially replace with specific repo error check
	ErrTokenInvalid       = errors.New("invalid or expired token")
	ErrAccountInactive    = errors.New("user account is inactive")
)


// AuthService defines the interface for authentication operations.
type AuthService interface {
	// LoginOrRegisterViaProvider handles user lookup or creation after successful OAuth.
	// It takes the UserInfo fetched by OAuthService and returns a JWT string.
	LoginOrRegisterViaProvider(ctx context.Context, userInfo *UserInfo) (string, error)

	// Added methods for password login and JWT validation
	LoginByPassword(ctx context.Context, email, password string) (string, error)
	ValidateJWTToken(tokenString string) (*JWTClaims, error)
	// TODO: Add methods for registration (e.g., RegisterByPassword) (Is this covered by UserService?)
	// TODO: Add methods for password management (e.g., ChangePassword, RequestPasswordReset)
}

// JWTClaims defines the structure for custom JWT claims
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
	    fmt.Println("Warning: JWT Secret is not configured or using default value.") // Replace with proper logging
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
		fmt.Printf("Error signing token: %v\n", err) // Replace with proper logging
		return "", errors.New("could not generate token")
	}

	return tokenString, nil
}


// --- Interface Implementations ---

// LoginOrRegisterViaProvider implements the logic to find an existing user
// based on provider info or create a new one, returning a JWT.
func (s *authService) LoginOrRegisterViaProvider(ctx context.Context, userInfo *UserInfo) (string, error) {
	// 1. Check if user exists by ProviderID
	// TODO: Replace placeholder error check with proper error type checking (e.g., errors.Is(err, gorm.ErrRecordNotFound))
	user, err := s.userService.FindUserByProvider(ctx, userInfo.Provider, userInfo.ProviderID)
	if err != nil && err.Error() != "user not found (placeholder)" { // Assuming placeholder error for now
		fmt.Printf("Error finding user by provider %s (%s): %v\n", userInfo.Provider, userInfo.ProviderID, err) // Log error
		return "", fmt.Errorf("database error checking provider identity")
	}

	// 2. If user exists by ProviderID, generate token.
	if user != nil {
		fmt.Printf("AuthService: Found existing user by provider: ID=%d\n", user.ID)
		// TODO: Potentially update user details (e.g., Name) if they differ from userInfo
		// if user.Name != userInfo.Name && userInfo.Name != "" { user.Name = userInfo.Name; /* call update */ }
		// updateErr := s.userService.UpdateUser(ctx, user) ... (handle error)
		return s.generateJWTToken(user) // Generate JWT
	}

	// 3. If user does not exist by ProviderID, check by Email (if available)
	if userInfo.Email != "" {
		// TODO: Replace placeholder error check with proper error type checking
		existingUserByEmail, emailErr := s.userService.FindUserByEmail(ctx, userInfo.Email)
		if emailErr != nil && emailErr.Error() != "user not found (placeholder)" { // Assuming placeholder error for now
			fmt.Printf("Error finding user by email %s: %v\n", userInfo.Email, emailErr) // Log error
			return "", fmt.Errorf("database error checking email")
		}

		if existingUserByEmail != nil {
			// User exists with this email but different/no provider link. Link them.
			fmt.Printf("AuthService: Found existing user by email, linking provider: ID=%d\n", existingUserByEmail.ID)
			existingUserByEmail.Provider = userInfo.Provider
			existingUserByEmail.ProviderID = userInfo.ProviderID
			// Potentially update Name if empty or different?
			needsUpdate := false
			if existingUserByEmail.Name == "" && userInfo.Name != "" {
				existingUserByEmail.Name = userInfo.Name
				needsUpdate = true
			}
			// Add other fields to potentially update here

			if needsUpdate {
				updateErr := s.userService.UpdateUser(ctx, existingUserByEmail)
				if updateErr != nil {
					// Log the error but proceed with login
					fmt.Printf("AuthService: Warning - failed to link provider %s to existing user %d: %v\n", userInfo.Provider, existingUserByEmail.ID, updateErr)
				}
			}
			return s.generateJWTToken(existingUserByEmail) // Generate JWT
		}
	}

	// 4. If user still not found, create a new user.
	fmt.Printf("AuthService: Creating new user for provider %s, email %s\n", userInfo.Provider, userInfo.Email)
	newUser := &models.User{
		Email:      userInfo.Email, // Ensure email is not empty if required by DB schema
		Name:       userInfo.Name,  // Ensure name is not empty if required
		Provider:   userInfo.Provider,
		ProviderID: userInfo.ProviderID,
		Role:       "user", // Default role
		Active:     true,   // Default active
		// Password will be empty/nil initially for OAuth users.
	}
	// Use CreateUser which should handle basic creation.
	createdUser, createErr := s.userService.CreateUser(ctx, newUser)
	if createErr != nil {
		fmt.Printf("Error creating new user for provider %s: %v\n", userInfo.Provider, createErr) // Log error
		return "", fmt.Errorf("failed to create new user")
	}
	fmt.Printf("AuthService: Created new user with ID: %d\n", createdUser.ID)

	// Assuming CreateUser returns the user with ID set.
	return s.generateJWTToken(createdUser) // Generate JWT
}

// LoginByPassword handles standard email/password authentication.
func (s *authService) LoginByPassword(ctx context.Context, email, password string) (string, error) {
	// TODO: Replace placeholder error check with proper error type checking
	user, err := s.userService.FindUserByEmail(ctx, email)
	if err != nil {
		// Use specific error check if available from repository/service
		if err.Error() == "user not found (placeholder)" || errors.Is(err, ErrUserNotFound) {
			return "", ErrInvalidCredentials // Don't reveal if user exists
		}
		fmt.Printf("Error finding user by email %s during login: %v\n", email, err) // Log internal error
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
	    fmt.Println("Warning: JWT Secret is not configured or using default value.") // Replace with proper logging
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
			fmt.Println("Malformed token:", err)
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenExpired) {
			fmt.Println("Token is expired:", err)
			return nil, ErrTokenInvalid
		} else if errors.Is(err, jwt.ErrTokenNotValidYet) {
			fmt.Println("Token not active yet:", err)
			return nil, ErrTokenInvalid
		} else {
			// Log other parsing errors
			fmt.Printf("Couldn't handle this token: %v\n", err) // Replace with proper logging
			return nil, ErrTokenInvalid // Generic invalid token error
		}
	}

	if !token.Valid {
		fmt.Println("Token validation failed.")
		return nil, ErrTokenInvalid
	}

	// We can trust claims now
	return claims, nil
}

// Add other AuthService method implementations here...