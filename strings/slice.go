package strings

func SliceInsertAt(xs []string, s string, pos int) []string {
	xs = append(xs, s)
	copy(xs[pos+1:], xs[pos:])
	xs[pos] = s
	return xs
}

func SliceDeleteAt(xs []string, pos int) []string {
	if pos < len(xs)-1 {
		copy(xs[pos:], xs[pos+1:])
	}
	return xs[:len(xs)-1]
}
