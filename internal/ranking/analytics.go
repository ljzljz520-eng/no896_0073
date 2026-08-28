package ranking

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"sort"
)

type PlayerSummary struct {
	PlayerID string
	Games    int
	Points   int
	Best     int
	Average  float64
}

func (s *Service) PlayerSummaries(ctx context.Context) ([]PlayerSummary, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	scores, e := s.db.Scores()
	if e != nil {
		return nil, e
	}
	m := map[string]*PlayerSummary{}
	for _, sc := range scores {
		v := m[sc.PlayerID]
		if v == nil {
			v = &PlayerSummary{PlayerID: sc.PlayerID}
			m[sc.PlayerID] = v
		}
		v.Games++
		v.Points += sc.Points
		if sc.Points > v.Best {
			v.Best = sc.Points
		}
	}
	out := make([]PlayerSummary, 0, len(m))
	for _, v := range m {
		v.Average = float64(v.Points) / float64(v.Games)
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Points == out[j].Points {
			return out[i].PlayerID < out[j].PlayerID
		}
		return out[i].Points > out[j].Points
	})
	return out, nil
}
func (s *Service) TopForPlayer(ctx context.Context, id string, n int) ([]model.Score, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	h, e := s.PlayerHistory(ctx, id)
	if e != nil {
		return nil, e
	}
	sort.Slice(h, func(i, j int) bool { return h[i].Points > h[j].Points })
	if n > 0 && len(h) > n {
		h = h[:n]
	}
	return h, nil
}
func (s *Service) PointsByGame(ctx context.Context, id string) (int, int, error) {
	if e := ctx.Err(); e != nil {
		return 0, 0, e
	}
	scores, e := s.db.Scores()
	if e != nil {
		return 0, 0, e
	}
	points, plays := 0, 0
	for _, sc := range scores {
		if sc.GameID == id {
			points += sc.Points
			plays++
		}
	}
	return points, plays, nil
}
