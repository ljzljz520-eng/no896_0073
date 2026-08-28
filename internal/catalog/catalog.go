package catalog

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"sort"
	"strings"
)

type Service struct{ db *store.Store }

func New(db *store.Store) *Service { return &Service{db: db} }
func (s *Service) Create(ctx context.Context, g model.Game) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !g.Valid() {
		return errors.New("invalid game")
	}
	if _, e := s.db.Game(g.ID); e == nil {
		return errors.New("game exists")
	}
	return s.db.SaveGame(g)
}
func (s *Service) Update(ctx context.Context, g model.Game) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !g.Valid() {
		return errors.New("invalid game")
	}
	old, e := s.db.Game(g.ID)
	if e != nil {
		return errors.New("game not found")
	}
	g.Version = old.Version + 1
	return s.db.SaveGame(g)
}
func (s *Service) Publish(ctx context.Context, id string, value bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g, e := s.db.Game(id)
	if e != nil {
		return e
	}
	g.Published = value
	return s.db.SaveGame(g)
}
func (s *Service) List(ctx context.Context, f model.Filter) ([]model.Game, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	all, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	out := make([]model.Game, 0, len(all))
	for _, g := range all {
		if f.Matches(g) {
			out = append(out, g)
		}
	}
	sort.Slice(out, func(i, j int) bool { return strings.ToLower(out[i].Title) < strings.ToLower(out[j].Title) })
	return out, nil
}
func (s *Service) Remove(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	return s.db.DeleteGame(id)
}
func (s *Service) Get(ctx context.Context, id string) (model.Game, error) {
	if err := ctx.Err(); err != nil {
		return model.Game{}, err
	}
	return s.db.Game(id)
}
