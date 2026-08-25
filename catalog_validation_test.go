package main_test

import (
	"context"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"testing"
)

func TestCatalogRejectsInvalid(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/d.db")
	defer s.Close()
	if e := catalog.New(s).Create(context.Background(), model.Game{ID: "x", Title: "bad"}); e == nil {
		t.Fatal("invalid accepted")
	}
}
