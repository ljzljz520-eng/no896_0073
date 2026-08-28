package catalog

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/model"
	"fmt"
	"strings"
	"time"
)

type Importer struct{ catalog *Service }

func NewImporter(c *Service) *Importer { return &Importer{catalog: c} }
func (i *Importer) Simulate(ctx context.Context, source string, rows []model.Game) (model.ImportRecord, error) {
	if err := ctx.Err(); err != nil {
		return model.ImportRecord{}, err
	}
	if strings.TrimSpace(source) == "" {
		return model.ImportRecord{}, errors.New("source required")
	}
	r := model.ImportRecord{ID: fmt.Sprintf("import-%d", len(rows)), Source: source, Received: len(rows), StartedAt: time.Unix(0, 0)}
	for _, g := range rows {
		if err := ctx.Err(); err != nil {
			r.Status = "failed"
			r.Error = err.Error()
			_ = i.catalog.db.SaveImport(r)
			return r, err
		}
		if e := i.catalog.Create(ctx, g); e != nil {
			r.Rejected++
		} else {
			r.Accepted++
		}
	}
	r.Status = "completed"
	r.FinishedAt = time.Unix(int64(r.Accepted), 0)
	return r, i.catalog.db.SaveImport(r)
}
