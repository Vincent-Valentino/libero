package service

import (
	"encoding/json"
	"fmt"
	"io"

	"libero-backend/internal/models"
	"libero-backend/internal/repository"
	"net/http"
	"strings"
	"sync"
	"time"
)

// FixturesService defines the interface for fixture-related operations.
type FixturesService interface {
	GetTodaysFixtures() ([]models.CompetitionFixturesDTO, error)
	GetFixturesSummary(competitionCode string) (models.FixturesSummaryDTO, error)
}

// fixturesService implements the FixturesService interface.
type fixturesService struct {
	apiKey      string
	baseURL     string
	cacheRepo   repository.CacheRepository
	rateLimiter *rateLimiter
}

// rateLimiter provides basic rate limiting functionality
type rateLimiter struct {
	mutex           sync.Mutex
	lastRequestTime time.Time
	minInterval     time.Duration // Minimum time between API requests
}

// newRateLimiter creates a new rate limiter with specified interval
func newRateLimiter(interval time.Duration) *rateLimiter {
	return &rateLimiter{
		minInterval:     interval,
		lastRequestTime: time.Now().Add(-interval), // Allow immediate first request
	}
}

// Wait blocks until rate limit allows a new request
func (r *rateLimiter) Wait() {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	now := time.Now()
	// If we've made a request recently, wait until we can make another
	if !r.lastRequestTime.IsZero() {
		elapsed := now.Sub(r.lastRequestTime)
		if elapsed < r.minInterval {
			waitTime := r.minInterval - elapsed
			time.Sleep(waitTime)
			now = time.Now() // Update now after sleeping
		}
	}

	r.lastRequestTime = now
}

// NewFixturesService creates a new instance of fixturesService with API configuration.
func NewFixturesService(apiKey, baseURL string, cacheRepo repository.CacheRepository) FixturesService {
	// Create rate limiter with 1 request per 1.1 seconds (to be safe with API's limit)
	limiter := newRateLimiter(1100 * time.Millisecond)

	return &fixturesService{
		apiKey:      apiKey,
		baseURL:     baseURL,
		cacheRepo:   cacheRepo,
		rateLimiter: limiter,
	}
}

