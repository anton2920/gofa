package ints

/* Clamp returns number clamped into a range from l inclusive to r exclusive. */
func Clamp(x int, l int, r int) int {
	if x > r-1 {
		x = r - 1
	}
	if x < l {
		x = l
	}
	return x
}
