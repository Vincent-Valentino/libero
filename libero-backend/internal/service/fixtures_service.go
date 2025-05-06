package service

import (
	"encoding/json"
	"fmt"
	"io"
	"libero-backend/internal/api/dto"
	"net/http"
	"strings"
	"time"
)

// FixturesService defines the interface for fixture-related operations.
type FixturesService interface {
	GetTodaysFixtures() ([]dto.CompetitionFixturesDTO, error)
	GetFixturesSummary(competitionCode string) (dto.FixturesSummaryDTO, error)
}

// fixturesService implements the FixturesService interface.
type fixturesService struct {
	apiKey  string
	baseURL string
}

// NewFixturesService creates a new instance of fixturesService with API configuration.
func NewFixturesService(apiKey, baseURL string) FixturesService {
	return &fixturesService{apiKey: apiKey, baseURL: baseURL}
}

// GetTodaysFixtures fetches and filters mock fixtures for the current date and specified leagues.
func (s *fixturesService) GetTodaysFixtures() ([]dto.CompetitionFixturesDTO, error) {
	// Build today's date filter
	today := time.Now().UTC().Format("2006-01-02")
	// Filter by relevant competition codes (PL, PD, SA, BL1, FL1, CL, EL)
	comps := []string{"PL","PD","SA","BL1","FL1","CL","EL"}
	compParam := strings.Join(comps, ",")
	// Build endpoint URL using date parameter (v4): returns matches on that date
	url := fmt.Sprintf("%s/matches?date=%s&competitions=%s", strings.TrimRight(s.baseURL, "/"), today, compParam)
	// Create HTTP request
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}
	// Set API key header (adjust header name per provider)
	req.Header.Set("X-Auth-Token", s.apiKey)
	// Perform request
	client := http.Client{Timeout: 10 * time.Second}
	resp, err := client.Do(req)
		if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status %d", resp.StatusCode)
	}
	// Decode provider response (assumes JSON of shape {matches: [...]})
	var raw struct { Matches []json.RawMessage `json:"matches"` }
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		return nil, err
	}
	// Use a map from competition code to DTOs
	grouped := make(map[string][]dto.FixtureMatchDTO)
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
			iv := int(v); homeScore = &iv
		}
		if v, ok := ft["awayTeam"].(float64); ok {
			iv := int(v); awayScore = &iv
		}
		// Extract date and status
		dateStr, _ := m["utcDate"].(string)
		matchTime, _ := time.Parse(time.RFC3339, dateStr)
		status, _ := m["status"].(string)
		// Build DTO with additional fields
		venue, _ := m["venue"].(string)
		homeCrest, _ := home["crest"].(string)
		awayCrest, _ := away["crest"].(string)
			matchDTO := dto.FixtureMatchDTO{
			MatchDate:     matchTime,               // UTC timestamp
			HomeTeamName:  home["name"].(string),   
			AwayTeamName:  away["name"].(string),   
			HomeScore:     homeScore,
			AwayScore:     awayScore,
			MatchStatus:   status,
			Venue:         venue,
			HomeLogoURL:   homeCrest,
			AwayLogoURL:   awayCrest,
		}
		grouped[compCode] = append(grouped[compCode], matchDTO)
	}
	// Build final DTO array
	result := make([]dto.CompetitionFixturesDTO, 0, len(grouped))
	for code, matches := range grouped {
		meta := compMeta[code]
		if len(matches) == 0 {
			continue
		}
		result = append(result, dto.CompetitionFixturesDTO{
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
	return result, nil
}

// Implement the fixtures summary: today, tomorrow, upcoming
func (s *fixturesService) GetFixturesSummary(competitionCode string) (dto.FixturesSummaryDTO, error) {
	// Ensure external API receives uppercase competition codes (e.g. 'PL')
	compCode := strings.ToUpper(competitionCode)
	base := strings.TrimRight(s.baseURL, "/")
	// 1. Fetch competition metadata via /competitions/{code}
	compURL := fmt.Sprintf("%s/competitions/%s", base, compCode)
	reqComp, err := http.NewRequest(http.MethodGet, compURL, nil)
	if err != nil {
		return dto.FixturesSummaryDTO{}, err
	}
	reqComp.Header.Set("X-Auth-Token", s.apiKey)
	// Perform request
	respComp, err := http.DefaultClient.Do(reqComp)
	if err != nil || respComp.StatusCode != http.StatusOK {
		// Read response body for diagnostics
		bodyBytes, _ := io.ReadAll(respComp.Body)
		return dto.FixturesSummaryDTO{}, fmt.Errorf("competition fetch failed: status %d, body: %s", respComp.StatusCode, string(bodyBytes))
	}
	defer respComp.Body.Close()
	var compRaw struct { Name string `json:"name"`; Emblem string `json:"emblem"` }
	if err := json.NewDecoder(respComp.Body).Decode(&compRaw); err != nil {
		return dto.FixturesSummaryDTO{}, err
	}
	// Helper to fetch match list by date or dateFrom
	fetch := func(param string, useDateFrom bool) ([]dto.FixtureMatchDTO, error) {
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
		if err != nil || resp.StatusCode != http.StatusOK {
			bodyBytes, _ := io.ReadAll(resp.Body)
			return nil, fmt.Errorf("match fetch failed: status %d, body: %s", resp.StatusCode, string(bodyBytes))
		}
		defer resp.Body.Close()
		var raw struct { Matches []json.RawMessage `json:"matches"` }
		if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
			return nil, err
		}
		var out []dto.FixtureMatchDTO
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
			if v, ok := ft["homeTeam"].(float64); ok { i := int(v); hs = &i }
			if v, ok := ft["awayTeam"].(float64); ok { i := int(v); as = &i }
			home, _ := m["homeTeam"].(map[string]interface{})
			away, _ := m["awayTeam"].(map[string]interface{})
			hcrest, _ := home["crest"].(string)
			acrest, _ := away["crest"].(string)
			out = append(out, dto.FixtureMatchDTO{
				MatchDate:   mt,
				HomeTeamName: home["name"].(string), AwayTeamName: away["name"].(string),
				HomeScore:   hs, AwayScore: as,
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
	// Fetch each bucket
	todayList, _ := fetch(today, false)
	tomorrowList, _ := fetch(tomorrow, false)
	upcomingList, _ := fetch(dayAfter, true)
	if len(upcomingList) > 4 { upcomingList = upcomingList[:4] }
	// Ensure slices are not nil to return empty arrays instead of null
	if todayList == nil {
		todayList = []dto.FixtureMatchDTO{}
	}
	if tomorrowList == nil {
		tomorrowList = []dto.FixtureMatchDTO{}
	}
	if upcomingList == nil {
		upcomingList = []dto.FixtureMatchDTO{}
	}
	// Build summary DTO
	return dto.FixturesSummaryDTO{
		CompetitionName: compRaw.Name,
		CompetitionCode: compCode,
		LogoURL:         compRaw.Emblem,
		Today:           todayList,
		Tomorrow:        tomorrowList,
		Upcoming:        upcomingList,
	}, nil
}