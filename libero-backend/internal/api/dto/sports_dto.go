package dto

import "time"

// MatchDTO represents the data structure for a match.
type MatchDTO struct {
	ID          string    `json:"id"`
	Competition string    `json:"competition"`
	HomeTeam    string    `json:"home_team"`
	AwayTeam    string    `json:"away_team"`
	Date        time.Time `json:"date"`
	Status      string    `json:"status"`
}

// ResultDTO represents the data structure for a match result.
// It embeds MatchDTO and adds result-specific fields.
type ResultDTO struct {
	MatchDTO
	HomeScore      int      `json:"home_score"`
	AwayScore      int      `json:"away_score"`
	PossessionHome *float64 `json:"possession_home,omitempty"` // Use pointer for optional field
	PossessionAway *float64 `json:"possession_away,omitempty"` // Use pointer for optional field
	// Status is inherited from MatchDTO, libero-ml defaults it to "finished" for results
}

// PlayerStatsDTO represents the data structure for player statistics.
type PlayerStatsDTO struct {
	PlayerID    string `json:"player_id"`
	PlayerName  string `json:"player_name"`
	Season      string `json:"season"`
	Appearances int    `json:"appearances"`
	Goals       int    `json:"goals"`
	Assists     int    `json:"assists"`
}

// FixtureMatchDTO represents a single match within a competition's fixtures.
type FixtureMatchDTO struct {
	MatchDate     time.Time `json:"match_date"`
	HomeTeamName  string    `json:"home_team_name"`
	AwayTeamName  string    `json:"away_team_name"`
	HomeScore     *int      `json:"home_score,omitempty"` // Use pointer for optional (e.g., not started)
	AwayScore     *int      `json:"away_score,omitempty"` // Use pointer for optional
	MatchStatus   string    `json:"match_status"`       // e.g., "Scheduled", "Live", "Finished"
	Venue         string    `json:"venue,omitempty"`
	HomeLogoURL   string    `json:"home_logo_url,omitempty"`
	AwayLogoURL   string    `json:"away_logo_url,omitempty"`
}

// CompetitionFixturesDTO groups fixtures by competition.
type CompetitionFixturesDTO struct {
	CompetitionName string            `json:"competition_name"`
	CompetitionCode string            `json:"competition_code"`
	LogoURL         string            `json:"logo_url,omitempty"`
	Matches         []FixtureMatchDTO `json:"matches"`
}

// FixturesSummaryDTO groups fixtures for a single competition into time buckets.
type FixturesSummaryDTO struct {
	CompetitionName string            `json:"competition_name"`
	CompetitionCode string            `json:"competition_code"`
	LogoURL         string            `json:"logo_url,omitempty"`
	Today           []FixtureMatchDTO `json:"today"`
	Tomorrow        []FixtureMatchDTO `json:"tomorrow"`
	Upcoming        []FixtureMatchDTO `json:"upcoming"`
}