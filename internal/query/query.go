package query

import (
	"context"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/ranking"
)

type Service struct {
	catalog *catalog.Service
	ranking *ranking.Service
}

func New(c *catalog.Service, r *ranking.Service) *Service { return &Service{catalog: c, ranking: r} }
func (s *Service) Browse(ctx context.Context, f model.Filter) ([]model.Game, error) {
	return s.catalog.List(ctx, f)
}
func (s *Service) Top(ctx context.Context, game string, limit int) ([]model.RankingEntry, error) {
	return s.ranking.Leaderboard(ctx, game, limit)
}
func (s *Service) Detail(ctx context.Context, id string) (model.Game, error) {
	if err := ctx.Err(); err != nil {
		return model.Game{}, err
	}
	return s.catalog.Get(ctx, id)
}
