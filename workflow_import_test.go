package main_test

import (
	"context"
	"example.com/online-game-rank/internal/admin"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/store"
	"testing"
)

func TestWorkflowThree(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/c.db")
	defer s.Close()
	c := catalog.New(s)
	a := admin.New(s, c)
	rows := []model.Game{{ID: "e1", Title: "Words", Category: "english", AgeMin: 6, AgeMax: 9, EntryURL: "/words", Published: true}, {ID: "bad", Title: "", Category: "english"}}
	r, e := a.ImportWithAudit(context.Background(), "fixture", rows)
	if e != nil || r.Accepted != 1 || r.Rejected != 1 {
		t.Fatalf("import %#v %v", r, e)
	}
}
