package util

//go:nosplit
func RoundUp(x int, quantum int) int {
	return (x + (quantum - 1)) & ^(quantum - 1)
}
