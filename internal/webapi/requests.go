package webapi

import (
	"example.com/online-game-rank/internal/model"
	"net/http"
	"strconv"
	"strings"
)

func filterFromRequest(r *http.Request) model.Filter {
	q := r.URL.Query()
	age, _ := strconv.Atoi(q.Get("age"))
	published := q.Get("published") != "false"
	return model.Filter{Category: strings.ToLower(q.Get("category")), Age: age, Device: strings.ToLower(q.Get("device")), Search: q.Get("q"), PublishedOnly: published}
}
func acceptsJSON(r *http.Request) bool {
	return strings.Contains(r.Header.Get("Accept"), "application/json") || r.Header.Get("Accept") == ""
}
func parseLimit(r *http.Request, defaultValue, max int) int {
	n, e := strconv.Atoi(r.URL.Query().Get("limit"))
	if e != nil || n <= 0 {
		n = defaultValue
	}
	if n > max {
		n = max
	}
	return n
}
func requestID(r *http.Request) string {
	if x := r.Header.Get("X-Request-ID"); x != "" {
		return x
	}
	return "anonymous"
}
func setCache(w http.ResponseWriter, seconds int) {
	if seconds < 0 {
		seconds = 0
	}
	w.Header().Set("Cache-Control", "public, max-age="+strconv.Itoa(seconds))
}
func wantsHTML(r *http.Request) bool { return strings.Contains(r.Header.Get("Accept"), "text/html") }
