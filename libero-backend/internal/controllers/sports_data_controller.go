package controllers

import (
	"encoding/json"
	"fmt"
	"libero-backend/internal/models"
	"libero-backend/internal/repository"
	"libero-backend/internal/service"
	"libero-backend/internal/utils"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
)

// SportsDataController handles HTTP requests for sports data.
type SportsDataController struct {
	mlService       service.MLService
	fixturesService service.FixturesService
	footballService *service.FootballService
	cacheRepo       repository.CacheRepository
}

// NewSportsDataController creates a new sports data controller instance.
func NewSportsDataController(
	mlService service.MLService,
	fixturesService service.FixturesService,
	footballService *service.FootballService,
	cacheRepo repository.CacheRepository,
) *SportsDataController {
	return &SportsDataController{
		mlService:       mlService,
		fixturesService: fixturesService,
		footballService: footballService,
		cacheRepo:       cacheRepo,
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

	cacheKey := fmt.Sprintf("standings_%s", competition)
	if cachedData, err := c.cacheRepo.Get(cacheKey); err == nil && cachedData != nil {
		var standings interface{}
		if err := json.Unmarshal(cachedData.Value, &standings); err == nil {
			w.Header().Set("ETag", cachedData.ETag)
			w.Header().Set("Last-Modified", cachedData.LastModified.Format(http.TimeFormat))
			utils.RespondWithJSON(w, http.StatusOK, standings)
			return
		}
	}

	standings, err := c.footballService.GetStandings(competition)
	if err != nil {
		fmt.Printf("Error fetching standings for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve standings data", http.StatusInternalServerError)
		return
	}

	newETag := fmt.Sprintf("%d", time.Now().UnixNano())
	if standingsJSON, err := json.Marshal(standings); err == nil {
		cacheItem := models.CacheItem{
			Key:          cacheKey,
			Value:        standingsJSON,
			ETag:         newETag,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		_ = c.cacheRepo.SetWithMetadata(cacheKey, cacheItem)
	}
	w.Header().Set("ETag", newETag)
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	utils.RespondWithJSON(w, http.StatusOK, standings)
}

// HandleGetTopScorers handles requests for top scorers
func (c *SportsDataController) HandleGetTopScorers(w http.ResponseWriter, r *http.Request) {
	competition := strings.ToUpper(r.URL.Query().Get("competition"))
	if competition == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}
	cacheKey := fmt.Sprintf("scorers_%s", competition)
	if cachedItem, err := c.cacheRepo.Get(cacheKey); err == nil && cachedItem != nil {
		var scorers interface{}
		if err := json.Unmarshal(cachedItem.Value, &scorers); err == nil {
			w.Header().Set("ETag", cachedItem.ETag)
			w.Header().Set("Last-Modified", cachedItem.LastModified.Format(http.TimeFormat))
			utils.RespondWithJSON(w, http.StatusOK, scorers)
			return
		}
	}
	scorers, err := c.footballService.GetTopScorers(competition)
	if err != nil {
		fmt.Printf("Error fetching top scorers for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve top scorers data", http.StatusInternalServerError)
		return
	}
	newETag := fmt.Sprintf("%d", time.Now().UnixNano())
	if scorersJSON, err := json.Marshal(scorers); err == nil {
		cacheItem := models.CacheItem{
			Key:          cacheKey,
			Value:        scorersJSON,
			ETag:         newETag,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		_ = c.cacheRepo.SetWithMetadata(cacheKey, cacheItem)
	}
	w.Header().Set("ETag", newETag)
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	utils.RespondWithJSON(w, http.StatusOK, scorers)
}

// HandleGetTodaysFixtures handles requests for today's fixtures.
func (c *SportsDataController) HandleGetTodaysFixtures(w http.ResponseWriter, r *http.Request) {
	cacheKey := fmt.Sprintf("todays_fixtures_%s", time.Now().Format("2006-01-02"))
	if cachedItem, err := c.cacheRepo.Get(cacheKey); err == nil && cachedItem != nil {
		var fixtures interface{}
		if err := json.Unmarshal(cachedItem.Value, &fixtures); err == nil {
			w.Header().Set("ETag", cachedItem.ETag)
			w.Header().Set("Last-Modified", cachedItem.LastModified.Format(http.TimeFormat))
			utils.RespondWithJSON(w, http.StatusOK, fixtures)
			return
		}
	}
	fixtures, err := c.fixturesService.GetTodaysFixtures()
	if err != nil {
		fmt.Printf("Error fetching today's fixtures: %v\n", err)
		http.Error(w, "Failed to retrieve today's fixtures", http.StatusInternalServerError)
		return
	}
	newETag := fmt.Sprintf("%d", time.Now().UnixNano())
	if fixturesJSON, err := json.Marshal(fixtures); err == nil {
		cacheItem := models.CacheItem{
			Key:          cacheKey,
			Value:        fixturesJSON,
			ETag:         newETag,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		_ = c.cacheRepo.SetWithMetadata(cacheKey, cacheItem)
	}
	w.Header().Set("ETag", newETag)
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	utils.RespondWithJSON(w, http.StatusOK, fixtures)
}

// HandleGetFixturesSummary handles requests for fixtures summary
func (c *SportsDataController) HandleGetFixturesSummary(w http.ResponseWriter, r *http.Request) {
	competition := r.URL.Query().Get("competition")
	if competition == "" {
		http.Error(w, "competition code is required", http.StatusBadRequest)
		return
	}
	cacheKey := fmt.Sprintf("fixtures_summary_%s", competition)
	if cachedItem, err := c.cacheRepo.Get(cacheKey); err == nil && cachedItem != nil {
		var summary interface{}
		if err := json.Unmarshal(cachedItem.Value, &summary); err == nil {
			w.Header().Set("ETag", cachedItem.ETag)
			w.Header().Set("Last-Modified", cachedItem.LastModified.Format(http.TimeFormat))
			utils.RespondWithJSON(w, http.StatusOK, summary)
			return
		}
	}
	summary, err := c.fixturesService.GetFixturesSummary(competition)
	if err != nil {
		fmt.Printf("Error fetching fixtures summary for %s: %v\n", competition, err)
		http.Error(w, "Failed to retrieve fixtures summary", http.StatusInternalServerError)
		return
	}
	newETag := fmt.Sprintf("%d", time.Now().UnixNano())
	if summaryJSON, err := json.Marshal(summary); err == nil {
		cacheItem := models.CacheItem{
			Key:          cacheKey,
			Value:        summaryJSON,
			ETag:         newETag,
			LastModified: time.Now(),
			ExpiresAt:    time.Now().Add(24 * time.Hour),
		}
		_ = c.cacheRepo.SetWithMetadata(cacheKey, cacheItem)
	}
	w.Header().Set("ETag", newETag)
	w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))
	utils.RespondWithJSON(w, http.StatusOK, summary)
}

// CacheMiddleware is a middleware for caching HTTP responses.
func (c *SportsDataController) CacheMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Check if the response is cached
		cacheKey := r.URL.Path
		cachedResponse, err := c.cacheRepo.Get(cacheKey)
		if err != nil {
			fmt.Printf("Error fetching from cache: %v\n", err)
		}

		if cachedResponse != nil {
			// Serve the cached response
			fmt.Println("Serving from cache:", cacheKey)
			w.Header().Set("Content-Type", "application/json")
			w.Header().Set("ETag", cachedResponse.ETag)
			w.Header().Set("Last-Modified", cachedResponse.LastModified.Format(http.TimeFormat))
			w.WriteHeader(http.StatusOK)
			w.Write(cachedResponse.Value)
			return
		}

		// If not cached, proceed with the request
		next.ServeHTTP(w, r)

		// Cache the response
		go func() {
			// Serialize the response body
			var responseBody []byte
			if r.Method == http.MethodGet {
				responseBody, err = json.Marshal(w)
				if err != nil {
					fmt.Printf("Error serializing response body: %v\n", err)
					return
				}

				// Store in cache with expiration
				err = c.cacheRepo.Set(cacheKey, responseBody, 10*time.Minute)
				if err != nil {
					fmt.Printf("Error storing in cache: %v\n", err)
				}
			}
		}()
	})
}
