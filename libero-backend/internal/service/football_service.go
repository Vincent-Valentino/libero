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
	competitionCode = mapCompetitionCode(competitionCode)
	<-s.rateLimiter.C
	url := fmt.Sprintf("%s/competitions/%s/standings", strings.TrimRight(s.baseURL, "/"), competitionCode)
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
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Always return empty DTO for any non-OK status
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
	}

	// Parse raw response first
	var rawStandings models.StandingsResponse
	if err := json.NewDecoder(resp.Body).Decode(&rawStandings); err != nil {
		return &models.CompetitionStandingsDTO{
			CompetitionName: "",
			CompetitionCode: competitionCode,
			Season:          0,
			Standings:       []models.StandingsTableDTO{},
		}, nil
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
