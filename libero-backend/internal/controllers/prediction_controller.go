package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"libero-backend/config"
	"libero-backend/internal/utils"
	"net/http"
)

// PredictionController handles HTTP requests for match predictions
type PredictionController struct {
	mlServiceURL string
	httpClient   *http.Client
}

// NewPredictionController creates a new prediction controller instance
func NewPredictionController(cfg *config.Config) *PredictionController {
	return &PredictionController{
		mlServiceURL: cfg.MLServiceURL,
		httpClient:   &http.Client{},
	}
}

// PredictMatchRequest represents the request payload for match prediction
type PredictMatchRequest struct {
	League   string `json:"league" binding:"required"`
	HomeTeam string `json:"home_team" binding:"required"`
	AwayTeam string `json:"away_team" binding:"required"`
}

// PredictMatchResponse represents the response from the ML service
type PredictMatchResponse struct {
	Prediction          int                `json:"prediction"`
	Probabilities       map[string]float64 `json:"probabilities"`
	ExpectedHomeGoals   float64            `json:"expected_home_goals"`
	ExpectedAwayGoals   float64            `json:"expected_away_goals"`
	MostLikelyHomeScore int                `json:"most_likely_home_score"`
	MostLikelyAwayScore int                `json:"most_likely_away_score"`
}

// PredictMatch handles match prediction requests
func (c *PredictionController) PredictMatch(w http.ResponseWriter, r *http.Request) {
	var request PredictMatchRequest

	// Parse request body
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	// Validate required fields
	if request.League == "" || request.HomeTeam == "" || request.AwayTeam == "" {
		http.Error(w, "League, home_team, and away_team are required", http.StatusBadRequest)
		return
	}

	// Call the FastAPI ML service
	prediction, err := c.callMLService(request)
	if err != nil {
		fmt.Printf("Error calling ML service: %v\n", err)
		http.Error(w, "Failed to get prediction", http.StatusInternalServerError)
		return
	}

	// Return the prediction result
	utils.RespondWithJSON(w, http.StatusOK, prediction)
}

// callMLService makes an HTTP call to the FastAPI ML service
func (c *PredictionController) callMLService(request PredictMatchRequest) (*PredictMatchResponse, error) {
	// Prepare the request payload for the ML service
	payload := map[string]string{
		"league":    request.League,
		"home_team": request.HomeTeam,
		"away_team": request.AwayTeam,
	}

	// Convert to JSON
	jsonData, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create HTTP request to ML service
	url := c.mlServiceURL + "/predict"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")

	// Make the request
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to call ML service: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ML service returned status %d", resp.StatusCode)
	}

	// Parse response
	var prediction PredictMatchResponse
	if err := json.NewDecoder(resp.Body).Decode(&prediction); err != nil {
		return nil, fmt.Errorf("failed to decode ML service response: %w", err)
	}

	return &prediction, nil
}

// GetAvailableTeams fetches available teams from the ML service
func (c *PredictionController) GetAvailableTeams(w http.ResponseWriter, r *http.Request) {
	url := c.mlServiceURL + "/teams"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching teams from ML service: %v\n", err)
		http.Error(w, "Failed to get teams", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "ML service unavailable", http.StatusServiceUnavailable)
		return
	}

	// FastAPI returns {"teams": [...]} so we need to decode the wrapper object
	var teamsResponse struct {
		Teams []string `json:"teams"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&teamsResponse); err != nil {
		fmt.Printf("Error decoding teams response: %v\n", err)
		http.Error(w, "Failed to parse teams data", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"teams": teamsResponse.Teams,
	})
}

// GetAvailableLeagues fetches available leagues from the ML service
func (c *PredictionController) GetAvailableLeagues(w http.ResponseWriter, r *http.Request) {
	url := c.mlServiceURL + "/leagues"

	resp, err := c.httpClient.Get(url)
	if err != nil {
		fmt.Printf("Error fetching leagues from ML service: %v\n", err)
		http.Error(w, "Failed to get leagues", http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		http.Error(w, "ML service unavailable", http.StatusServiceUnavailable)
		return
	}

	// FastAPI returns {"leagues": [...]} so we need to decode the wrapper object
	var leaguesResponse struct {
		Leagues []string `json:"leagues"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&leaguesResponse); err != nil {
		fmt.Printf("Error decoding leagues response: %v\n", err)
		http.Error(w, "Failed to parse leagues data", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, map[string]interface{}{
		"leagues": leaguesResponse.Leagues,
	})
}
