package service

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json" // Added for potential user info parsing
	"errors"        // Added for error handling
	"fmt"
	"io"
	"net/http"

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
	// Debug OAuth configuration
	fmt.Printf("🔧 OAuth Configuration Debug:\n")
	fmt.Printf("  Google Client ID: %s (len: %d)\n", maskSecret(cfg.Google.ClientID), len(cfg.Google.ClientID))
	fmt.Printf("  Google Client Secret: %s (len: %d)\n", maskSecret(cfg.Google.ClientSecret), len(cfg.Google.ClientSecret))
	fmt.Printf("  Google Redirect URL: %s\n", cfg.Google.RedirectURL)
	fmt.Printf("  Facebook Client ID: %s (len: %d)\n", maskSecret(cfg.Facebook.ClientID), len(cfg.Facebook.ClientID))
	fmt.Printf("  Facebook Client Secret: %s (len: %d)\n", maskSecret(cfg.Facebook.ClientSecret), len(cfg.Facebook.ClientSecret))
	fmt.Printf("  Facebook Redirect URL: %s\n", cfg.Facebook.RedirectURL)
	fmt.Printf("  GitHub Client ID: %s (len: %d)\n", maskSecret(cfg.GitHub.ClientID), len(cfg.GitHub.ClientID))
	fmt.Printf("  GitHub Client Secret: %s (len: %d)\n", maskSecret(cfg.GitHub.ClientSecret), len(cfg.GitHub.ClientSecret))
	fmt.Printf("  GitHub Redirect URL: %s\n", cfg.GitHub.RedirectURL)

	// Validate OAuth configurations
	if cfg.Google.ClientID == "" || cfg.Google.ClientSecret == "" {
		fmt.Printf("⚠️  WARNING: Google OAuth is not properly configured (missing ClientID or ClientSecret)\n")
	}
	if cfg.GitHub.ClientID == "" || cfg.GitHub.ClientSecret == "" {
		fmt.Printf("⚠️  WARNING: GitHub OAuth is not properly configured (missing ClientID or ClientSecret)\n")
	}
	if cfg.Facebook.ClientID != "" && cfg.Facebook.ClientSecret != "" {
		fmt.Printf("✅ Facebook OAuth is properly configured\n")
	}

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

