package service

import (
	"encoding/json"
	"fmt"
	"libero-backend/config"
	"libero-backend/internal/api/dto"
	"net/http"
	"net/url"
	"time"
)

// MLService defines the interface for ML-related operations.
type MLService interface {
	GetUpcomingMatches() ([]dto.MatchDTO, error)
	GetResults() ([]dto.ResultDTO, error)
	GetPlayerStats(playerID string) (*dto.PlayerStatsDTO, error)
}

// mlService implements the MLService interface.
type mlService struct {
	baseURL    string
	httpClient *http.Client
}

// NewMLService creates a new MLService instance.
func NewMLService(cfg *config.Config) MLService {
	return &mlService{
		baseURL: cfg.MLServiceURL,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// GetUpcomingMatches fetches upcoming matches from the ML service.
func (s *mlService) GetUpcomingMatches() ([]dto.MatchDTO, error) {
	endpoint := "/matches/upcoming"
	targetURL := s.baseURL + endpoint

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", targetURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// Consider reading the body here for more detailed error messages from ml-service if needed
		return nil, fmt.Errorf("ml service returned non-OK status (%d) for %s", resp.StatusCode, targetURL)
	}

	var matches []dto.MatchDTO
	if err := json.NewDecoder(resp.Body).Decode(&matches); err != nil {
		return nil, fmt.Errorf("failed to decode upcoming matches response from %s: %w", targetURL, err)
	}

	return matches, nil
}

// GetResults fetches match results from the ML service.
func (s *mlService) GetResults() ([]dto.ResultDTO, error) {
	endpoint := "/matches/results"
	targetURL := s.baseURL + endpoint

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", targetURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service returned non-OK status (%d) for %s", resp.StatusCode, targetURL)
	}

	var results []dto.ResultDTO
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode results response from %s: %w", targetURL, err)
	}

	return results, nil
}

// GetPlayerStats fetches player stats from the ML service.
func (s *mlService) GetPlayerStats(playerID string) (*dto.PlayerStatsDTO, error) {
	// Ensure playerID is URL-safe, although less critical for path segments than query params
	safePlayerID := url.PathEscape(playerID)
	endpoint := fmt.Sprintf("/players/%s/stats", safePlayerID)
	targetURL := s.baseURL + endpoint

	req, err := http.NewRequest("GET", targetURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request to %s: %w", targetURL, err)
	}

	resp, err := s.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request to %s: %w", targetURL, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, nil // Return nil, nil to indicate not found, controller can handle 404
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ml service returned non-OK status (%d) for %s", resp.StatusCode, targetURL)
	}

	var stats dto.PlayerStatsDTO
	if err := json.NewDecoder(resp.Body).Decode(&stats); err != nil {
		return nil, fmt.Errorf("failed to decode player stats response from %s: %w", targetURL, err)
	}

	return &stats, nil
}