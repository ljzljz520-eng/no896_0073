package store

import (
	"encoding/json"
	"errors"
	"example.com/online-game-rank/internal/model"
	"go.etcd.io/bbolt"
)

type Transaction struct {
	tx     *bbolt.Tx
	closed bool
}

func (s *Store) BeginWrite(fn func(*Transaction) error) error {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.db == nil {
		return errors.New("store closed")
	}
	return s.db.Update(func(tx *bbolt.Tx) error { return fn(&Transaction{tx: tx}) })
}
func (t *Transaction) SaveGame(g model.Game) error {
	if t.closed {
		return errors.New("transaction closed")
	}
	b, e := json.Marshal(g)
	if e != nil {
		return e
	}
	return t.tx.Bucket([]byte("games")).Put([]byte(g.ID), b)
}
func (t *Transaction) SaveImport(r model.ImportRecord) error {
	if t.closed {
		return errors.New("transaction closed")
	}
	b, e := json.Marshal(r)
	if e != nil {
		return e
	}
	return t.tx.Bucket([]byte("imports")).Put([]byte(r.ID), b)
}
func (t *Transaction) Finish() { t.closed = true }