// GetTodaysFixtures fetches and filters fixtures for the current date and specified leagues,
// using cache when available.
func (s *fixturesService) GetTodaysFixtures() ([]models.CompetitionFixturesDTO, error) {
	// Try to get data from cache first
	cachedData, err := s.cacheRepo.GetCachedTodayFixtures()
	if err == nil && cachedData != nil {
		// Cache hit - convert JSONB to our DTO
		var fixtures []models.CompetitionFixturesDTO
		dataBytes, err := json.Marshal(cachedData.Data)
		if err == nil {
			if err := json.Unmarshal(dataBytes, &fixtures); err == nil {
				return fixtures, nil
			}
		}
	}

	// Cache miss or error - fetch from API
	today := time.Now().UTC().Format("2006-01-02")
	// Filter by relevant competition codes (PL, PD, SA, BL1, FL1, CL, EL)
	comps := []string{"PL", "PD", "SA", "BL1", "FL1", "CL", "EL"}
	compParam := strings.Join(comps, ",")

	// Wait for rate limiter before making request
	s.rateLimiter.Wait()
	// Build endpoint URL using date parameter
	url := fmt.Sprintf("%s/matches?dateFrom=%s&dateTo=%s&competitions=%s",
		s.baseURL,
		today,
		today,
		compParam,
	)

	// Debug URL (remove in production)
	fmt.Printf("Fetching matches from: %s\n", url)

	// Create HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set required headers
	req.Header.Set("X-Auth-Token", s.apiKey)
	req.Header.Set("Accept", "application/json")

	// Perform request with proper timeout
	client := http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusTooManyRequests {
		// Rate limited - try to get from cache even if it's expired
		cachedData, err = s.cacheRepo.GetCachedTodayFixturesIgnoringExpiry()
		if err == nil && cachedData != nil {
			var fixtures []models.CompetitionFixturesDTO
			dataBytes, err := json.Marshal(cachedData.Data)
			if err == nil {
				if err := json.Unmarshal(dataBytes, &fixtures); err == nil {
					return fixtures, nil
				}
			}
		}

		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("rate limited: %s", string(bodyBytes))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}

	// Decode provider response (assumes JSON of shape {matches: [...]})
	var raw struct {
		Matches []json.RawMessage `json:"matches"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	// Use a map from competition code to DTOs
	grouped := make(map[string][]models.FixtureMatchDTO)
	// Track competition metadata: code and emblem
	compMeta := make(map[string]struct{ Name, Code, Emblem string })
	for _, rawMatch := range raw.Matches {
		// Unmarshal into a generic map to extract required fields
		var m map[string]interface{}
		if err := json.Unmarshal(rawMatch, &m); err != nil {
			continue
		}
		// Extract competition
		comp, _ := m["competition"].(map[string]interface{})
		compName, _ := comp["name"].(string)
		compCode, _ := comp["code"].(string)
		emblem, _ := comp["emblem"].(string)
		// Extract teams
		home, _ := m["homeTeam"].(map[string]interface{})
		away, _ := m["awayTeam"].(map[string]interface{})
		// Record competition metadata for later
		if _, ok := compMeta[compCode]; !ok {
			compMeta[compCode] = struct{ Name, Code, Emblem string }{Name: compName, Code: compCode, Emblem: emblem}
		}
		// Extract score
		scoreObj, _ := m["score"].(map[string]interface{})
		ft, _ := scoreObj["fullTime"].(map[string]interface{})
		var homeScore *int
		var awayScore *int
		if v, ok := ft["homeTeam"].(float64); ok {
			iv := int(v)
			homeScore = &iv
		}
		if v, ok := ft["awayTeam"].(float64); ok {
			iv := int(v)
			awayScore = &iv
		}
		// Extract date and status
		dateStr, _ := m["utcDate"].(string)
		matchTime, _ := time.Parse(time.RFC3339, dateStr)
		status, _ := m["status"].(string)
		// Build DTO with additional fields
		venue, _ := m["venue"].(string)
		homeCrest, _ := home["crest"].(string)
		awayCrest, _ := away["crest"].(string)
		matchDTO := models.FixtureMatchDTO{
			MatchDate:    matchTime, // UTC timestamp
			HomeTeamName: home["name"].(string),
			AwayTeamName: away["name"].(string),
			HomeScore:    homeScore,
			AwayScore:    awayScore,
			MatchStatus:  status,
			Venue:        venue,
			HomeLogoURL:  homeCrest,
			AwayLogoURL:  awayCrest,
		}
		grouped[compCode] = append(grouped[compCode], matchDTO)
	}
	// Build final DTO array
	result := make([]models.CompetitionFixturesDTO, 0, len(grouped))
	for code, matches := range grouped {
		meta := compMeta[code]
		if len(matches) == 0 {
			continue
		}
		result = append(result, models.CompetitionFixturesDTO{
			CompetitionName: meta.Name,
			CompetitionCode: meta.Code,
			LogoURL:         meta.Emblem,
			Matches:         matches,
		})
	}
	// After building final DTO array, normalize PD name
	for i, cf := range result {
		// Map "PD" code to La Liga display name
		if cf.CompetitionCode == "PD" {
			result[i].CompetitionName = "La Liga"
		}
	}

	// Store in cache for future use (with 2 hour TTL)
	if len(result) > 0 {
		resultJSON, err := json.Marshal(result)
		if err == nil {
			var jsonData models.JSONB
			err = json.Unmarshal(resultJSON, &jsonData)
			if err == nil {
				s.cacheRepo.StoreCachedTodayFixtures(jsonData, 2*time.Hour)
			}
		}
	}

	return result, nil
}

// Implement the fixtures summary: today, tomorrow, upcoming
func (s *fixturesService) GetFixturesSummary(competitionCode string) (models.FixturesSummaryDTO, error) {
	// Ensure external API receives uppercase competition codes (e.g. 'PL')
	compCode := strings.ToUpper(competitionCode)

	// Try to get data from cache first
	cachedData, err := s.cacheRepo.GetCachedFixtures(compCode, "fixtures_summary")
	if err == nil && cachedData != nil {
		// Cache hit - convert JSONB to our DTO
		var summary models.FixturesSummaryDTO
		dataBytes, err := json.Marshal(cachedData.Data)
		if err == nil {
			if err := json.Unmarshal(dataBytes, &summary); err == nil {
				return summary, nil
			}
		}
	}

	// Cache miss or error - fetch from API
	base := strings.TrimRight(s.baseURL, "/")

	// Wait for rate limiter before making first request
	s.rateLimiter.Wait()

	// 1. Fetch competition metadata via /competitions/{code}
	compURL := fmt.Sprintf("%s/competitions/%s", base, compCode)
	reqComp, err := http.NewRequest(http.MethodGet, compURL, nil)
	if err != nil {
		return models.FixturesSummaryDTO{}, err
	}
	reqComp.Header.Set("X-Auth-Token", s.apiKey)
	// Perform request
	respComp, err := http.DefaultClient.Do(reqComp)
	if err != nil {
		return models.FixturesSummaryDTO{}, err
	}
	defer respComp.Body.Close()

	// Handle rate limit response
	if respComp.StatusCode == http.StatusTooManyRequests {
		// Try to get from cache even if it's expired
		expiredCache, expErr := s.cacheRepo.GetCachedFixturesIgnoringExpiry(compCode, "fixtures_summary")
		if expErr == nil && expiredCache != nil {
			var summary models.FixturesSummaryDTO
			dataBytes, err := json.Marshal(expiredCache.Data)
			if err == nil {
				if err := json.Unmarshal(dataBytes, &summary); err == nil {
					return summary, nil
				}
			}
		}

		bodyBytes, _ := io.ReadAll(respComp.Body)
		return models.FixturesSummaryDTO{}, fmt.Errorf("competition fetch failed: status %d, body: %s", respComp.StatusCode, string(bodyBytes))
	}

	if respComp.StatusCode != http.StatusOK {
		// Read response body for diagnostics
		bodyBytes, _ := io.ReadAll(respComp.Body)
		return models.FixturesSummaryDTO{}, fmt.Errorf("competition fetch failed: status %d, body: %s", respComp.StatusCode, string(bodyBytes))
	}

	var compRaw struct {
		Name   string `json:"name"`
		Emblem string `json:"emblem"`
	}
	if err := json.NewDecoder(respComp.Body).Decode(&compRaw); err != nil {
		return models.FixturesSummaryDTO{}, err
	}

	// Helper to fetch match list by date or dateFrom
	fetch := func(param string, useDateFrom bool) ([]models.FixtureMatchDTO, error) {
		// Wait for rate limit before each request
		s.rateLimiter.Wait()

		var url string
		if useDateFrom {
			url = fmt.Sprintf("%s/matches?dateFrom=%s&competitions=%s", base, param, compCode)
		} else {
			url = fmt.Sprintf("%s/matches?date=%s&competitions=%s", base, param, compCode)
		}
		req, err := http.NewRequest(http.MethodGet, url, nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Auth-Token", s.apiKey)
		// Perform request
		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		// Handle rate limit response
		if resp.StatusCode == http.StatusTooManyRequests {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("match fetch failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		}

		if resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("match fetch failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		}

		var raw struct {
			Matches []json.RawMessage `json:"matches"`
		}
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return nil, err
		}
		var out []models.FixtureMatchDTO
		for _, rm := range raw.Matches {
			var m map[string]interface{}
			if err := json.Unmarshal(rm, &m); err != nil {
				continue
			}
			// parse date, status, teams, scores, venue, crests
			dateStr, _ := m["utcDate"].(string)
			mt, _ := time.Parse(time.RFC3339, dateStr)
			venue, _ := m["venue"].(string)
			status, _ := m["status"].(string)
			scoreObj, _ := m["score"].(map[string]interface{})
			ft, _ := scoreObj["fullTime"].(map[string]interface{})
			var hs, as *int
			if v, ok := ft["homeTeam"].(float64); ok {
				i := int(v)
				hs = &i
			}
			if v, ok := ft["awayTeam"].(float64); ok {
				i := int(v)
				as = &i
			}
			home, _ := m["homeTeam"].(map[string]interface{})
			away, _ := m["awayTeam"].(map[string]interface{})
			hcrest, _ := home["crest"].(string)
			acrest, _ := away["crest"].(string)
			out = append(out, models.FixtureMatchDTO{
				MatchDate:    mt,
				HomeTeamName: home["name"].(string), AwayTeamName: away["name"].(string),
				HomeScore: hs, AwayScore: as,
				MatchStatus: status, Venue: venue,
				HomeLogoURL: hcrest, AwayLogoURL: acrest,
			})
		}
		return out, nil
	}

	// Dates for buckets
	now := time.Now().UTC()
	today := now.Format("2006-01-02")
	tomorrow := now.Add(24 * time.Hour).Format("2006-01-02")
	dayAfter := now.Add(48 * time.Hour).Format("2006-01-02")

	// Fetch each bucket with error handling
	var todayList, tomorrowList, upcomingList []models.FixtureMatchDTO

	todayList, err = fetch(today, false)
	if err != nil {
		// Log error but continue with other buckets
		fmt.Printf("Error fetching today's fixtures: %v\n", err)
	}

	tomorrowList, err = fetch(tomorrow, false)
	if err != nil {
		fmt.Printf("Error fetching tomorrow's fixtures: %v\n", err)
	}

	upcomingList, err = fetch(dayAfter, true)
	if err != nil {
		fmt.Printf("Error fetching upcoming fixtures: %v\n", err)
	}

	if len(upcomingList) > 4 {
		upcomingList = upcomingList[:4]
	}

	// Ensure slices are not nil to return empty arrays instead of null
	if todayList == nil {
		todayList = []models.FixtureMatchDTO{}
	}
	if tomorrowList == nil {
		tomorrowList = []models.FixtureMatchDTO{}
	}
	if upcomingList == nil {
		upcomingList = []models.FixtureMatchDTO{}
	}

	// Build summary DTO
	summary := models.FixturesSummaryDTO{
		CompetitionName: compRaw.Name,
		CompetitionCode: compCode,
		LogoURL:         compRaw.Emblem,
		Today:           todayList,
		Tomorrow:        tomorrowList,
		Upcoming:        upcomingList,
	}

	// Store in cache for future use (with 2 hour TTL)
	summaryJSON, err := json.Marshal(summary)
	if err == nil {
		var jsonData models.JSONB
		err = json.Unmarshal(summaryJSON, &jsonData)
		if err == nil {
			s.cacheRepo.StoreCachedFixtures(compCode, "fixtures_summary", jsonData, 2*time.Hour)
		}
	}

	return summary, nil
}
