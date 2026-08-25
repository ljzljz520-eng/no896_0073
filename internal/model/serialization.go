package model

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

func (g Game) Key() string         { return "game:" + g.ID }
func (p Player) Key() string       { return "player:" + p.ID }
func (s Score) Key() string        { return "score:" + s.ID }
func (r ImportRecord) Key() string { return "import:" + r.ID }
func EncodeGame(g Game) ([]byte, error) {
	if !g.Valid() {
		return nil, errors.New("invalid game")
	}
	return json.Marshal(g)
}
func DecodeGame(data []byte) (Game, error) {
	var g Game
	if len(data) == 0 {
		return g, errors.New("empty game")
	}
	e := json.Unmarshal(data, &g)
	if e != nil {
		return g, e
	}
	if !g.Valid() {
		return g, errors.New("invalid game")
	}
	return g, nil
}
func EncodePlayer(p Player) ([]byte, error) {
	if !p.Eligible() {
		return nil, errors.New("invalid player")
	}
	return json.Marshal(p)
}
func DecodePlayer(data []byte) (Player, error) {
	var p Player
	e := json.Unmarshal(data, &p)
	if e != nil {
		return p, e
	}
	if !p.Eligible() {
		return p, errors.New("invalid player")
	}
	return p, nil
}
func EncodeScore(s Score) ([]byte, error) {
	if !s.Valid() {
		return nil, errors.New("invalid score")
	}
	return json.Marshal(s)
}
func DecodeScore(data []byte) (Score, error) {
	var s Score
	e := json.Unmarshal(data, &s)
	if e != nil {
		return s, e
	}
	if !s.Valid() {
		return s, errors.New("invalid score")
	}
	return s, nil
}
func Slug(title string) string {
	parts := strings.Fields(strings.ToLower(title))
	if len(parts) == 0 {
		return ""
	}
	out := parts[0]
	for _, p := range parts[1:] {
		out += "-" + p
	}
	return out
}
func DisplayAge(min, max int) string {
	if min == max {
		return fmt.Sprintf("age %d", min)
	}
	return fmt.Sprintf("ages %d-%d", min, max)
}
func DeviceList(g Game) []string {
	out := []string{}
	if g.WebEnabled {
		out = append(out, "web")
	}
	if g.TabletEnabled {
		out = append(out, "tablet")
	}
	if g.MobileEnabled {
		out = append(out, "mobile")
	}
	return out
}
