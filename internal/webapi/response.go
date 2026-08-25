package webapi

import (
	"encoding/json"
	"net/http"
)

type ErrorResponse struct {
	Error string `json:"error"`
	Code  string `json:"code"`
}

func writeError(w http.ResponseWriter, status int, code, detail string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorResponse{Error: detail, Code: code})
}
func methodAllowed(w http.ResponseWriter, ok string, r *http.Request) bool {
	if r.Method != ok {
		w.Header().Set("Allow", ok)
		writeError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return false
	}
	return true
}
