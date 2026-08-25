package model

import "time"

type ImportRecord struct {
	ID         string    `json:"id"`
	Source     string    `json:"source"`
	Received   int       `json:"received"`
	Accepted   int       `json:"accepted"`
	Rejected   int       `json:"rejected"`
	Status     string    `json:"status"`
	StartedAt  time.Time `json:"started_at"`
	FinishedAt time.Time `json:"finished_at"`
	Error      string    `json:"error,omitempty"`
}

func (r ImportRecord) Complete() bool { return r.Status == "completed" || r.Status == "failed" }

type Filter struct {
	Category      string
	Age           int
	Device        string
	Search        string
	PublishedOnly bool
}

func (f Filter) Matches(g Game) bool {
	if f.Category != "" && g.Category != f.Category {
		return false
	}
	if f.Age > 0 && !g.FitsAge(f.Age) {
		return false
	}
	if f.Device != "" && !g.IsPlayable(f.Device) {
		return false
	}
	if f.PublishedOnly && !g.Published {
		return false
	}
	if f.Search != "" && !containsFold(g.Title, f.Search) {
		return false
	}
	return true
}

func containsFold(text, query string) bool {
	for i := 0; i+len(query) <= len(text); i++ {
		ok := true
		for j := range query {
			a, b := text[i+j], query[j]
			if a >= 'A' && a <= 'Z' {
				a += 'a' - 'A'
			}
			if b >= 'A' && b <= 'Z' {
				b += 'a' - 'A'
			}
			if a != b {
				ok = false
				break
			}
		}
		if ok {
			return true
		}
	}
	return false
}
