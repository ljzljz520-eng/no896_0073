package admin

import (
	"context"
	"errors"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"os"
	"path/filepath"
	"testing"
)

// validGame builds a Game that passes model.Game.Valid() and the default policy.
func validGame(id, title string) model.Game {
	return model.Game{
		ID:         id,
		Title:      title,
		Category:   "math",
		AgeMin:     3,
		AgeMax:     8,
		EntryURL:   "https://example.com/" + id,
		Published:  true,
		WebEnabled: true,
	}
}

func newStore(t *testing.T) (*store.Store, func()) {
	t.Helper()
	dir := t.TempDir()
	db, err := store.Open(filepath.Join(dir, "test.db"))
	if err != nil {
		t.Fatalf("open store: %v", err)
	}
	cleanup := func() {
		_ = db.Close()
		_ = os.RemoveAll(dir)
	}
	return db, cleanup
}

// TestImportWithAuditCancelledStopsSideEffects verifies record isolation: a
// cancelled import must not persist records. The request's cancellation state
// must be inherited by every subsequent record write inside Simulate (and the
// per-row Create call), so the write never runs.
//
// This is the regression test for the bug where ImportWithAudit replaced the
// request context with context.Background(), so a cancelled import kept calling
// catalog.Create -> store.SaveGame (business side effects) after cancellation.
func TestImportWithAuditCancelledStopsSideEffects(t *testing.T) {
	db, cleanup := newStore(t)
	defer cleanup()

	cat := catalog.New(db)
	svc := New(db, cat)

	rows := []model.Game{validGame("g1", "Alpha"), validGame("g2", "Beta"), validGame("g3", "Gamma")}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // cancel up front: every subsequent operation must observe cancellation

	rec, err := svc.ImportWithAudit(ctx, "manual", rows)
	if err == nil {
		t.Fatalf("expected cancellation error, got nil (record=%+v)", rec)
	}
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", err)
	}

	// No record may have been persisted: a cancelled request must not produce
	// business side effects. Simulate checks ctx.Err() before the row loop and
	// Create() inherits the (cancelled) request context, so SaveGame never runs.
	games, gErr := db.Games()
	if gErr != nil {
		t.Fatalf("list games: %v", gErr)
	}
	if len(games) != 0 {
		t.Fatalf("expected zero persisted games after cancellation, got %d: %+v", len(games), games)
	}
}

// TestImportWithAuditCancelledDoesNotPersistImportsBucket confirms the cancelled
// request leaves no side-effect ImportRecord behind in the imports bucket either.
func TestImportWithAuditCancelledDoesNotPersistImportsBucket(t *testing.T) {
	db, cleanup := newStore(t)
	defer cleanup()

	cat := catalog.New(db)
	svc := New(db, cat)

	rows := []model.Game{validGame("g1", "Alpha"), validGame("g2", "Beta")}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	if _, err := svc.ImportWithAudit(ctx, "manual", rows); err == nil {
		t.Fatalf("expected cancellation error")
	}

	imports, iErr := db.Imports()
	if iErr != nil {
		t.Fatalf("list imports: %v", iErr)
	}
	if len(imports) != 0 {
		t.Fatalf("expected no import records after pre-cancellation, got %d: %+v", len(imports), imports)
	}
}

// TestImportWithAuditUncancelledPersistsAll is the positive control: without
// cancellation the same path persists every accepted row, confirming the fix
// did not break the happy path.
func TestImportWithAuditUncancelledPersistsAll(t *testing.T) {
	db, cleanup := newStore(t)
	defer cleanup()

	cat := catalog.New(db)
	svc := New(db, cat)

	rows := []model.Game{validGame("g1", "Alpha"), validGame("g2", "Beta")}
	rec, err := svc.ImportWithAudit(context.Background(), "manual", rows)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rec.Status != "completed" {
		t.Fatalf("expected status %q, got %q", "completed", rec.Status)
	}
	if rec.Accepted != 2 {
		t.Fatalf("expected 2 accepted, got %d", rec.Accepted)
	}
	games, _ := db.Games()
	if len(games) != 2 {
		t.Fatalf("expected 2 persisted games, got %d", len(games))
	}
}
