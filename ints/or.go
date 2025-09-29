package ints

func Or(vs ...int) int {
	for i := 0; i < len(vs); i++ {
		if vs[i] > 0 {
			return vs[i]
		}
	}
	return 0
}
