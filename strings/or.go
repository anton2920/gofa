package strings

func Or(vs ...string) string {
	for i := 0; i < len(vs); i++ {
		if len(vs[i]) > 0 {
			return vs[i]
		}
	}
	return ""
}
