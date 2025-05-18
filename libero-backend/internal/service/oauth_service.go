package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json" // Added for potential user info parsing
	"errors"        // Added for error handling
	"fmt"
	"io"

	"libero-backend/config"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
)

// OAuthService defines the interface for OAuth related operations
type OAuthService interface {
	GetGoogleLoginURL() (string, string)
	HandleGoogleCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error)
	GetFacebookLoginURL() (string, string)
	HandleFacebookCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error)
	GetGitHubLoginURL() (string, string)
	HandleGitHubCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error)
}

// oauthService implements the OAuthService interface
type oauthService struct {
	GoogleConfig   *oauth2.Config
	FacebookConfig *oauth2.Config
	GitHubConfig   *oauth2.Config
	authService    AuthService // Added AuthService dependency
}

// NewOAuthService creates a new OAuthService instance
func NewOAuthService(cfg *config.Config, authService AuthService) OAuthService { // Return interface type
	googleCfg := &oauth2.Config{
		ClientID:     cfg.Google.ClientID,
		ClientSecret: cfg.Google.ClientSecret,
		RedirectURL:  cfg.Google.RedirectURL,
		Scopes:       []string{"https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/userinfo.profile"},
		Endpoint:     google.Endpoint,
	}

	facebookCfg := &oauth2.Config{
		ClientID:     cfg.Facebook.ClientID,
		ClientSecret: cfg.Facebook.ClientSecret,
		RedirectURL:  cfg.Facebook.RedirectURL,
		Scopes:       []string{"email", "public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	githubCfg := &oauth2.Config{
		ClientID:     cfg.GitHub.ClientID,
		ClientSecret: cfg.GitHub.ClientSecret,
		RedirectURL:  cfg.GitHub.RedirectURL,
		Scopes:       []string{"read:user", "user:email"},
		Endpoint:     github.Endpoint,
	}

	return &oauthService{ // Return pointer to the concrete implementation struct
		GoogleConfig:   googleCfg,
		FacebookConfig: facebookCfg,
		GitHubConfig:   githubCfg,
		authService:    authService, // Store injected AuthService
	}
}

// --- Helper ---

// generateStateOauthCookie generates a random state string for CSRF protection.
func generateStateOauthCookie() string {
	b := make([]byte, 16)
	rand.Read(b)
	return base64.URLEncoding.EncodeToString(b)
}

// UserInfo represents common user details fetched from an OAuth provider.
// This might be refined or moved later, possibly into an auth or user service package.
type UserInfo struct {
	Provider     string
	ProviderID   string
	Email        string
	Name         string
	AccessToken  string                 // The provider's access token
	RefreshToken string                 // Optional: The provider's refresh token
	RawData      map[string]interface{} // Raw data from provider for flexibility
}

// --- Google ---

// GetGoogleLoginURL generates the Google OAuth login URL.
// It returns the URL and the state string (which should be stored temporarily, e.g., in a session/cookie).
func (s *oauthService) GetGoogleLoginURL() (string, string) {
	state := generateStateOauthCookie()
	// Add AccessTypeOffline to request a refresh token
	// Add PromptSelectAccount to force account selection
	url := s.GoogleConfig.AuthCodeURL(state, oauth2.AccessTypeOffline, oauth2.ApprovalForce) // Consider oauth2.ApprovalForce if needed
	return url, state
}

// HandleGoogleCallback handles the callback from Google, exchanges the code for a token,
// fetches user info, and then calls AuthService to login or register the user.
// It returns a session token/ID (placeholder string for now) or an error.
func (s *oauthService) HandleGoogleCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error) {
	if receivedState != storedState {
		return "", errors.New("invalid oauth state")
	}

	token, err := s.GoogleConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	// Fetch user info from Google API
	client := s.GoogleConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo")
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %w", err)
	}
	defer response.Body.Close()

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading user info response body: %w", err)
	}

	var googleUserInfo map[string]interface{}
	if err := json.Unmarshal(contents, &googleUserInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal google user info: %w", err)
	}

	// Extract necessary fields
	userInfo := &UserInfo{
		Provider:     "google",
		ProviderID:   fmt.Sprintf("%v", googleUserInfo["id"]),
		Email:        fmt.Sprintf("%v", googleUserInfo["email"]),
		Name:         fmt.Sprintf("%v", googleUserInfo["name"]),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		RawData:      googleUserInfo,
	}

	// Call AuthService to handle login or registration
	tokenString, err := s.authService.LoginOrRegisterViaProvider(ctx, userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to login or register via google: %w", err)
	}

	return tokenString, nil
}

