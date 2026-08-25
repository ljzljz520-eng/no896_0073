package store

import (
	"encoding/json"
	"errors"
	"example.com/online-game-rank/internal/model"
	"go.etcd.io/bbolt"
	"path/filepath"
	"sync"
	"time"
)

var buckets = [][]byte{[]byte("games"), []byte("players"), []byte("scores"), []byte("imports")}

type Store struct {
	db *bbolt.DB
	mu sync.RWMutex
}

func Open(path string) (*Store, error) {
	db, err := bbolt.Open(filepath.Clean(path), 0600, &bbolt.Options{Timeout: time.Second})
	if err != nil {
		return nil, err
	}
	s := &Store{db: db}
	err = db.Update(func(tx *bbolt.Tx) error {
		for _, b := range buckets {
			if _, e := tx.CreateBucketIfNotExists(b); e != nil {
				return e
			}
		}
		return nil
	})
	if err != nil {
		db.Close()
		return nil, err
	}
	return s, nil
}
func (s *Store) Close() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.db == nil {
		return nil
	}
	err := s.db.Close()
	s.db = nil
	return err
}
func (s *Store) put(bucket, key string, v any) error {
	data, e := json.Marshal(v)
	if e != nil {
		return e
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Put([]byte(key), data) })
}
func (s *Store) get(bucket, key string, v any) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		raw := tx.Bucket([]byte(bucket)).Get([]byte(key))
		if raw == nil {
			return bbolt.ErrBucketNotFound
		}
		return json.Unmarshal(raw, v)
	})
}
func (s *Store) all(bucket string, decode func([]byte) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.View(func(tx *bbolt.Tx) error {
		return tx.Bucket([]byte(bucket)).ForEach(func(_, v []byte) error {
			if v == nil {
				return nil
			}
			return decode(v)
		})
	})
}
func (s *Store) del(bucket, key string) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return tx.Bucket([]byte(bucket)).Delete([]byte(key)) })
}

func (s *Store) SaveGame(g model.Game) error { return s.put("games", g.ID, g) }
func (s *Store) Game(id string) (model.Game, error) {
	var g model.Game
	e := s.get("games", id, &g)
	return g, e
}
func (s *Store) DeleteGame(id string) error { return s.del("games", id) }
func (s *Store) Games() ([]model.Game, error) {
	out := []model.Game{}
	e := s.all("games", func(b []byte) error {
		var g model.Game
		if x := json.Unmarshal(b, &g); x != nil {
			return x
		}
		out = append(out, g)
		return nil
	})
	return out, e
}
func (s *Store) SavePlayer(p model.Player) error { return s.put("players", p.ID, p) }
func (s *Store) Player(id string) (model.Player, error) {
	var p model.Player
	e := s.get("players", id, &p)
	return p, e
}
func (s *Store) Players() ([]model.Player, error) {
	out := []model.Player{}
	e := s.all("players", func(b []byte) error {
		var p model.Player
		if x := json.Unmarshal(b, &p); x != nil {
			return x
		}
		out = append(out, p)
		return nil
	})
	return out, e
}
func (s *Store) SaveScore(sc model.Score) error { return s.put("scores", sc.ID, sc) }
func (s *Store) Scores() ([]model.Score, error) {
	out := []model.Score{}
	e := s.all("scores", func(b []byte) error {
		var x model.Score
		if z := json.Unmarshal(b, &x); z != nil {
			return z
		}
		out = append(out, x)
		return nil
	})
	return out, e
}
func (s *Store) SaveImport(r model.ImportRecord) error { return s.put("imports", r.ID, r) }
func (s *Store) Imports() ([]model.ImportRecord, error) {
	out := []model.ImportRecord{}
	e := s.all("imports", func(b []byte) error {
		var x model.ImportRecord
		if z := json.Unmarshal(b, &x); z != nil {
			return z
		}
		out = append(out, x)
		return nil
	})
	return out, e
}
