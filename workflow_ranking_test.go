package main_test

import (
	"context"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/ranking"
	"example.com/online-game-rank/internal/store"
	"testing"
	"time"
)

func TestWorkflowTwo(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/b.db")
	defer s.Close()
	c := catalog.New(s)
	r := ranking.New(s)
	g := model.Game{ID: "logic-1", Title: "Logic", Category: "logic", AgeMin: 8, AgeMax: 14, EntryURL: "/logic", Published: true}
	_ = c.Create(context.Background(), g)
	_ = s.SavePlayer(model.Player{ID: "p1", DisplayName: "A", Age: 10, Active: true})
	_ = s.SavePlayer(model.Player{ID: "p2", DisplayName: "B", Age: 10, Active: true})
	_ = r.Record(context.Background(), model.Score{ID: "s1", GameID: g.ID, PlayerID: "p1", Points: 80, DurationSeconds: 30, PlayedAt: time.Unix(1, 0)})
	_ = r.Record(context.Background(), model.Score{ID: "s2", GameID: g.ID, PlayerID: "p2", Points: 90, DurationSeconds: 30, PlayedAt: time.Unix(2, 0)})
	top, e := r.Leaderboard(context.Background(), g.ID, 10)
	if e != nil || top[0].PlayerID != "p2" {
		t.Fatalf("ranking %#v %v", top, e)
	}
}
