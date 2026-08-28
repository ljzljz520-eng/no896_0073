package main

import (
	"example.com/online-game-rank/internal/admin"
	"example.com/online-game-rank/internal/catalog"
	"example.com/online-game-rank/internal/query"
	"example.com/online-game-rank/internal/ranking"
	"example.com/online-game-rank/internal/store"
	"example.com/online-game-rank/internal/webapi"
	"log"
	"net/http"
	"os"
)

func main() {
	path := os.Getenv("GAME_RANK_DB")
	if path == "" {
		path = "games.db"
	}
	db, e := store.Open(path)
	if e != nil {
		log.Fatal(e)
	}
	defer db.Close()
	c := catalog.New(db)
	r := ranking.New(db)
	q := query.New(c, r)
	a := admin.New(db, c)
	api := webapi.New(q)
	mux := http.NewServeMux()
	mux.Handle("/", api.Handler())
	mux.Handle("/admin/", audit(api, a))
	log.Println(http.ListenAndServe(":8080", mux))
}
func audit(api *webapi.Server, a *admin.Service) http.Handler { return api.AdminHandler(a) }
