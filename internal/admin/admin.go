package admin

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
)

type Service struct {
	db      *store.Store
	catalog *catalog.Service
}

func New(db *store.Store, c *catalog.Service) *Service { return &Service{db: db, catalog: c} }
func (s *Service) AddPlayer(ctx context.Context, p model.Player) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !p.Eligible() {
		return errors.New("invalid player")
	}
	if _, e := s.db.Player(p.ID); e == nil {
		return errors.New("player exists")
	}
	return s.db.SavePlayer(p)
}
func (s *Service) SetDeviceAvailability(ctx context.Context, id string, web, tablet, mobile bool) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g, e := s.db.Game(id)
	if e != nil {
		return e
	}
	g.WebEnabled, g.TabletEnabled, g.MobileEnabled = web, tablet, mobile
	return s.db.SaveGame(g)
}
func (s *Service) Archive(ctx context.Context, id string) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	g, e := s.db.Game(id)
	if e != nil {
		return e
	}
	g.Published = false
	return s.db.SaveGame(g)
}
func (s *Service) ImportWithAudit(ctx context.Context, source string, rows []model.Game) (model.ImportRecord, error) {
	// Propagate the request context so a cancelled import stops mutating records.
	// Simulate checks ctx.Err() per row and inside each Create call, so subsequent
	// record writes inherit the original request's cancellation state instead of
	// outliving it and continuing to produce side effects.
	return catalog.NewImporter(s.catalog).Simulate(ctx, source, rows)
}
