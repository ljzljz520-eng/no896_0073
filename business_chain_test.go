package main_test

import (
	"context"
	"example.com/online-game-rank/internal/admin"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"testing"
)

func TestBusinessChain32(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/bug.db")
	defer s.Close()
	c := catalog.New(s)
	a := admin.New(s, c)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, e := a.ImportWithAudit(ctx, "cancelled", []model.Game{{ID: "side", Title: "Side Effect", Category: "science", AgeMin: 8, AgeMax: 12, EntryURL: "/side"}})
	if e == nil {
		t.Fatal("cancelled request should return an error")
	}
	games, _ := s.Games()
	if len(games) != 0 {
		t.Fatalf("cancelled request created %d games", len(games))
	}
}
