package strings

import "strings"

/* TODO(anton2920): rewrite using SIMD. */
func FindChar(s string, c byte) int {
	return strings.IndexByte(s, c)
}

/* TODO(anton2920): rewrite using SIMD. */
func FindSubstring(a, b string) int {
	return strings.Index(a, b)
}

/* TODO(anton2920): rewrite using (DF=1 and REP SCASB) or SIMD. */
func FindCharReverse(s string, c byte) int {
	for i := len(s) - 1; i >= 0; i-- {
		if s[i] == c {
			return i
		}
	}
	return -1
}
