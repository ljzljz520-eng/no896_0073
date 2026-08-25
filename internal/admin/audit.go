package admin

import (
	"context"
	"example.com/online-game-rank/internal/model"
)

type AuditEvent struct {
	Action   string
	EntityID string
	Outcome  string
}

func (s *Service) ReviewImport(ctx context.Context) ([]model.ImportRecord, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	return s.db.Imports()
}
func (s *Service) Events(ctx context.Context, id string) ([]AuditEvent, error) {
	if e := ctx.Err(); e != nil {
		return nil, e
	}
	records, e := s.db.Imports()
	if e != nil {
		return nil, e
	}
	out := []AuditEvent{}
	for _, r := range records {
		if id == "" || r.ID == id {
			out = append(out, AuditEvent{Action: "import", EntityID: r.ID, Outcome: r.Status})
		}
	}
	return out, nil
}
