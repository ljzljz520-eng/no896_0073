package catalog

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
	"strings"
)

type CategoryReport struct {
	Category  string
	Total     int
	Published int
	Playable  int
	MinAge    int
	MaxAge    int
}

func (s *Service) Report(ctx context.Context, device string) ([]CategoryReport, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	games, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	m := map[string]*CategoryReport{}
	for _, g := range games {
		r := m[g.Category]
		if r == nil {
			r = &CategoryReport{Category: g.Category, MinAge: g.AgeMin, MaxAge: g.AgeMax}
			m[g.Category] = r
		}
		r.Total++
		if g.Published {
			r.Published++
		}
		if g.IsPlayable(device) {
			r.Playable++
		}
		if g.AgeMin < r.MinAge {
			r.MinAge = g.AgeMin
		}
		if g.AgeMax > r.MaxAge {
			r.MaxAge = g.AgeMax
		}
	}
	out := make([]CategoryReport, 0, len(m))
	for _, r := range m {
		out = append(out, *r)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Category < out[j].Category })
	return out, nil
}
func (s *Service) SearchTitles(ctx context.Context, term string) ([]model.Game, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	term = strings.TrimSpace(term)
	return s.List(ctx, model.Filter{Search: term})
}
func (s *Service) PublishedByAge(ctx context.Context, age int) ([]model.Game, error) {
	if age < 3 || age > 18 {
		return []model.Game{}, nil
	}
	return s.List(ctx, model.Filter{Age: age, PublishedOnly: true})
}
func (s *Service) TogglePublication(ctx context.Context, ids []string, value bool) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	changed := 0
	for _, id := range ids {
		if e := s.Publish(ctx, id, value); e != nil {
			continue
		}
		changed++
	}
	return changed, nil
}
func (s *Service) Titles(games []model.Game) []string {
	out := make([]string, 0, len(games))
	for _, g := range games {
		out = append(out, g.Title)
	}
	sort.Strings(out)
	return out
}
