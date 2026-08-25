package query

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type Insight struct {
	Kind   string
	Label  string
	Value  int
	Detail string
}

func (s *Service) Insights(ctx context.Context) ([]Insight, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	games, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	scores, e := s.ranking.Summaries(ctx)
	if e != nil {
		return nil, e
	}
	out := []Insight{}
	for _, c := range []string{"math", "english", "science", "logic"} {
		n := 0
		for _, g := range games {
			if g.Category == c && g.Published {
				n++
			}
		}
		out = append(out, Insight{Kind: "category", Label: c, Value: n, Detail: "published games"})
	}
	for _, x := range scores {
		out = append(out, Insight{Kind: "game", Label: x.GameID, Value: x.Plays, Detail: "plays"})
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Kind == out[j].Kind {
			return out[i].Label < out[j].Label
		}
		return out[i].Kind < out[j].Kind
	})
	return out, nil
}
func (s *Service) AgeDistribution(ctx context.Context) (map[int]int, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	games, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	out := map[int]int{}
	for _, g := range games {
		for age := g.AgeMin; age <= g.AgeMax; age++ {
			out[age]++
		}
	}
	return out, nil
}
func (s *Service) DeviceAvailability(ctx context.Context) (map[string]int, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	games, e := s.catalog.List(ctx, model.Filter{})
	if e != nil {
		return nil, e
	}
	out := map[string]int{"web": 0, "tablet": 0, "mobile": 0}
	for _, g := range games {
		if g.WebEnabled {
			out["web"]++
		}
		if g.TabletEnabled {
			out["tablet"]++
		}
		if g.MobileEnabled {
			out["mobile"]++
		}
	}
	return out, nil
}
func (s *Service) Search(ctx context.Context, term string, age int) ([]model.Game, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return s.catalog.List(ctx, model.Filter{Search: term, Age: age, PublishedOnly: true})
}
func (s *Service) Recommended(ctx context.Context, age int, device, category string, limit int) ([]model.Game, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	all, e := s.catalog.List(ctx, model.Filter{Age: age, Device: device, Category: category, PublishedOnly: true})
	if e != nil {
		return nil, e
	}
	sort.SliceStable(all, func(i, j int) bool {
		if all[i].Version == all[j].Version {
			return all[i].Title < all[j].Title
		}
		return all[i].Version > all[j].Version
	})
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all, nil
}
func (s *Service) ValidateBrowse(ctx context.Context, f model.Filter) (bool, error) {
	if e := ctx.Err(); e != nil {
		return false, e
	}
	if f.Age < 0 || f.Age > 18 {
		return false, nil
	}
	if f.Device != "" && f.Device != "web" && f.Device != "tablet" && f.Device != "mobile" {
		return false, nil
	}
	return true, nil
}
func (s *Service) CountPublished(ctx context.Context) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	g, e := s.catalog.List(ctx, model.Filter{PublishedOnly: true})
	return len(g), e
}
func (s *Service) CountPlayable(ctx context.Context, d string) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	g, e := s.catalog.List(ctx, model.Filter{Device: d, PublishedOnly: true})
	return len(g), e
}
