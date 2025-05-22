package repository

import (
	"encoding/json"
	"fmt"
	"io"
	"libero-backend/config"
	"libero-backend/internal/models"
	"net/http"
	"time"
)

type FootballAPIRepository interface {
	GetStandings(competitionCode string) (models.StandingsResponse, error)
	GetTopScorers(competitionCode string) (models.ScorersResponse, error)
}

type footballAPIRepository struct {
	apiKey    string
	baseURL   string
	client    *http.Client
	cacheRepo CacheRepository
}

func NewFootballAPIRepository(cfg *config.Config, cacheRepo CacheRepository) FootballAPIRepository {
	return &footballAPIRepository{
		apiKey:    cfg.ThirdPartyAPIKey,
		baseURL:   cfg.ThirdPartyBaseURL,
		client:    &http.Client{},
		cacheRepo: cacheRepo,
	}
}

func (r *footballAPIRepository) makeRequest(endpoint string) (*http.Response, error) {
	req, err := http.NewRequest("GET", endpoint, nil)
	if err != nil {
		return nil, fmt.Errorf("error creating request: %w", err)
	}

	req.Header.Set("X-Auth-Token", r.apiKey)

	// Try the request
	resp, err := r.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("error making request: %w", err)
	}

	// Handle rate limiting with exponential backoff
	if resp.StatusCode == http.StatusTooManyRequests {
		retryAfter := resp.Header.Get("Retry-After")
		if retryAfter != "" {
			waitSeconds, _ := time.ParseDuration(retryAfter + "s")
			time.Sleep(waitSeconds)
			return r.client.Do(req)
		}
		// If no Retry-After header, use default backoff
		time.Sleep(10 * time.Second)
		return r.client.Do(req)
	}

	return resp, nil
}

func (r *footballAPIRepository) GetStandings(competitionCode string) (models.StandingsResponse, error) {
	var response models.StandingsResponse

	// Try to get from cache first
	cachedData, err := r.cacheRepo.GetCachedFixtures(competitionCode, "standings")
	if err == nil && cachedData != nil {
		// Convert cached JSONB data back to response struct
		cachedBytes, err := json.Marshal(cachedData.Data)
		if err == nil {
			if err := json.Unmarshal(cachedBytes, &response); err == nil {
				return response, nil
			}
		}
	}

	// Cache miss or error - fetch from API
	endpoint := fmt.Sprintf("%s/competitions/%s/standings", r.baseURL, competitionCode)
	resp, err := r.makeRequest(endpoint)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("error decoding response: %w", err)
	}

	// Cache the response
	// First convert response to raw JSON
	responseBytes, err := json.Marshal(response)
	if err == nil {
		var jsonData models.JSONB
		if err := json.Unmarshal(responseBytes, &jsonData); err == nil {
			r.cacheRepo.StoreCachedFixtures(competitionCode, "standings", jsonData, 1*time.Hour)
		}
	}

	return response, nil
}

func (r *footballAPIRepository) GetTopScorers(competitionCode string) (models.ScorersResponse, error) {
	var response models.ScorersResponse

	// Try to get from cache first
	cachedData, err := r.cacheRepo.GetCachedFixtures(competitionCode, "scorers")
	if err == nil && cachedData != nil {
		// Convert cached JSONB data back to response struct
		cachedBytes, err := json.Marshal(cachedData.Data)
		if err == nil {
			if err := json.Unmarshal(cachedBytes, &response); err == nil {
				return response, nil
			}
		}
	}

	// Cache miss or error - fetch from API
	endpoint := fmt.Sprintf("%s/competitions/%s/scorers", r.baseURL, competitionCode)
	resp, err := r.makeRequest(endpoint)
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return response, fmt.Errorf("API returned non-200 status: %d, body: %s", resp.StatusCode, string(body))
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return response, fmt.Errorf("error decoding response: %w", err)
	}

	// Cache the response
	// First convert response to raw JSON
	responseBytes, err := json.Marshal(response)
	if err == nil {
		var jsonData models.JSONB
		if err := json.Unmarshal(responseBytes, &jsonData); err == nil {
			r.cacheRepo.StoreCachedFixtures(competitionCode, "scorers", jsonData, 1*time.Hour)
		}
	}

	return response, nil
}
