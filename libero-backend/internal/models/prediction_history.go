package models

import (
	"time"
)

// PredictionHistory represents a user's match prediction stored in the database
type PredictionHistory struct {
	ID                 uint      `gorm:"primaryKey" json:"id"`
	UserID             uint      `gorm:"not null;index;column:user_id" json:"userId"`
	HomeTeam           string    `gorm:"not null;column:home_team" json:"homeTeam"`
	AwayTeam           string    `gorm:"not null;column:away_team" json:"awayTeam"`
	HomeLeague         string    `gorm:"not null;column:home_league" json:"homeLeague"`
	AwayLeague         string    `gorm:"not null;column:away_league" json:"awayLeague"`
	PredictedHomeScore int       `gorm:"not null;column:predicted_home_score" json:"predictedHomeScore"`
	PredictedAwayScore int       `gorm:"not null;column:predicted_away_score" json:"predictedAwayScore"`
	ExpectedHomeGoals  float64   `gorm:"not null;column:expected_home_goals" json:"expectedHomeGoals"`
	ExpectedAwayGoals  float64   `gorm:"not null;column:expected_away_goals" json:"expectedAwayGoals"`
	HomeWinProbability float64   `gorm:"not null;column:home_win_probability" json:"homeWinProbability"`
	DrawProbability    float64   `gorm:"not null;column:draw_probability" json:"drawProbability"`
	AwayWinProbability float64   `gorm:"not null;column:away_win_probability" json:"awayWinProbability"`
	PredictedResult    string    `gorm:"not null;column:predicted_result" json:"predictedResult"`
	CreatedAt          time.Time `gorm:"column:created_at" json:"createdAt"`
	UpdatedAt          time.Time `gorm:"column:updated_at" json:"updatedAt"`

	// Relationship
	User User `gorm:"foreignKey:UserID" json:"user,omitempty"`
}

// --- DTOs for Prediction History API ---

// CreatePredictionRequest defines the request body for creating a new prediction
type CreatePredictionRequest struct {
	HomeTeam           string  `json:"homeTeam,omitempty" binding:"-"`
	AwayTeam           string  `json:"awayTeam,omitempty" binding:"-"`
	HomeLeague         string  `json:"homeLeague,omitempty" binding:"-"`
	AwayLeague         string  `json:"awayLeague,omitempty" binding:"-"`
	PredictedHomeScore int     `json:"predictedHomeScore,omitempty" binding:"-"`
	PredictedAwayScore int     `json:"predictedAwayScore,omitempty" binding:"-"`
	ExpectedHomeGoals  float64 `json:"expectedHomeGoals,omitempty" binding:"-"`
	ExpectedAwayGoals  float64 `json:"expectedAwayGoals,omitempty" binding:"-"`
	HomeWinProbability float64 `json:"homeWinProbability,omitempty" binding:"-"`
	DrawProbability    float64 `json:"drawProbability,omitempty" binding:"-"`
	AwayWinProbability float64 `json:"awayWinProbability,omitempty" binding:"-"`
	PredictedResult    string  `json:"predictedResult,omitempty" binding:"-"`

	// Support snake_case for backward compatibility
	HomeTeamSnake           string  `json:"home_team,omitempty" binding:"-"`
	AwayTeamSnake           string  `json:"away_team,omitempty" binding:"-"`
	HomeLeagueSnake         string  `json:"home_league,omitempty" binding:"-"`
	AwayLeagueSnake         string  `json:"away_league,omitempty" binding:"-"`
	PredictedHomeScoreSnake int     `json:"predicted_home_score,omitempty" binding:"-"`
	PredictedAwayScoreSnake int     `json:"predicted_away_score,omitempty" binding:"-"`
	ExpectedHomeGoalsSnake  float64 `json:"expected_home_goals,omitempty" binding:"-"`
	ExpectedAwayGoalsSnake  float64 `json:"expected_away_goals,omitempty" binding:"-"`
	HomeWinProbabilitySnake float64 `json:"home_win_probability,omitempty" binding:"-"`
	DrawProbabilitySnake    float64 `json:"draw_probability,omitempty" binding:"-"`
	AwayWinProbabilitySnake float64 `json:"away_win_probability,omitempty" binding:"-"`
	PredictedResultSnake    string  `json:"predicted_result,omitempty" binding:"-"`
}