// maskSecret masks a secret string for logging, showing only first/last chars
func maskSecret(secret string) string {
	if len(secret) <= 8 {
		return "***"
	}
	return secret[:4] + "***" + secret[len(secret)-4:]
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
	// Validate configuration before proceeding
	if s.GoogleConfig.ClientID == "" || s.GoogleConfig.ClientSecret == "" {
		fmt.Printf("ERROR: Google OAuth not configured - missing ClientID or ClientSecret\n")
		return "", ""
	}

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

	// Fetch user info from Google API - using the current v2 userinfo endpoint
	client := s.GoogleConfig.Client(ctx, token)
	response, err := client.Get("https://www.googleapis.com/oauth2/v2/userinfo?access_token=" + token.AccessToken)
	if err != nil {
		return "", fmt.Errorf("failed getting user info: %w", err)
	}
	defer response.Body.Close()

	// Check response status
	if response.StatusCode != http.StatusOK {
		return "", fmt.Errorf("google userinfo API returned status %d", response.StatusCode)
	}

	contents, err := io.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading user info response body: %w", err)
	}

	var googleUserInfo map[string]interface{}
	if err := json.Unmarshal(contents, &googleUserInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal google user info: %w", err)
	}

	// Validate that we have required fields
	if googleUserInfo["id"] == nil || googleUserInfo["email"] == nil {
		return "", fmt.Errorf("google response missing required fields (id or email)")
	}

	// Extract necessary fields with better type checking
	userInfo := &UserInfo{
		Provider:     "google",
		ProviderID:   fmt.Sprintf("%v", googleUserInfo["id"]),
		Email:        fmt.Sprintf("%v", googleUserInfo["email"]),
		Name:         fmt.Sprintf("%v", googleUserInfo["name"]),
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		RawData:      googleUserInfo,
	}

	// Validate email is not empty
	if userInfo.Email == "" || userInfo.Email == "<nil>" {
		return "", fmt.Errorf("google did not provide a valid email address")
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
	// Validate configuration before proceeding
	if s.GitHubConfig.ClientID == "" || s.GitHubConfig.ClientSecret == "" {
		fmt.Printf("ERROR: GitHub OAuth not configured - missing ClientID or ClientSecret\n")
		return "", ""
	}

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

	// Check response status
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("github user API returned status %d", resp.StatusCode)
	}

	contents, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed reading github user info response body: %w", err)
	}

	var ghUserInfo map[string]interface{}
	if err := json.Unmarshal(contents, &ghUserInfo); err != nil {
		return "", fmt.Errorf("failed to unmarshal github user info: %w", err)
	}

	// Validate that we have required fields
	if ghUserInfo["id"] == nil {
		return "", fmt.Errorf("github response missing required field (id)")
	}

	// GitHub might not return email directly from /user, may need /user/emails
	email := ""
	if val, ok := ghUserInfo["email"].(string); ok && val != "" {
		email = val
	} else {
		// Fetch primary email from /user/emails endpoint
		emailsResp, emailsErr := client.Get("https://api.github.com/user/emails")
		if emailsErr == nil {
			defer emailsResp.Body.Close()
			if emailsResp.StatusCode == http.StatusOK {
				emailsContents, emailsReadErr := io.ReadAll(emailsResp.Body)
				if emailsReadErr == nil {
					var emails []map[string]interface{}
					if json.Unmarshal(emailsContents, &emails) == nil {
						// Find primary email
						for _, emailObj := range emails {
							if isPrimary, ok := emailObj["primary"].(bool); ok && isPrimary {
								if primaryEmail, ok := emailObj["email"].(string); ok && primaryEmail != "" {
									email = primaryEmail
									break
								}
							}
						}
						// If no primary found, use the first verified email
						if email == "" {
							for _, emailObj := range emails {
								if isVerified, ok := emailObj["verified"].(bool); ok && isVerified {
									if verifiedEmail, ok := emailObj["email"].(string); ok && verifiedEmail != "" {
										email = verifiedEmail
										break
									}
								}
							}
						}
					}
				}
			}
		}
	}

	// GitHub ID can be a number, handle it properly
	var providerID string
	switch id := ghUserInfo["id"].(type) {
	case float64:
		providerID = fmt.Sprintf("%.0f", id)
	case int64:
		providerID = fmt.Sprintf("%d", id)
	case int:
		providerID = fmt.Sprintf("%d", id)
	default:
		providerID = fmt.Sprintf("%v", id)
	}

	// Extract name with fallback to login
	name := ""
	if nameVal, ok := ghUserInfo["name"].(string); ok && nameVal != "" {
		name = nameVal
	} else if login, ok := ghUserInfo["login"].(string); ok && login != "" {
		name = login // Use GitHub username as fallback
	}

	// Extract necessary fields
	userInfo := &UserInfo{
		Provider:     "github",
		ProviderID:   providerID,
		Email:        email,
		Name:         name,
		AccessToken:  token.AccessToken,
		RefreshToken: token.RefreshToken,
		RawData:      ghUserInfo,
	}

	// Note: GitHub OAuth might not always provide email if user's email privacy settings
	// don't allow it. In such cases, we'll proceed without email but it might cause issues
	// if the application requires email for user creation.
	if userInfo.Email == "" {
		// Log this case for debugging but don't fail - some users have private emails
		// In production, you might want to request additional permissions or handle this case
		fmt.Printf("WARNING: GitHub user %s (%s) did not provide email address\n", userInfo.Name, userInfo.ProviderID)
	}

	// Call AuthService to handle login or registration
	tokenString, err := s.authService.LoginOrRegisterViaProvider(ctx, userInfo)
	if err != nil {
		return "", fmt.Errorf("failed to login or register via github: %w", err)
	}

	return tokenString, nil
}
