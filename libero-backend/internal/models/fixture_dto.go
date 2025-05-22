package models

import "time"

// FixtureMatchDTO represents a single match within a competition's fixtures.
type FixtureMatchDTO struct {
	MatchDate    time.Time `json:"match_date"`
	HomeTeamName string    `json:"home_team_name"`
	AwayTeamName string    `json:"away_team_name"`
	HomeScore    *int      `json:"home_score,omitempty"`
	AwayScore    *int      `json:"away_score,omitempty"`
	MatchStatus  string    `json:"match_status"`
	Venue        string    `json:"venue,omitempty"`
	HomeLogoURL  string    `json:"home_logo_url,omitempty"`
	AwayLogoURL  string    `json:"away_logo_url,omitempty"`
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
