package main_test

import (
	"context"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"testing"
)

func TestWorkflowOne(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/a.db")
	defer s.Close()
	c := catalog.New(s)
	g := model.Game{ID: "math-1", Title: "Number Quest", Category: "math", AgeMin: 7, AgeMax: 10, WebEnabled: true, EntryURL: "https://x", Published: true}
	if e := c.Create(context.Background(), g); e != nil {
		t.Fatal(e)
	}
	items, e := c.List(context.Background(), model.Filter{Category: "math", Age: 8, Device: "web", PublishedOnly: true})
	if e != nil || len(items) != 1 {
		t.Fatalf("list %v %d", e, len(items))
	}
	if e = c.Publish(context.Background(), g.ID, false); e != nil {
		t.Fatal(e)
	}
	items, _ = c.List(context.Background(), model.Filter{PublishedOnly: true})
	if len(items) != 0 {
		t.Fatal("unpublished game listed")
	}
}
