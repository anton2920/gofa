package ints

func InsertAt(vs []int, v int, i int) []int {
	vs = append(vs, 0)
	copy(vs[i+1:], vs[i:])
	vs[i] = v
	return vs
}
