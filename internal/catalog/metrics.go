package catalog

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type AgeBand struct {
	Label string
	Min   int
	Max   int
	Count int
}

func (s *Service) AgeBands(ctx context.Context, b []AgeBand) ([]AgeBand, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	o := append([]AgeBand(nil), b...)
	for i := range o {
		o[i].Count = 0
		for _, x := range g {
			if x.AgeMin <= o[i].Max && x.AgeMax >= o[i].Min {
				o[i].Count++
			}
		}
	}
	sort.Slice(o, func(i, j int) bool { return o[i].Min < o[j].Min })
	return o, nil
}
func (s *Service) CountByDevice(ctx context.Context) (map[string]int, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	o := map[string]int{"web": 0, "tablet": 0, "mobile": 0}
	for _, x := range g {
		for _, d := range []string{"web", "tablet", "mobile"} {
			if x.IsPlayable(d) {
				o[d]++
			}
		}
	}
	return o, nil
}
func (s *Service) CategoryTitles(ctx context.Context, c string) ([]string, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.List(ctx, model.Filter{Category: c})
	if e != nil {
		return nil, e
	}
	o := []string{}
	for _, x := range g {
		o = append(o, x.Title)
	}
	sort.Strings(o)
	return o, nil
}
func (s *Service) Versions(ctx context.Context) (map[string]int, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	o := map[string]int{}
	for _, x := range g {
		o[x.ID] = x.Version
	}
	return o, nil
}
func (s *Service) Eligible(ctx context.Context, age int, d string) ([]model.Game, error) {
	return s.List(ctx, model.Filter{Age: age, Device: d, PublishedOnly: true})
}
func (s *Service) MissingLinks(ctx context.Context) ([]string, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	g, e := s.db.Games()
	if e != nil {
		return nil, e
	}
	o := []string{}
	for _, x := range g {
		if x.EntryURL == "" {
			o = append(o, x.ID)
		}
	}
	return o, nil
}
func (s *Service) PublishAll(ctx context.Context, c string) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	g, e := s.db.Games()
	if e != nil {
		return 0, e
	}
	n := 0
	for _, x := range g {
		if c != "" && x.Category != c {
			continue
		}
		if !x.Published {
			x.Published = true
			if e = s.db.SaveGame(x); e == nil {
				n++
			}
		}
	}
	return n, nil
}
func (s *Service) ArchiveOlder(ctx context.Context, v int) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	g, e := s.db.Games()
	if e != nil {
		return 0, e
	}
	n := 0
	for _, x := range g {
		if x.Version <= v && x.Published {
			x.Published = false
			if e = s.db.SaveGame(x); e == nil {
				n++
			}
		}
	}
	return n, nil
}
