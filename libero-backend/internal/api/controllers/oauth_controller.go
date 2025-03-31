package controllers

import (
	"crypto/rand"
	"encoding/base64"
	// "encoding/json" // Removed previously
	// "fmt" // Removing fmt
	"libero-backend/internal/service"
	"net/http"
	// Keep necessary imports like net/http, encoding/base64, service
)

// OAuthController handles OAuth authentication requests
type OAuthController struct {
	oauthService service.OAuthService // Only depends on the OAuthService interface
	// userService is removed as the logic is handled within OAuthService/AuthService
}

// NewOAuthController creates a new OAuthController instance
func NewOAuthController(oauthService service.OAuthService) *OAuthController { // Only takes OAuthService
	return &OAuthController{
		oauthService: oauthService,
		// userService removed
	}
}

// generateStateOauthCookie generates a random state string and sets it as a cookie
// (Kept as it's relevant to the controller's HTTP handling)
func generateStateOauthCookie(w http.ResponseWriter) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	// Set cookie with SameSite=Lax for security
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate",
		Value:    state,
		MaxAge:   3600, // 1 hour
		Path:     "/",
		HttpOnly: true,
		// Secure: true, // Enable in production with HTTPS
		SameSite: http.SameSiteLaxMode,
	})
	return state
}

// --- Google Handlers ---

// GoogleLogin initiates the Google OAuth flow
func (ctrl *OAuthController) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	url, state := ctrl.oauthService.GetGoogleLoginURL()
	http.SetCookie(w, &http.Cookie{
		Name:     "oauthstate", Value: state, MaxAge: 3600, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GoogleCallback handles the callback from Google OAuth
func (ctrl *OAuthController) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	oauthStateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "invalid session: state cookie missing", http.StatusBadRequest)
		return
	}

	receivedState := r.FormValue("state")
	if receivedState != oauthStateCookie.Value {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}

	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "authorization code missing", http.StatusBadRequest)
		return
	}

	tokenString, err := ctrl.oauthService.HandleGoogleCallback(r.Context(), oauthStateCookie.Value, receivedState, code)
	if err != nil {
		// TODO: Replace with proper logging
		http.Error(w, "Authentication failed.", http.StatusInternalServerError)
		return
	}

	// Return token
	// Assuming respondWithJSON is available from user_controller.go in the same package
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// --- Facebook Handlers ---

// FacebookLogin initiates the Facebook OAuth flow
func (ctrl *OAuthController) FacebookLogin(w http.ResponseWriter, r *http.Request) {
	url, state := ctrl.oauthService.GetFacebookLoginURL()
	http.SetCookie(w, &http.Cookie{
		Name: "oauthstate", Value: state, MaxAge: 3600, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// FacebookCallback handles the callback from Facebook OAuth
func (ctrl *OAuthController) FacebookCallback(w http.ResponseWriter, r *http.Request) {
	oauthStateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "invalid session: state cookie missing", http.StatusBadRequest)
		return
	}
	receivedState := r.FormValue("state")
	if receivedState != oauthStateCookie.Value {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "authorization code missing", http.StatusBadRequest)
		return
	}

	tokenString, err := ctrl.oauthService.HandleFacebookCallback(r.Context(), oauthStateCookie.Value, receivedState, code)
	if err != nil {
		// TODO: Replace with proper logging
		http.Error(w, "Authentication failed.", http.StatusInternalServerError)
		return
	}
	// Assuming respondWithJSON is available from user_controller.go in the same package
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// --- GitHub Handlers ---

// GitHubLogin initiates the GitHub OAuth flow
func (ctrl *OAuthController) GitHubLogin(w http.ResponseWriter, r *http.Request) {
	url, state := ctrl.oauthService.GetGitHubLoginURL()
	http.SetCookie(w, &http.Cookie{
		Name: "oauthstate", Value: state, MaxAge: 3600, Path: "/", HttpOnly: true, SameSite: http.SameSiteLaxMode,
	})
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// GitHubCallback handles the callback from GitHub OAuth
func (ctrl *OAuthController) GitHubCallback(w http.ResponseWriter, r *http.Request) {
	oauthStateCookie, err := r.Cookie("oauthstate")
	if err != nil {
		http.Error(w, "invalid session: state cookie missing", http.StatusBadRequest)
		return
	}
	receivedState := r.FormValue("state")
	if receivedState != oauthStateCookie.Value {
		http.Error(w, "invalid oauth state", http.StatusBadRequest)
		return
	}
	code := r.FormValue("code")
	if code == "" {
		http.Error(w, "authorization code missing", http.StatusBadRequest)
		return
	}

	tokenString, err := ctrl.oauthService.HandleGitHubCallback(r.Context(), oauthStateCookie.Value, receivedState, code)
	if err != nil {
		// TODO: Replace with proper logging
		http.Error(w, "Authentication failed.", http.StatusInternalServerError)
		return
	}
	// Assuming respondWithJSON is available from user_controller.go in the same package
	respondWithJSON(w, http.StatusOK, map[string]string{"token": tokenString})
}

// Removed helper functions and unused imports.