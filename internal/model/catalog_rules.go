package model

var Categories = []string{"math", "english", "science", "logic"}
var Devices = []string{"web", "tablet", "mobile"}

func CategoryAllowed(v string) bool {
	for _, x := range Categories {
		if x == v {
			return true
		}
	}
	return false
}
func DeviceAllowed(v string) bool {
	for _, x := range Devices {
		if x == v {
			return true
		}
	}
	return false
}
func NormalizeCategory(v string) string {
	for _, x := range Categories {
		if len(v) == len(x) && lower(v) == x {
			return x
		}
	}
	return ""
}
func lower(v string) string {
	b := []byte(v)
	for i, c := range b {
		if c >= 'A' && c <= 'Z' {
			b[i] = c + 32
		}
	}
	return string(b)
}
