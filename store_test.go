package main_test

import (
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"path/filepath"
	"testing"
)

func TestPersistenceSurvivesReopen(t *testing.T) {
	p := filepath.Join(t.TempDir(), "rank.db")
	s, e := store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	g := model.Game{ID: "g1", Title: "Fractions", Category: "math", AgeMin: 8, AgeMax: 12, EntryURL: "/g1", Published: true}
	if e = s.SaveGame(g); e != nil {
		t.Fatal(e)
	}
	s.Close()
	s, e = store.Open(p)
	if e != nil {
		t.Fatal(e)
	}
	defer s.Close()
	got, e := s.Game("g1")
	if e != nil || got.Title != g.Title {
		t.Fatalf("reopen: %#v %v", got, e)
	}
}
