package query

import (
	"context"
	"encoding/json"
	"example.com/online-game-rank/internal/model"
)

type CatalogEnvelope struct {
	Games      []model.Game `json:"games"`
	Count      int          `json:"count"`
	Categories []string     `json:"categories"`
}

func (s *Service) ExportCatalog(ctx context.Context) ([]byte, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	c, e := s.Categories(ctx)
	if e != nil {
		return nil, e
	}
	return json.Marshal(CatalogEnvelope{Games: g, Count: len(g), Categories: c})
}
func (s *Service) ExportPlayer(ctx context.Context, id string) ([]byte, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	h, e := s.ranking.PlayerHistory(ctx, id)
	if e != nil {
		return nil, e
	}
	return json.Marshal(struct {
		PlayerID string        `json:"player_id"`
		Scores   []model.Score `json:"scores"`
		Count    int           `json:"count"`
	}{id, h, len(h)})
}
func (s *Service) GameIDs(ctx context.Context, f model.Filter) ([]string, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.catalog.List(ctx, f)
	if e != nil {
		return nil, e
	}
	o := make([]string, 0, len(g))
	for _, x := range g {
		o = append(o, x.ID)
	}
	return o, nil
}
func (s *Service) Empty(ctx context.Context) bool {
	if ctx.Err() != nil {
		return true
	}
	return false
}
func (s *Service) HasCategory(ctx context.Context, category string) (bool, error) {
	if e := ctx.Err(); e != nil {
		return false, e
	}
	items, e := s.catalog.List(ctx, model.Filter{Category: category})
	return len(items) > 0, e
}
func (s *Service) HasDevice(ctx context.Context, device string) (bool, error) {
	if e := ctx.Err(); e != nil {
		return false, e
	}
	items, e := s.catalog.List(ctx, model.Filter{Device: device, PublishedOnly: true})
	return len(items) > 0, e
}
func (s *Service) PublishedIDs(ctx context.Context) ([]string, error) {
	return s.GameIDs(ctx, model.Filter{PublishedOnly: true})
}
func (s *Service) CountByCategory(ctx context.Context) (map[string]int, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	items, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	out := map[string]int{}
	for _, g := range items {
		out[g.Category]++
	}
	return out, nil
}
func (s *Service) CountByAge(ctx context.Context, age int) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	items, e := s.catalog.List(ctx, model.Filter{Age: age})
	return len(items), e
}
func (s *Service) CountBySearch(ctx context.Context, term string) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	items, e := s.catalog.SearchTitles(ctx, term)
	return len(items), e
}
func (s *Service) IsEmpty(ctx context.Context) (bool, error) {
	if e := ctx.Err(); e != nil {
		return false, e
	}
	items, e := s.catalog.List(ctx, model.Filter{})
	return len(items) == 0, e
}
func (s *Service) CategoryExists(ctx context.Context, c string) (bool, error) {
	if e := ctx.Err(); e != nil {
		return false, e
	}
	counts, e := s.CountByCategory(ctx)
	return counts[c] > 0, e
}
func (s *Service) DeviceCounts(ctx context.Context) (map[string]int, error) {
	return s.DeviceAvailability(ctx)
}
func (s *Service) RecommendedIDs(ctx context.Context, age int, device string) ([]string, error) {
	g, e := s.Recommended(ctx, age, device, "", 0)
	if e != nil {
		return nil, e
	}
	o := []string{}
	for _, x := range g {
		o = append(o, x.ID)
	}
	return o, nil
}
func (s *Service) PublishedCategory(ctx context.Context, c string) ([]model.Game, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return s.catalog.List(ctx, model.Filter{Category: c, PublishedOnly: true})
}
func (s *Service) PublishedAge(ctx context.Context, a int) ([]model.Game, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return s.catalog.List(ctx, model.Filter{Age: a, PublishedOnly: true})
}
func (s *Service) PublishedDevice(ctx context.Context, d string) ([]model.Game, error) {
	return s.catalog.List(ctx, model.Filter{Device: d, PublishedOnly: true})
}
