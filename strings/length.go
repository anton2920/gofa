package strings

import "unicode/utf8"

func LengthInRange(s string, min, max int) bool {
	return (utf8.RuneCountInString(s) >= min) && (utf8.RuneCountInString(s) <= max)
}
