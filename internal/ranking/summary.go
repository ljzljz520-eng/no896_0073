package ranking

import (
	"context"
	"sort"
)

type GameSummary struct {
	GameID        string
	Plays         int
	TotalPoints   int
	AveragePoints float64
}

func (s *Service) Summaries(ctx context.Context) ([]GameSummary, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	scores, e := s.db.Scores()
	if e != nil {
		return nil, e
	}
	m := map[string]*GameSummary{}
	for _, sc := range scores {
		v := m[sc.GameID]
		if v == nil {
			v = &GameSummary{GameID: sc.GameID}
			m[sc.GameID] = v
		}
		v.Plays++
		v.TotalPoints += sc.Points
	}
	out := make([]GameSummary, 0, len(m))
	for _, v := range m {
		v.AveragePoints = float64(v.TotalPoints) / float64(v.Plays)
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].TotalPoints > out[j].TotalPoints })
	return out, nil
}
