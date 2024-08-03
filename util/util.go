package util

import "unsafe"

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

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

/* GetCallerPC returns a value of %IP register that is going to be used by RET instruction. arg0 is the address of the first agrument function of interest accepts. */
func GetCallerPC(arg0 unsafe.Pointer) uintptr {
	return *((*uintptr)(unsafe.Add(arg0, -8)))
}

/* Noescape hides a pointer from escape analysis. Noescape is the identity function but escape analysis doesn't think the output depends on the input. Noescape is inlined and currently compiles down to zero instructions. */
//go:nosplit
func Noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func RoundUp(x int, quantum int) int {
	return (x + (quantum - 1)) & ^(quantum - 1)
}

func Memset[T any](mem []T, v T) {
	var i int
	for i = 0; i < len(mem)>>2; i++ {
		mem[(i<<2)+0] = v
		mem[(i<<2)+1] = v
		mem[(i<<2)+2] = v
		mem[(i<<2)+3] = v
	}
	for j := 0; j < len(mem)&3; j++ {
		mem[(i<<2)+j] = v
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

func SwapBytesInWord(x uint16) uint16 {
	return ((x << 8) & 0xFF00) | (x >> 8)
}
