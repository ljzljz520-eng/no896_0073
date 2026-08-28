package catalog

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type ScheduleItem struct {
	GameID   string
	Priority int
	Reason   string
}

func (s *Service) Schedule(ctx context.Context, age int, device string) ([]ScheduleItem, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	games, e := s.List(ctx, model.Filter{Age: age, Device: device, PublishedOnly: true})
	if e != nil {
		return nil, e
	}
	o := make([]ScheduleItem, 0, len(games))
	for _, g := range games {
		p := g.Version
		if g.Category == "math" {
			p += 3
		}
		if g.Category == "english" {
			p += 2
		}
		if g.AgeMin == age {
			p++
		}
		o = append(o, ScheduleItem{GameID: g.ID, Priority: p, Reason: model.CategoryLabel(g.Category)})
	}
	sort.Slice(o, func(i, j int) bool {
		if o[i].Priority == o[j].Priority {
			return o[i].GameID < o[j].GameID
		}
		return o[i].Priority > o[j].Priority
	})
	return o, nil
}
func (s *Service) Next(ctx context.Context, age int, device string) (model.Game, error) {
	a, e := s.Schedule(ctx, age, device)
	if e != nil {
		return model.Game{}, e
	}
	if len(a) == 0 {
		return model.Game{}, nil
	}
	return s.Get(ctx, a[0].GameID)
}
func (s *Service) RefreshLinks(ctx context.Context, updates map[string]string) (int, error) {
	if e := ctx.Err(); e != nil {
		return 0, e
	}
	n := 0
	for id, url := range updates {
		g, e := s.db.Game(id)
		if e != nil {
			continue
		}
		if url == "" {
			continue
		}
		g.EntryURL = url
		g.Version++
		if e = s.db.SaveGame(g); e == nil {
			n++
		}
	}
	return n, nil
}
