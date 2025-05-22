package models

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
type ResultDTO struct {
	MatchDTO
	HomeScore      int      `json:"home_score"`
	AwayScore      int      `json:"away_score"`
	PossessionHome *float64 `json:"possession_home,omitempty"`
	PossessionAway *float64 `json:"possession_away,omitempty"`
}
