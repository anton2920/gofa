package strings

func And(vs ...string) string {
	for i := 1; i < len(vs); i++ {
		if len(vs[i]) == 0 {
			return vs[i-1]
		}
	}
	return vs[0]
}
