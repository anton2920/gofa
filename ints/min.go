package ints

func Min(a int, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func Mins(xs ...int) int {
	min := xs[0]
	for i := 1; i < len(xs); i++ {
		if xs[i] < min {
			min = xs[i]
		}
	}
	return min
}
