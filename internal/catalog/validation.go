package catalog

import (
	"errors"
	"example.com/online-game-rank/internal/model"
	"strings"
)

func ValidateBatch(rows []model.Game) (accepted, rejected int, reasons []string) {
	for _, g := range rows {
		if strings.TrimSpace(g.Title) == "" {
			rejected++
			reasons = append(reasons, g.ID+": title")
			continue
		}
		if !g.Valid() {
			rejected++
			reasons = append(reasons, g.ID+": fields")
			continue
		}
		accepted++
	}
	return
}
func EnsureUnique(rows []model.Game) error {
	seen := map[string]bool{}
	for _, g := range rows {
		if seen[g.ID] {
			return errors.New("duplicate game: " + g.ID)
		}
		seen[g.ID] = true
	}
	return nil
}
