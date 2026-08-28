package store

import (
	"errors"
	"go.etcd.io/bbolt"
	"time"
)

type Health struct {
	Open      bool
	Buckets   map[string]int
	CheckedAt time.Time
}

func (s *Store) Health() (Health, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return Health{}, errors.New("store closed")
	}
	h := Health{Open: true, Buckets: map[string]int{}, CheckedAt: time.Unix(0, 0)}
	e := s.db.View(func(tx *bbolt.Tx) error {
		for _, b := range []string{"games", "players", "scores", "imports"} {
			n := 0
			bucket := tx.Bucket([]byte(b))
			if bucket == nil {
				continue
			}
			_ = bucket.ForEach(func(k, v []byte) error {
				if v != nil {
					n++
				}
				return nil
			})
			h.Buckets[b] = n
		}
		return nil
	})
	return h, e
}
func (s *Store) Has(bucket, key string) (bool, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return false, errors.New("store closed")
	}
	found := false
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket missing")
		}
		found = b.Get([]byte(key)) != nil
		return nil
	})
	return found, e
}
func (s *Store) Keys(bucket string) ([]string, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return nil, errors.New("store closed")
	}
	out := []string{}
	e := s.db.View(func(tx *bbolt.Tx) error {
		b := tx.Bucket([]byte(bucket))
		if b == nil {
			return errors.New("bucket missing")
		}
		return b.ForEach(func(k, v []byte) error {
			if v != nil {
				out = append(out, string(k))
			}
			return nil
		})
	})
	return out, e
}
