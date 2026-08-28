package model

import "time"

type Game struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	AgeMin        int       `json:"age_min"`
	AgeMax        int       `json:"age_max"`
	WebEnabled    bool      `json:"web_enabled"`
	TabletEnabled bool      `json:"tablet_enabled"`
	MobileEnabled bool      `json:"mobile_enabled"`
	EntryURL      string    `json:"entry_url"`
	Published     bool      `json:"published"`
	Version       int       `json:"version"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func (g Game) IsPlayable(device string) bool {
	if !g.Published {
		return false
	}
	switch device {
	case "web":
		return g.WebEnabled
	case "tablet":
		return g.TabletEnabled
	case "mobile":
		return g.MobileEnabled
	default:
		return false
	}
}

func (g Game) FitsAge(age int) bool { return age >= g.AgeMin && age <= g.AgeMax }

func (g Game) Valid() bool {
	if g.ID == "" || g.Title == "" || g.EntryURL == "" {
		return false
	}
	if g.AgeMin < 3 || g.AgeMax < g.AgeMin || g.AgeMax > 18 {
		return false
	}
	switch g.Category {
	case "math", "english", "science", "logic":
		return true
	default:
		return false
	}
}
