package strings

func InsertAt(xs []string, s string, pos int) []string {
	xs = append(xs, s)
	copy(xs[pos+1:], xs[pos:])
	xs[pos] = s
	return xs
}
