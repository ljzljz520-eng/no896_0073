package query

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type Page struct {
	Items   []model.Game
	Offset  int
	Limit   int
	Total   int
	HasMore bool
}

func (s *Service) BrowsePage(ctx context.Context, f model.Filter, offset, limit int) (Page, error) {
	if e := ctx.Err(); e != nil {
		return Page{}, e
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 20
	}
	all, e := s.Browse(ctx, f)
	if e != nil {
		return Page{}, e
	}
	sort.Slice(all, func(i, j int) bool { return all[i].ID < all[j].ID })
	if offset > len(all) {
		offset = len(all)
	}
	end := offset + limit
	if end > len(all) {
		end = len(all)
	}
	return Page{Items: all[offset:end], Offset: offset, Limit: limit, Total: len(all), HasMore: end < len(all)}, nil
}
func (s *Service) Categories(ctx context.Context) ([]string, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	all, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	seen := map[string]bool{}
	out := []string{}
	for _, g := range all {
		if !seen[g.Category] {
			seen[g.Category] = true
			out = append(out, g.Category)
		}
	}
	sort.Strings(out)
	return out, nil
}
func (s *Service) Devices(ctx context.Context) ([]string, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return []string{"web", "tablet", "mobile"}, nil
}
