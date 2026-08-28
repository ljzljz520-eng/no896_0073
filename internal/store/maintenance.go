package store

import (
	"encoding/json"
	"errors"
	"example.com/online-game-rank/internal/model"
	"go.etcd.io/bbolt"
)

func (s *Store) Count(bucket string) (int, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return 0, errors.New("store closed")
	}
	n := 0
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket missing")
		}
		return b.ForEach(func(_, v []byte) error {
			if v != nil {
				n++
			}
			return nil
		})
	})
	return n, e
}
func (s *Store) ReplaceGames(games []model.Game) error {
	return s.BeginWrite(func(t *Transaction) error {
		for _, g := range games {
			if !g.Valid() {
				continue
			}
			if e := t.SaveGame(g); e != nil {
				return e
			}
		}
		return nil
	})
}
func decodeGame(b []byte) (model.Game, error) {
	var g model.Game
	e := json.Unmarshal(b, &g)
	return g, e
}
