package util

import "unsafe"

//go:nosplit
func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

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

func RemoveAtIndex[T any](ts []T, i int) []T {
	if (len(ts) == 0) || (i < 0) || (i >= len(ts)) {
		return ts
	}

	if i < len(ts)-1 {
		copy(ts[i:], ts[i+1:])
	}
	return ts[:len(ts)-1]
}

//go:nosplit
func SwapBytesInWord(x uint16) uint16 {
	return ((x << 8) & 0xFF00) | (x >> 8)
}
