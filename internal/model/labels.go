package model

func CategoryLabel(c string) string {
	switch c {
	case "math":
		return "Mathematics"
	case "english":
		return "English"
	case "science":
		return "Science"
	case "logic":
		return "Logic"
	default:
		return "Unknown"
	}
}
func DeviceLabel(d string) string {
	switch d {
	case "web":
		return "Web browser"
	case "tablet":
		return "Tablet"
	case "mobile":
		return "Mobile"
	default:
		return "Unknown device"
	}
}
func AgeBandLabel(a int) string {
	if a < 3 {
		return "not eligible"
	}
	if a <= 5 {
		return "early learners"
	}
	if a <= 8 {
		return "young learners"
	}
	if a <= 12 {
		return "middle learners"
	}
	if a <= 15 {
		return "teen learners"
	}
	if a <= 18 {
		return "older learners"
	}
	return "not eligible"
}
func ScoreLabel(p int) string {
	if p < 0 {
		return "invalid"
	}
	if p < 25 {
		return "starting"
	}
	if p < 50 {
		return "developing"
	}
	if p < 75 {
		return "strong"
	}
	if p < 100 {
		return "excellent"
	}
	return "mastery"
}
func DurationLabel(s int) string {
	if s <= 0 {
		return "none"
	}
	m := s / 60
	r := s % 60
	if m == 0 {
		return "under a minute"
	}
	if r == 0 {
		return plural(m, "minute")
	}
	return plural(m, "minute") + " " + plural(r, "second")
}
func plural(n int, w string) string {
	if n == 1 {
		return "1 " + w
	}
	return fmtInt(n) + " " + w + "s"
}
func fmtInt(n int) string {
	if n == 0 {
		return "0"
	}
	d := []byte{}
	for n > 0 {
		d = append([]byte{byte('0' + n%10)}, d...)
		n /= 10
	}
	return string(d)
}
func EligibleAges(min, max int) []int {
	if min < 3 {
		min = 3
	}
	if max > 18 {
		max = 18
	}
	o := []int{}
	for a := min; a <= max; a++ {
		o = append(o, a)
	}
	return o
}
func DeviceMatrix(g Game) map[string]bool {
	return map[string]bool{"web": g.WebEnabled && g.Published, "tablet": g.TabletEnabled && g.Published, "mobile": g.MobileEnabled && g.Published}
}
func CloneGame(g Game) Game {
	return Game{ID: g.ID, Title: g.Title, Description: g.Description, Category: g.Category, AgeMin: g.AgeMin, AgeMax: g.AgeMax, WebEnabled: g.WebEnabled, TabletEnabled: g.TabletEnabled, MobileEnabled: g.MobileEnabled, EntryURL: g.EntryURL, Published: g.Published, Version: g.Version, CreatedAt: g.CreatedAt, UpdatedAt: g.UpdatedAt}
}
