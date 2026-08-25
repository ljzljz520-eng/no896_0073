package main_test

import (
	"example.com/online-game-rank/internal/model"
	"testing"
)

func TestFilterSearch(t *testing.T) {
	g := model.Game{Title: "Number Quest", Category: "math", AgeMin: 7, AgeMax: 10, EntryURL: "x", Published: true, WebEnabled: true}
	if !(model.Filter{Search: "quest"}).Matches(g) {
		t.Fatal("case insensitive search failed")
	}
}
