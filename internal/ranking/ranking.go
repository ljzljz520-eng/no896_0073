package ranking

import (
	"context"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"sort"
)

type Service struct{ db *store.Store }

func New(db *store.Store) *Service { return &Service{db: db} }
func (s *Service) Record(ctx context.Context, sc model.Score) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !sc.Valid() {
		return ErrInvalidScore
	}
	if _, e := s.db.Game(sc.GameID); e != nil {
		return e
	}
	if _, e := s.db.Player(sc.PlayerID); e != nil {
		return e
	}
	return s.db.SaveScore(sc)
}
func (s *Service) Leaderboard(ctx context.Context, gameID string, limit int) ([]model.RankingEntry, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	scores, e := s.db.Scores()
	if e != nil {
		return nil, e
	}
	players, e := s.db.Players()
	if e != nil {
		return nil, e
	}
	names := map[string]string{}
	for _, p := range players {
		names[p.ID] = p.DisplayName
	}
	totals := map[string]*model.RankingEntry{}
	for _, sc := range scores {
		if gameID != "" && sc.GameID != gameID {
			continue
		}
		v := totals[sc.PlayerID]
		if v == nil {
			v = &model.RankingEntry{PlayerID: sc.PlayerID, PlayerName: names[sc.PlayerID]}
			totals[sc.PlayerID] = v
		}
		v.Points += sc.Points
		v.GamesPlayed++
	}
	out := make([]model.RankingEntry, 0, len(totals))
	for _, v := range totals {
		out = append(out, *v)
	}
	sort.Slice(out, func(i, j int) bool {
		if out[i].Points == out[j].Points {
			return out[i].PlayerName < out[j].PlayerName
		}
		return out[i].Points > out[j].Points
	})
	for n := range out {
		out[n].Rank = n + 1
	}
	if limit > 0 && len(out) > limit {
		out = out[:limit]
	}
	return out, nil
}
func (s *Service) PlayerHistory(ctx context.Context, id string) ([]model.Score, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	all, e := s.db.Scores()
	if e != nil {
		return nil, e
	}
	out := []model.Score{}
	for _, x := range all {
		if x.PlayerID == id {
			out = append(out, x)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].PlayedAt.Before(out[j].PlayedAt) })
	return out, nil
}
