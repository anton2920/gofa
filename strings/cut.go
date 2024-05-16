package strings

/* Cut slices s around the first instance of sep, returning the text before and after sep. The found result reports whether sep appears in s. If sep does not appear in s, cut returns s, "", false. */
func Cut(s, sep string) (before, after string, found bool) {
	if i := FindSubstring(s, sep); i >= 0 {
		return s[:i], s[i+len(sep):], true
	}
	return s, "", false
}
