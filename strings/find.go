package strings

/* TODO(anton2920): rewrite using SIMD. */
func FindChar(s string, c byte) int {
	for i := 0; i < len(s); i++ {
		if s[i] == c {
			return i
		}
	}
	return -1
}

/* TODO(anton2920): rewrite using SIMD. */
func FindSubstring(haystack string, needle string) int {
	for i := 0; i < len(haystack)-len(needle); i++ {
		if StartsWith(haystack[i:], needle) {
			return i
		}
	}
	return -1
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
