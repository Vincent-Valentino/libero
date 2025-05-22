package models

// TeamResponse represents a football team in the API responses
type TeamResponse struct {
	ID        int    `json:"id"`
	Name      string `json:"name"`
	ShortName string `json:"shortName"`
	Crest     string `json:"crest"`
}

// PlayerResponse represents a football player in the API responses
type PlayerResponse struct {
	ID          int    `json:"id"`
	Name        string `json:"name"`
	Position    string `json:"position"`
	Nationality string `json:"nationality"`
}

// PlayerStatsDTO represents the data structure for player statistics
type PlayerStatsDTO struct {
	PlayerID    string `json:"player_id"`
	PlayerName  string `json:"player_name"`
	Season      string `json:"season"`
	Appearances int    `json:"appearances"`
	Goals       int    `json:"goals"`
	Assists     int    `json:"assists"`
}
