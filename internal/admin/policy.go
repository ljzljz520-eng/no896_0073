package admin

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/model"
	"strings"
)

type Policy struct {
	AllowedCategories map[string]bool
	MinimumAge        int
	MaximumAge        int
}

func DefaultPolicy() Policy {
	return Policy{AllowedCategories: map[string]bool{"math": true, "english": true, "science": true, "logic": true}, MinimumAge: 3, MaximumAge: 18}
}
func (p Policy) Check(g model.Game) error {
	if p.AllowedCategories == nil || !p.AllowedCategories[g.Category] {
		return errors.New("category not permitted")
	}
	if g.AgeMin < p.MinimumAge || g.AgeMax > p.MaximumAge {
		return errors.New("age outside policy")
	}
	return nil
}
func (s *Service) ValidateRows(ctx context.Context, rows []model.Game, p Policy) ([]model.Game, []string, error) {
	if e := ctx.Err(); e != nil {
		return nil, nil, e
	}
	accepted := []model.Game{}
	reasons := []string{}
	for _, g := range rows {
		if e := p.Check(g); e != nil {
			reasons = append(reasons, g.ID+": "+e.Error())
			continue
		}
		if e := s.catalog.Create(ctx, g); e != nil {
			reasons = append(reasons, g.ID+": "+e.Error())
			continue
		}
		accepted = append(accepted, g)
	}
	return accepted, reasons, nil
}
func (s *Service) Rename(ctx context.Context, id, title string) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if strings.TrimSpace(title) == "" {
		return errors.New("title required")
	}
	g, e := s.db.Game(id)
	if e != nil {
		return e
	}
	g.Title = strings.TrimSpace(title)
	g.Version++
	return s.db.SaveGame(g)
}
func (s *Service) DeactivatePlayer(ctx context.Context, id string) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	p, e := s.db.Player(id)
	if e != nil {
		return e
	}
	p.Active = false
	return s.db.SavePlayer(p)
}
