package main_test

import (
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/query"
	"example.com/online-game-rank/internal/ranking"
	"example.com/online-game-rank/internal/store"
	"example.com/online-game-rank/internal/webapi"
	"net/http/httptest"
	"testing"
)

func TestHealthEndpoint(t *testing.T) {
	s, _ := store.Open(t.TempDir() + "/f.db")
	defer s.Close()
	q := query.New(catalog.New(s), ranking.New(s))
	rr := httptest.NewRecorder()
	req := httptest.NewRequest("GET", "/health", nil)
	webapi.New(q).Handler().ServeHTTP(rr, req)
	if rr.Code != 200 {
		t.Fatalf("status %d", rr.Code)
	}
}
