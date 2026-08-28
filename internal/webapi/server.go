package webapi

import (
	"encoding/json"
	"example.com/online-game-rank/internal/model"
	"example.com/online-game-rank/internal/query"
	"net/http"
	"strconv"
)

type Server struct{ query *query.Service }

func New(q *query.Service) *Server { return &Server{query: q} }
func (s *Server) Handler() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/health", s.health)
	mux.HandleFunc("/api/games", s.games)
	mux.HandleFunc("/api/ranking", s.ranking)
	return mux
}
func (s *Server) health(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok"}`))
}
func (s *Server) games(w http.ResponseWriter, r *http.Request) {
	f := model.Filter{Category: r.URL.Query().Get("category"), Device: r.URL.Query().Get("device"), Search: r.URL.Query().Get("q"), PublishedOnly: true}
	if a, e := strconv.Atoi(r.URL.Query().Get("age")); e == nil {
		f.Age = a
	}
	items, e := s.query.Browse(r.Context(), f)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	write(w, items)
}
func (s *Server) ranking(w http.ResponseWriter, r *http.Request) {
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	items, e := s.query.Top(r.Context(), r.URL.Query().Get("game"), limit)
	if e != nil {
		http.Error(w, e.Error(), http.StatusBadRequest)
		return
	}
	write(w, items)
}
func write(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
