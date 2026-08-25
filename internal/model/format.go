package model

import (
	"strings"
	"time"
)

func CleanTitle(v string) string { return strings.Join(strings.Fields(strings.TrimSpace(v)), " ") }
func CleanURL(v string) string   { return strings.TrimSpace(strings.TrimRight(v, "/")) }
func AgeRange(min, max int) (int, int) {
	if min < 3 {
		min = 3
	}
	if max > 18 {
		max = 18
	}
	if max < min {
		max = min
	}
	return min, max
}
func NewGame(id, title, category, url string, min, max int) Game {
	min, max = AgeRange(min, max)
	return Game{ID: id, Title: CleanTitle(title), Category: strings.ToLower(strings.TrimSpace(category)), EntryURL: CleanURL(url), AgeMin: min, AgeMax: max, CreatedAt: time.Unix(0, 0), UpdatedAt: time.Unix(0, 0)}
}
func (g Game) WithDevices(web, tablet, mobile bool) Game {
	g.WebEnabled = web
	g.TabletEnabled = tablet
	g.MobileEnabled = mobile
	return g
}
func (g Game) WithDescription(v string) Game { g.Description = CleanTitle(v); return g }
func (g Game) WithPublished(v bool) Game     { g.Published = v; return g }
func (g Game) Touch() Game                   { g.Version++; g.UpdatedAt = time.Unix(int64(g.Version), 0); return g }
func (p Player) WithActive(v bool) Player    { p.Active = v; return p }
func NewPlayer(id, name string, age int) Player {
	return Player{ID: id, DisplayName: CleanTitle(name), Age: age, JoinedAt: time.Unix(0, 0), Active: true}
}
