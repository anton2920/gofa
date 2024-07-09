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

func Memset[T any](mem []T, v T) {
	var i int
	for i = 0; i < len(mem)>>2; i += 4 {
		mem[i+0] = v
		mem[i+1] = v
		mem[i+2] = v
		mem[i+3] = v
	}
	for ; i < len(mem); i++ {
		mem[i] = v
	}
}
