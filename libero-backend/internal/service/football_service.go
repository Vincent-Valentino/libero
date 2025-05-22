package service

import (
	"encoding/json"
	"fmt"
	"libero-backend/internal/models"
	"net/http"
	"strings"
	"time"
)

type FootballService struct {
	baseURL     string
	apiKey      string
	rateLimiter *time.Ticker
}

func NewFootballService(baseURL, apiKey string) *FootballService {
	// Initialize with rate limiting (3 requests per second as per provider limits)
	return &FootballService{
		baseURL:     baseURL,
		apiKey:      apiKey,
		rateLimiter: time.NewTicker(time.Second / 3),
	}
}

// GetStandings retrieves the current standings for a competition
func (s *FootballService) GetStandings(competitionCode string) (*models.CompetitionStandingsDTO, error) {
	// Wait for rate limiter before making request
	<-s.rateLimiter.C

	// Build URL for standings endpoint
	url := fmt.Sprintf("%s/competitions/%s/standings", strings.TrimRight(s.baseURL, "/"), competitionCode)

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create standings request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Auth-Token", s.apiKey)

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute standings request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("standings request failed with status %d", resp.StatusCode)
	}

	// Parse raw response first
	var rawStandings models.StandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawStandings); err != nil {
		return nil, fmt.Errorf("failed to decode standings response: %w", err)
	}

	// Convert to our DTO format
	result := &models.CompetitionStandingsDTO{
		CompetitionName: rawStandings.Competition.Name,
		CompetitionCode: competitionCode,
		Season:          rawStandings.Season.ID,
		Standings:       make([]models.StandingsTableDTO, 0),
	}

	// Get the total standings (usually first group for league format)
	if len(rawStandings.Standings) > 0 {
		var totalStandings struct {
			Stage string `json:"stage"`
			Table []struct {
				Position       int                 `json:"position"`
				Team           models.TeamResponse `json:"team"`
				PlayedGames    int                 `json:"playedGames"`
				Won            int                 `json:"won"`
				Draw           int                 `json:"draw"`
				Lost           int                 `json:"lost"`
				Points         int                 `json:"points"`
				GoalsFor       int                 `json:"goalsFor"`
				GoalsAgainst   int                 `json:"goalsAgainst"`
				GoalDifference int                 `json:"goalDifference"`
			} `json:"table"`
		}

		// Find regular season standings
		for _, s := range rawStandings.Standings {
			if s.Stage == "REGULAR_SEASON" {
				totalStandings = s
				break
			}
		}
		// If no regular season found, use first standings group
		if totalStandings.Table == nil && len(rawStandings.Standings) > 0 {
			totalStandings = rawStandings.Standings[0]
		}

		// Map standings data
		for _, row := range totalStandings.Table {
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

	return result, nil
}

// GetTopScorers retrieves the top scorers for a competition
func (s *FootballService) GetTopScorers(competitionCode string) (*models.CompetitionScorersDTO, error) {
	// Wait for rate limiter before making request
	<-s.rateLimiter.C

	// Build URL for scorers endpoint
	url := fmt.Sprintf("%s/competitions/%s/scorers", strings.TrimRight(s.baseURL, "/"), competitionCode)

	// Create request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create scorers request: %w", err)
	}

	// Set headers
	req.Header.Set("X-Auth-Token", s.apiKey)

	// Perform request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute scorers request: %w", err)
	}
	defer resp.Body.Close()

	// Handle non-200 responses
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
