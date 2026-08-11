package strings

import "strings"

func Trim(s string, cutset string) string {
	return strings.Trim(s, cutset)
}

func TrimSpace(s string) string {
	return strings.TrimSpace(s)
}