// --- Facebook ---

// GetFacebookLoginURL generates the Facebook OAuth login URL.
func (s *oauthService) GetFacebookLoginURL() (string, string) {
	state := generateStateOauthCookie()
	url := s.FacebookConfig.AuthCodeURL(state)
	return url, state
}

// HandleFacebookCallback handles the callback from Facebook.
// It returns a session token/ID (placeholder string for now) or an error.
func (s *oauthService) HandleFacebookCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error) {
	if receivedState != storedState {
		return "", errors.New("invalid oauth state")
	}

	token, err := s.FacebookConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	// Fetch user info from Facebook Graph API
	client := s.FacebookConfig.Client(ctx, token)
	// Use the token directly in the URL as Facebook's client might not automatically add it
	resp, err := client.Get("https://graph.facebook.com/me?fields=id,name,email&access_token=" + token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed getting user info from facebook: %w", err)
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading facebook user info response body: %w", err)
	}

	var fbUserInfo map[string]interface{}
	if err := json.Unmarshal(contents, &fbUserInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal facebook user info: %w", err)
	}

	// Extract necessary fields
	userInfo := &UserInfo{
		Provider:     "facebook",
		ProviderID:   fmt.Sprintf("%v", fbUserInfo["id"]),
		Email:        fmt.Sprintf("%v", fbUserInfo["email"]), // May be empty if user didn't grant permission
		Name:         fmt.Sprintf("%v", fbUserInfo["name"]),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken, // Facebook might not provide refresh tokens by default
		RawData:      fbUserInfo,
	}

	// Call AuthService to handle login or registration
	tokenString, err := s.authService.LoginOrRegisterViaProvider(ctx, userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to login or register via facebook: %w", err)
	}

	return tokenString, nil
}

// --- GitHub ---

// GetGitHubLoginURL generates the GitHub OAuth login URL.
func (s *oauthService) GetGitHubLoginURL() (string, string) {
	state := generateStateOauthCookie()
	url := s.GitHubConfig.AuthCodeURL(state)
	return url, state
}

// HandleGitHubCallback handles the callback from GitHub.
// It returns a session token/ID (placeholder string for now) or an error.
func (s *oauthService) HandleGitHubCallback(ctx context.Context, storedState string, receivedState string, code string) (string, error) {
	if receivedState != storedState {
		return "", errors.New("invalid oauth state")
	}

	token, err := s.GitHubConfig.Exchange(ctx, code)
	if err != nil {
		return "", fmt.Errorf("code exchange failed: %w", err)
	}

	// Fetch user info from GitHub API
	client := s.GitHubConfig.Client(ctx, token)
	resp, err := client.Get("https://api.github.com/user")
	if err != nil {
		return "", fmt.Errorf("failed getting user info from github: %w", err)
	}
	defer resp.Body.Close()

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading github user info response body: %w", err)
	}

	var ghUserInfo map[string]interface{}
	if err := json.Unmarshal(contents, &ghUserInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal github user info: %w", err)
	}

	// GitHub might not return email directly from /user, may need /user/emails
	// For simplicity, we'll try to get it from the main payload first.
	email := ""
	if val, ok := ghUserInfo["email"].(string); ok && val != "" {
		email = val
	} else {
		// TODO: Optionally make a separate call to /user/emails if email scope granted
		// Example: Fetch emails if primary email is null
		// emailsResp, emailsErr := client.Get("https://api.github.com/user/emails")
		// Handle emailsResp and emailsErr... find primary email
	}

	// Extract necessary fields
	userInfo := &UserInfo{
		Provider:     "github",
		ProviderID:   fmt.Sprintf("%.0f", ghUserInfo["id"]), // GitHub ID is often a number
		Email:        email,
		Name:         fmt.Sprintf("%v", ghUserInfo["name"]), // May be empty
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		RawData:      ghUserInfo,
	}

	// Call AuthService to handle login or registration
	tokenString, err := s.authService.LoginOrRegisterViaProvider(ctx, userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to login or register via github: %w", err)
	}

	return tokenString, nil
}
