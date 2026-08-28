package ranking

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

func (s *Service) BestScore(ctx context.Context, p, g string) (model.Score, bool, error) {
	if e := ctx.Err(); e != nil {
		return model.Score{}, false, e
	}
	a, e := s.db.Scores()
	if e != nil {
		return model.Score{}, false, e
	}
	var b model.Score
	f := false
	for _, x := range a {
		if x.PlayerID != p || (g != "" && x.GameID != g) {
			continue
		}
		if !f || x.Points > b.Points || (x.Points == b.Points && x.ID < b.ID) {
			b = x
			f = true
		}
	}
	return b, f, nil
}
func (s *Service) GameScores(ctx context.Context, g string) ([]model.Score, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	a, e := s.db.Scores()
	if e != nil {
		return nil, e
	}
	o := []model.Score{}
	for _, x := range a {
		if x.GameID == g {
			o = append(o, x)
		}
	}
	sort.Slice(o, func(i, j int) bool {
		if o[i].Points == o[j].Points {
			return o[i].PlayedAt.Before(o[j].PlayedAt)
		}
		return o[i].Points > o[j].Points
	})
	return o, nil
}
func (s *Service) ScorePercentile(ctx context.Context, g string, p int) (float64, error) {
	a, e := s.GameScores(ctx, g)
	if e != nil {
		return 0, e
	}
	if len(a) == 0 {
		return 0, nil
	}
	n := 0
	for _, x := range a {
		if x.Points <= p {
			n++
		}
	}
	return float64(n*100) / float64(len(a)), nil
}
func (s *Service) DistinctPlayers(ctx context.Context, g string) (int, error) {
	a, e := s.GameScores(ctx, g)
	if e != nil {
		return 0, e
	}
	m := map[string]bool{}
	for _, x := range a {
		m[x.PlayerID] = true
	}
	return len(m), nil
}
func (s *Service) AverageDuration(ctx context.Context, g string) (float64, error) {
	a, e := s.GameScores(ctx, g)
	if e != nil {
		return 0, e
	}
	if len(a) == 0 {
		return 0, nil
	}
	n := 0
	for _, x := range a {
		n += x.DurationSeconds
	}
	return float64(n) / float64(len(a)), nil
}
func (s *Service) NormalizePoints(ctx context.Context, g string) (map[string]float64, error) {
	a, e := s.GameScores(ctx, g)
	if e != nil {
		return nil, e
	}
	m := 0
	for _, x := range a {
		if x.Points > m {
			m = x.Points
		}
	}
	o := map[string]float64{}
	for _, x := range a {
		if m == 0 {
			o[x.ID] = 0
		} else {
			o[x.ID] = float64(x.Points) / float64(m)
		}
	}
	return o, nil
}
