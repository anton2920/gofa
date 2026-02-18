package strings

func DeleteAt(xs []string, pos int) []string {
	if pos < len(xs)-1 {
		copy(xs[pos:], xs[pos+1:])
	}
	return xs[:len(xs)-1]
}
