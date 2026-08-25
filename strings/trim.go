package strings

/* TODO(anton2920): add support for UTF-8. */
func Trim(s string, cutset string) string {
	isInSet := func(ch byte, set string) bool {
		for _, r := range set {
			if rune(ch) == r {
				return true
			}
		}
		return false
	}

	var i int
	for i = 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			panic("UTF-8 is not supported")
		}
		if !isInSet(s[i], cutset) {
			break
		}
	}

	var j int
	for j = len(s) - 1; j > i; j-- {
		if s[i] >= 0x80 {
			panic("UTF-8 is not supported")
		}
		if !isInSet(s[j], cutset) {
			break
		}
	}

	return s[i : j+1]
}

/* TODO(anton2920): add support for UTF-8. */
func TrimSpace(s string) string {
	isSpace := func(ch byte) bool {
		return (ch == ' ') || (ch == '\t') || (ch == '\n') || (ch == '\r')
	}

	var i int
	for i = 0; i < len(s); i++ {
		if s[i] >= 0x80 {
			panic("UTF-8 is not supported")
		}
		if !isSpace(s[i]) {
			break
		}
	}

	var j int
	for j = len(s) - 1; j > i; j-- {
		if s[i] >= 0x80 {
			panic("UTF-8 is not supported")
		}
		if !isSpace(s[j]) {
			break
		}
	}

	return s[i : j+1]
}
