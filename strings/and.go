package strings

func And(vs ...string) string {
	var i int
	for i = 1; i < len(vs); i++ {
		if len(vs[i]) == 0 {
			break
		}
	}
	return vs[i-1]
}
