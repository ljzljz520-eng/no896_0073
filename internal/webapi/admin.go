package webapi

import (
	"encoding/json"
	"example.com/online-game-rank/internal/admin"
	"example.com/online-game-rank/internal/model"
	"net/http"
)

func (s *Server) AdminHandler(a *admin.Service) http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/admin/games", func(w http.ResponseWriter, r *http.Request) {
		var g model.Game
		if json.NewDecoder(r.Body).Decode(&g) != nil {
			http.Error(w, "bad json", 400)
			return
		}
		var e error
		if r.Method == http.MethodPost {
			_, e = a.ImportWithAudit(r.Context(), "manual", []model.Game{g})
		} else {
			e = a.SetDeviceAvailability(r.Context(), g.ID, g.WebEnabled, g.TabletEnabled, g.MobileEnabled)
		}
		if e != nil {
			http.Error(w, e.Error(), 400)
			return
		}
		write(w, map[string]string{"status": "accepted"})
	})
	return mux
}
