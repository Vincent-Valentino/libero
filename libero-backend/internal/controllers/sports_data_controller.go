package controllers

import (
	"fmt"
	"libero-backend/internal/service"
	"libero-backend/internal/utils"

	"net/http"
	"strings"

	"github.com/gorilla/mux"
)

// SportsDataController handles HTTP requests for sports data.
type SportsDataController struct {
	mlService       service.MLService
	fixturesService service.FixturesService
	footballService *service.FootballService
}

// NewSportsDataController creates a new sports data controller instance.
func NewSportsDataController(mlService service.MLService, fixturesService service.FixturesService, footballService *service.FootballService) *SportsDataController {
	return &SportsDataController{
		mlService:       mlService,
		fixturesService: fixturesService,
		footballService: footballService,
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

	utils.RespondWithJSON(w, http.StatusOK, matches)
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

	utils.RespondWithJSON(w, http.StatusOK, results)
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

	utils.RespondWithJSON(w, http.StatusOK, stats)
}

// HandleGetStandings handles requests for league standings
func (c *SportsDataController) HandleGetStandings(w http.ResponseWriter, r *http.Request) {
	competition := strings.ToUpper(r.URL.Query().Get("competition"))
	if competition == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}

	standings, err := c.footballService.GetStandings(competition)
	if err != nil {
		fmt.Printf("Error fetching standings for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve standings data", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, standings)
}

// HandleGetTopScorers handles requests for top scorers
func (c *SportsDataController) HandleGetTopScorers(w http.ResponseWriter, r *http.Request) {
	competition := strings.ToUpper(r.URL.Query().Get("competition"))
	if competition == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}

	scorers, err := c.footballService.GetTopScorers(competition)
	if err != nil {
		fmt.Printf("Error fetching top scorers for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve top scorers data", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, scorers)
}

// HandleGetTodaysFixtures handles requests for today's fixtures.
func (c *SportsDataController) HandleGetTodaysFixtures(w http.ResponseWriter, r *http.Request) {
	fixtures, err := c.fixturesService.GetTodaysFixtures()
	if err != nil {
		fmt.Printf("Error fetching today's fixtures: %v\n", err)
		http.Error(w, "Failed to retrieve today's fixtures", http.StatusInternalServerError)
		return
	}
	utils.RespondWithJSON(w, http.StatusOK, fixtures)
}

// HandleGetFixturesSummary handles requests for fixtures summary
func (c *SportsDataController) HandleGetFixturesSummary(w http.ResponseWriter, r *http.Request) {
	competition := r.URL.Query().Get("competition")
	if competition == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}

	summary, err := c.fixturesService.GetFixturesSummary(competition)
	if err != nil {
		fmt.Printf("Error fetching fixtures summary for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve fixtures summary", http.StatusInternalServerError)
		return
	}

	utils.RespondWithJSON(w, http.StatusOK, summary)
}
