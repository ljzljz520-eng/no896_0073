package admin

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/model"
)

type Role string

const (
	RoleViewer Role = "viewer"
	RoleEditor Role = "editor"
	RoleOwner  Role = "owner"
)

func Can(role Role, action string) bool {
	if role == RoleOwner {
		return true
	}
	if role == RoleEditor {
		return action != "delete_store" && action != "manage_users"
	}
	return action == "browse" || action == "report"
}
func Require(role Role, action string) error {
	if !Can(role, action) {
		return errors.New("permission denied")
	}
	return nil
}
func (s *Service) UpdateAs(ctx context.Context, role Role, g model.Game) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if e := Require(role, "edit_game"); e != nil {
		return e
	}
	return s.catalog.Update(ctx, g)
}
func (s *Service) PublishAs(ctx context.Context, role Role, id string, v bool) error {
	if e := ctx.Err(); e != nil {
		return e
	}
	if e := Require(role, "publish_game"); e != nil {
		return e
	}
	return s.catalog.Publish(ctx, id, v)
}
func (s *Service) ImportAs(ctx context.Context, role Role, source string, rows []model.Game) (model.ImportRecord, error) {
	if e := ctx.Err(); e != nil {
		return model.ImportRecord{}, e
	}
	if e := Require(role, "import_games"); e != nil {
		return model.ImportRecord{}, e
	}
	return s.ImportWithAudit(ctx, source, rows)
}
