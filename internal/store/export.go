package store

import (
	"encoding/json"
	"errors"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type Snapshot struct {
	Games   []model.Game         `json:"games"`
	Players []model.Player       `json:"players"`
	Scores  []model.Score        `json:"scores"`
	Imports []model.ImportRecord `json:"imports"`
}

func (s *Store) Snapshot() (Snapshot, error) {
	games, e := s.Games()
	if e != nil {
		return Snapshot{}, e
	}
	players, e := s.Players()
	if e != nil {
		return Snapshot{}, e
	}
	scores, e := s.Scores()
	if e != nil {
		return Snapshot{}, e
	}
	imports, e := s.Imports()
	if e != nil {
		return Snapshot{}, e
	}
	sort.Slice(games, func(i, j int) bool { return games[i].ID < games[j].ID })
	sort.Slice(players, func(i, j int) bool { return players[i].ID < players[j].ID })
	sort.Slice(scores, func(i, j int) bool { return scores[i].ID < scores[j].ID })
	return Snapshot{games, players, scores, imports}, nil
}
func (s *Store) ExportJSON() ([]byte, error) {
	snap, e := s.Snapshot()
	if e != nil {
		return nil, e
	}
	return json.MarshalIndent(snap, "", "  ")
}
func (s *Store) ImportSnapshot(snap Snapshot) error {
	if len(snap.Games) == 0 && len(snap.Players) == 0 && len(snap.Scores) == 0 {
		return errors.New("empty snapshot")
	}
	return s.BeginWrite(func(t *Transaction) error {
		for _, g := range snap.Games {
			if e := t.SaveGame(g); e != nil {
				return e
			}
		}
		for _, r := range snap.Imports {
			if e := t.SaveImport(r); e != nil {
				return e
			}
		}
		return nil
	})
}
