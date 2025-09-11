package ints

func AlignUp(x int, quantum int) int {
	return (x + (quantum - 1)) & ^(quantum - 1)
}
