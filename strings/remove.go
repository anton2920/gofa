package strings

func RemoveAt(vs []string, i int) []string {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}
