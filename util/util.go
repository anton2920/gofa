package util

import "unsafe"

/* NOTE(anton2920): Noescape hides a pointer from escape analysis. Noescape is the identity function but escape analysis doesn't think the output depends on the input. Noescape is inlined and currently compiles down to zero instructions. */
//go:nosplit
func Noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

//go:nosplit
func RoundUp(x int, quantum int) int {
	return (x + (quantum - 1)) & ^(quantum - 1)
}
