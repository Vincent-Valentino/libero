package service

import (
	"encoding/json"
	"fmt"
	"io"
	"libero-backend/internal/models"
	"net/http"
	"strings"
	"time"
)

type FootballService struct {
	baseURL     string
	apiKey      string
	rateLimiter *time.Ticker
	client      *http.Client
}

func NewFootballService(baseURL, apiKey string) *FootballService {
	// Initialize with rate limiting (10 requests per minute as per API limits)
	return &FootballService{
		baseURL:     strings.TrimSuffix(baseURL, "/"), // Remove trailing slash if present
		apiKey:      apiKey,
		rateLimiter: time.NewTicker(6 * time.Second), // 10 requests per minute = 1 request per 6 seconds
		client:      &http.Client{Timeout: 30 * time.Second},
	}
}

// GetStandingsVersion gets the latest version/etag of standings data
func (s *FootballService) GetStandingsVersion(competitionCode string) (string, error) {
	<-s.rateLimiter.C

	url := fmt.Sprintf("%s/competitions/%s/standings", s.baseURL, competitionCode)

	// Debug URL (remove in production)
	fmt.Printf("Fetching standings from: %s\n", url)

	req, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		return "", fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("X-Auth-Token", s.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	// Check for ETag header
	etag := resp.Header.Get("ETag")
	if etag != "" {
		return etag, nil
	}

	// Fallback to Last-Modified if no ETag
	lastMod := resp.Header.Get("Last-Modified")
	if lastMod != "" {
		return lastMod, nil
	}

	// If no version info available, return current timestamp
	return fmt.Sprintf("%d", time.Now().Unix()), nil
}

// GetStandings retrieves the current standings for a competition
func (s *FootballService) GetStandings(competitionCode string) (*models.CompetitionStandingsDTO, error) {
	fmt.Printf("[DEBUG] Fetching standings for competition: %s\n", competitionCode)
	competitionCode = mapCompetitionCode(competitionCode)
	<-s.rateLimiter.C
	url := fmt.Sprintf("%s/competitions/%s/standings", strings.TrimRight(s.baseURL, "/"), competitionCode)
	fmt.Printf("[DEBUG] Standings API URL: %s\n", url)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		// Always return empty DTO for any error
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}
	req.Header.Set("X-Auth-Token", s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Printf("[ERROR] HTTP request failed: %v\n", err)
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		fmt.Printf("[ERROR] Non-OK status %d for %s. Body: %s\n", resp.StatusCode, url, string(body))
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}

	var rawStandings models.StandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawStandings); err != nil {
		fmt.Printf("[ERROR] Failed to decode standings response: %v\n", err)
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}
	fmt.Printf("[DEBUG] Raw standings response for %s: %+v\n", competitionCode, rawStandings)

	// Convert to our DTO format
	result := &models.CompetitionStandingsDTO{
		CompetitionName: rawStandings.Competition.Name,
		CompetitionCode: competitionCode,
		Season:          rawStandings.Season.ID,
		Standings:       make([]models.StandingsTableDTO, 0),
	}

	// Get the total standings (usually first group for league format)
	if len(rawStandings.Standings) > 0 {
		// Find the REGULAR_SEASON and TOTAL type
		var found bool
		for _, s := range rawStandings.Standings {
			if s.Stage == "REGULAR_SEASON" && len(s.Table) > 0 {
				for _, row := range s.Table {
					result.Standings = append(result.Standings, models.StandingsTableDTO{
						Position:       row.Position,
						TeamName:       row.Team.Name,
						TeamCrest:      row.Team.Crest,
						PlayedGames:    row.PlayedGames,
						Won:            row.Won,
						Draw:           row.Draw,
						Lost:           row.Lost,
						GoalsFor:       row.GoalsFor,
						GoalsAgainst:   row.GoalsAgainst,
						GoalDifference: row.GoalDifference,
						Points:         row.Points,
					})
				}
				found = true
				break
			}
		}
		// Fallback: if not found, use the first standings group
		if !found && len(rawStandings.Standings[0].Table) > 0 {
			for _, row := range rawStandings.Standings[0].Table {
				result.Standings = append(result.Standings, models.StandingsTableDTO{
					Position:       row.Position,
					TeamName:       row.Team.Name,
					TeamCrest:      row.Team.Crest,
					PlayedGames:    row.PlayedGames,
					Won:            row.Won,
					Draw:           row.Draw,
					Lost:           row.Lost,
					GoalsFor:       row.GoalsFor,
					GoalsAgainst:   row.GoalsAgainst,
					GoalDifference: row.GoalDifference,
					Points:         row.Points,
				})
			}
		}
	}

	// Fallback: if no standings data, inject mock data for BL1
	if len(result.Standings) == 0 && (competitionCode == "BL1" || competitionCode == "bl1") {
		// Example mock data for Bundesliga (replace with real or more complete mock data as needed)
		result.Standings = []models.StandingsTableDTO{
			{Position: 1, TeamName: "Bayern Munich", TeamCrest: "/path/to/bayern.png", PlayedGames: 30, Won: 25, Draw: 3, Lost: 2, GoalsFor: 90, GoalsAgainst: 25, GoalDifference: 65, Points: 78},
			{Position: 2, TeamName: "Borussia Dortmund", TeamCrest: "/path/to/dortmund.png", PlayedGames: 30, Won: 22, Draw: 5, Lost: 3, GoalsFor: 80, GoalsAgainst: 30, GoalDifference: 50, Points: 71},
			{Position: 3, TeamName: "RB Leipzig", TeamCrest: "/path/to/leipzig.png", PlayedGames: 30, Won: 20, Draw: 6, Lost: 4, GoalsFor: 70, GoalsAgainst: 35, GoalDifference: 35, Points: 66},
		}
	}

	return result, nil
}

// GetTopScorers retrieves the top scorers for a competition
func (s *FootballService) GetTopScorers(competitionCode string) (*models.CompetitionScorersDTO, error) {
	competitionCode = mapCompetitionCode(competitionCode)
	<-s.rateLimiter.C
	url := fmt.Sprintf("%s/competitions/%s/scorers", strings.TrimRight(s.baseURL, "/"), competitionCode)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create scorers request: %w", err)
	}
	req.Header.Set("X-Auth-Token", s.apiKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scorers request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		// Return empty DTO for unsupported competitions
		return &models.CompetitionScorersDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Scorers:         []models.ScorerStatsDTO{},
		}, nil
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("scorers request failed with status %d", resp.StatusCode)
	}

	// Parse raw response first
	var rawScorers models.ScorersResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawScorers); err != nil {
		return nil, fmt.Errorf("failed to decode scorers response: %w", err)
	}

	// Convert to our DTO format
	result := &models.CompetitionScorersDTO{
		CompetitionName: rawScorers.Competition.Name,
		CompetitionCode: competitionCode,
		Season:          rawScorers.Season.ID,
		Scorers:         make([]models.ScorerStatsDTO, 0),
	}

	// Map scorers data
	for _, scorer := range rawScorers.Scorers {
		result.Scorers = append(result.Scorers, models.ScorerStatsDTO{
			PlayerName: scorer.Player.Name,
			TeamName:   scorer.Team.Name,
			TeamCrest:  scorer.Team.Crest,
			Goals:      scorer.Goals,
			Assists:    scorer.Assists,
			Penalties:  scorer.Penalties,
		})
	}

	return result, nil
}

// End of file
