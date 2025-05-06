package controllers

import (
	"fmt"
	"libero-backend/internal/service"
	"net/http"

	"github.com/gorilla/mux"
)

// SportsDataController handles HTTP requests for sports data.
type SportsDataController struct {
	mlService       *service.MLService      // Use pointer for service dependency
	fixturesService service.FixturesService // Add FixturesService dependency
}

// NewSportsDataController creates a new sports data controller instance.
func NewSportsDataController(mlService *service.MLService, fixturesService service.FixturesService) *SportsDataController {
	return &SportsDataController{
		mlService:       mlService,
		fixturesService: fixturesService,
	}
}

// HandleGetUpcomingMatches handles requests for upcoming matches.
func (c *SportsDataController) HandleGetUpcomingMatches(w http.ResponseWriter, r *http.Request) {
	matches, err := c.mlService.GetUpcomingMatches()
	if err != nil {
		// Log the error server-side (replace with proper logging)
		fmt.Printf("Error fetching upcoming matches from ML service: %v\n", err)
		http.Error(w, "Failed to retrieve upcoming matches", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, matches)
}

// HandleGetResults handles requests for match results.
func (c *SportsDataController) HandleGetResults(w http.ResponseWriter, r *http.Request) {
	results, err := c.mlService.GetResults()
	if err != nil {
		// Log the error server-side
		fmt.Printf("Error fetching results from ML service: %v\n", err)
		http.Error(w, "Failed to retrieve match results", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, results)
}

// HandleGetPlayerStats handles requests for player statistics.
func (c *SportsDataController) HandleGetPlayerStats(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	playerID := vars["player_id"] // Extract player_id from URL path

	if playerID == "" {
		http.Error(w, "Player ID is required", http.StatusBadRequest)
		return
	}

	stats, err := c.mlService.GetPlayerStats(playerID)
	if err != nil {
		// Log the error server-side
		fmt.Printf("Error fetching player stats for ID %s from ML service: %v\n", playerID, err)
		http.Error(w, "Failed to retrieve player statistics", http.StatusInternalServerError)
		return
	}

	// Check if stats are nil (meaning not found by the service)
	if stats == nil {
		http.Error(w, fmt.Sprintf("Stats not found for player_id: %s", playerID), http.StatusNotFound)
		return
	}

	respondWithJSON(w, http.StatusOK, stats)
}

// HandleGetTodaysFixtures handles requests for today's fixtures.
func (c *SportsDataController) HandleGetTodaysFixtures(w http.ResponseWriter, r *http.Request) {
	fixtures, err := c.fixturesService.GetTodaysFixtures()
	if err != nil {
		// Log the error server-side
		fmt.Printf("Error fetching today's fixtures: %v\n", err)
		http.Error(w, "Failed to retrieve today's fixtures", http.StatusInternalServerError)
		return
	}

	respondWithJSON(w, http.StatusOK, fixtures)
}

// HandleGetFixturesSummary handles requests for a league summary (today/tomorrow/upcoming).
func (c *SportsDataController) HandleGetFixturesSummary(w http.ResponseWriter, r *http.Request) {
	// Expect query param: ?competition={code}
	code := r.URL.Query().Get("competition")
	if code == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}
	summary, err := c.fixturesService.GetFixturesSummary(code)
	if err != nil {
		fmt.Printf("Error fetching fixtures summary for %s: %v\n", code, err)
		http.Error(w, "Failed to retrieve fixtures summary", http.StatusInternalServerError)
		return
	}
	respondWithJSON(w, http.StatusOK, summary)
}

/*
// respondWithJSON is likely defined in another controller file within the same package.
// If not, uncomment and use this helper function.
func respondWithJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
*/