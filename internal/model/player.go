package model

import "time"

type Player struct {
	ID          string    `json:"id"`
	DisplayName string    `json:"display_name"`
	Age         int       `json:"age"`
	JoinedAt    time.Time `json:"joined_at"`
	Active      bool      `json:"active"`
}

func (p Player) Eligible() bool {
	return p.ID != "" && p.DisplayName != "" && p.Age >= 3 && p.Age <= 18 && p.Active
}

type Score struct {
	ID              string    `json:"id"`
	GameID          string    `json:"game_id"`
	PlayerID        string    `json:"player_id"`
	Points          int       `json:"points"`
	DurationSeconds int       `json:"duration_seconds"`
	PlayedAt        time.Time `json:"played_at"`
}

func (s Score) Valid() bool {
	return s.ID != "" && s.GameID != "" && s.PlayerID != "" && s.Points >= 0 && s.DurationSeconds > 0
}

type RankingEntry struct {
	Rank        int    `json:"rank"`
	PlayerID    string `json:"player_id"`
	PlayerName  string `json:"player_name"`
	Points      int    `json:"points"`
	GamesPlayed int    `json:"games_played"`
}
