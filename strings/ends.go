package strings

func EndsWith(s string, suffix string) bool {
	return (len(s) >= len(suffix)) && (s[len(s)-len(suffix):] == suffix)
}
