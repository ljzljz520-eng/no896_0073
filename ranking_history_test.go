package main_test

import (
	"context"
	"example.com/online-game-rank/internal/ranking"
	"example.com/online-game-rank/internal/store"
	"testing"
)

func TestHistoryEmpty(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/e.db")
	defer s.Close()
	h, e := ranking.New(s).PlayerHistory(context.Background(), "none")
	if e != nil || len(h) != 0 {
		t.Fatalf("history %v %d", e, len(h))
	}
}
