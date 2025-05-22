package models

// StandingsResponse represents the league standings data
type StandingsResponse struct {
	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"competition"`
	Season struct {
		ID      int  `json:"id"`
		Current bool `json:"current"`
	} `json:"season"`
	Standings []struct {
		Stage string `json:"stage"`
		Table []struct {
			Position       int          `json:"position"`
			Team           TeamResponse `json:"team"`
			PlayedGames    int          `json:"playedGames"`
			Won            int          `json:"won"`
			Draw           int          `json:"draw"`
			Lost           int          `json:"lost"`
			Points         int          `json:"points"`
			GoalsFor       int          `json:"goalsFor"`
			GoalsAgainst   int          `json:"goalsAgainst"`
			GoalDifference int          `json:"goalDifference"`
		} `json:"table"`
	} `json:"standings"`
}

// ScorersResponse represents the top scorers data
type ScorersResponse struct {
	Competition struct {
		ID   int    `json:"id"`
		Name string `json:"name"`
	} `json:"competition"`
	Season struct {
		ID      int  `json:"id"`
		Current bool `json:"current"`
	} `json:"season"`
	Scorers []struct {
		Player    PlayerResponse `json:"player"`
		Team      TeamResponse   `json:"team"`
		Goals     int            `json:"goals"`
		Assists   int            `json:"assists"`
		Penalties int            `json:"penalties"`
	} `json:"scorers"`
}

// CompetitionStandingsDTO represents formatted standings data for a competition
type CompetitionStandingsDTO struct {
	CompetitionName string              `json:"competition_name"`
	CompetitionCode string              `json:"competition_code"`
	Season          int                 `json:"season"`
	Standings       []StandingsTableDTO `json:"standings"`
}

// StandingsTableDTO represents a single standings table entry
type StandingsTableDTO struct {
	Position       int    `json:"position"`
	TeamName       string `json:"team_name"`
	TeamCrest      string `json:"team_crest"`
	PlayedGames    int    `json:"played"`
	Won            int    `json:"won"`
	Draw           int    `json:"drawn"`
	Lost           int    `json:"lost"`
	GoalsFor       int    `json:"goals_for"`
	GoalsAgainst   int    `json:"goals_against"`
	GoalDifference int    `json:"goal_difference"`
	Points         int    `json:"points"`
}

// CompetitionScorersDTO represents formatted top scorers data for a competition
type CompetitionScorersDTO struct {
	CompetitionName string           `json:"competition_name"`
	CompetitionCode string           `json:"competition_code"`
	Season          int              `json:"season"`
	Scorers         []ScorerStatsDTO `json:"scorers"`
}

// ScorerStatsDTO represents stats for a single scorer
type ScorerStatsDTO struct {
	PlayerName string `json:"player_name"`
	TeamName   string `json:"team_name"`
	TeamCrest  string `json:"team_crest"`
	Goals      int    `json:"goals"`
	Assists    int    `json:"assists"`
	Penalties  int    `json:"penalties"`
}