// Normalize ensures that camelCase fields take precedence over snake_case
// This allows the frontend to send camelCase while maintaining backward compatibility
func (r *CreatePredictionRequest) Normalize() {
	// Use camelCase fields if available, otherwise fall back to snake_case
	if r.HomeTeam == "" && r.HomeTeamSnake != "" {
		r.HomeTeam = r.HomeTeamSnake
	}
	if r.AwayTeam == "" && r.AwayTeamSnake != "" {
		r.AwayTeam = r.AwayTeamSnake
	}
	if r.HomeLeague == "" && r.HomeLeagueSnake != "" {
		r.HomeLeague = r.HomeLeagueSnake
	}
	if r.AwayLeague == "" && r.AwayLeagueSnake != "" {
		r.AwayLeague = r.AwayLeagueSnake
	}
	if r.PredictedHomeScore == 0 && r.PredictedHomeScoreSnake != 0 {
		r.PredictedHomeScore = r.PredictedHomeScoreSnake
	}
	if r.PredictedAwayScore == 0 && r.PredictedAwayScoreSnake != 0 {
		r.PredictedAwayScore = r.PredictedAwayScoreSnake
	}
	if r.ExpectedHomeGoals == 0 && r.ExpectedHomeGoalsSnake != 0 {
		r.ExpectedHomeGoals = r.ExpectedHomeGoalsSnake
	}
	if r.ExpectedAwayGoals == 0 && r.ExpectedAwayGoalsSnake != 0 {
		r.ExpectedAwayGoals = r.ExpectedAwayGoalsSnake
	}
	if r.HomeWinProbability == 0 && r.HomeWinProbabilitySnake != 0 {
		r.HomeWinProbability = r.HomeWinProbabilitySnake
	}
	if r.DrawProbability == 0 && r.DrawProbabilitySnake != 0 {
		r.DrawProbability = r.DrawProbabilitySnake
	}
	if r.AwayWinProbability == 0 && r.AwayWinProbabilitySnake != 0 {
		r.AwayWinProbability = r.AwayWinProbabilitySnake
	}
	if r.PredictedResult == "" && r.PredictedResultSnake != "" {
		r.PredictedResult = r.PredictedResultSnake
	}
}

// PredictionHistoryResponse defines the response format for prediction history
type PredictionHistoryResponse struct {
	ID                 uint      `json:"id"`
	HomeTeam           string    `json:"homeTeam"`
	AwayTeam           string    `json:"awayTeam"`
	HomeLeague         string    `json:"homeLeague"`
	AwayLeague         string    `json:"awayLeague"`
	PredictedHomeScore int       `json:"predictedHomeScore"`
	PredictedAwayScore int       `json:"predictedAwayScore"`
	ExpectedHomeGoals  float64   `json:"expectedHomeGoals"`
	ExpectedAwayGoals  float64   `json:"expectedAwayGoals"`
	HomeWinProbability float64   `json:"homeWinProbability"`
	DrawProbability    float64   `json:"drawProbability"`
	AwayWinProbability float64   `json:"awayWinProbability"`
	PredictedResult    string    `json:"predictedResult"`
	CreatedAt          time.Time `json:"createdAt"`
}

// ToResponse converts PredictionHistory to PredictionHistoryResponse
func (p *PredictionHistory) ToResponse() PredictionHistoryResponse {
	return PredictionHistoryResponse{
		ID:                 p.ID,
		HomeTeam:           p.HomeTeam,
		AwayTeam:           p.AwayTeam,
		HomeLeague:         p.HomeLeague,
		AwayLeague:         p.AwayLeague,
		PredictedHomeScore: p.PredictedHomeScore,
		PredictedAwayScore: p.PredictedAwayScore,
		ExpectedHomeGoals:  p.ExpectedHomeGoals,
		ExpectedAwayGoals:  p.ExpectedAwayGoals,
		HomeWinProbability: p.HomeWinProbability,
		DrawProbability:    p.DrawProbability,
		AwayWinProbability: p.AwayWinProbability,
		PredictedResult:    p.PredictedResult,
		CreatedAt:          p.CreatedAt,
	}
}

// PredictionStatistics represents aggregated statistics for user's predictions
type PredictionStatistics struct {
	Total             int     `json:"total"`
	HomeWins          int     `json:"homeWins"`
	Draws             int     `json:"draws"`
	AwayWins          int     `json:"awayWins"`
	HomeWinPercentage float64 `json:"homeWinPercentage"`
	DrawPercentage    float64 `json:"drawPercentage"`
	AwayWinPercentage float64 `json:"awayWinPercentage"`
}
