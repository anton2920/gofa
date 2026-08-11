package strings

func ToInt(s string) int {
	var sign int

	switch s[0] {
	case '-':
		sign = -1
		fallthrough
	case '+':
		s = s[1:]
	}

	var n int
	for len(s) > 0 {
		n = n*10 + int(s[0]-'0')
		s = s[1:]
	}

	return sign * n
}
