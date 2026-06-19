package pointers

import "unsafe"

/* Noescape hides a pointer from escape analysis. Noescape is the identity function but escape analysis doesn't think the output depends on the input. Noescape is inlined and currently compiles down to zero instructions. */
//go:nosplit
func UnsafeNoescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

//go:nosplit
func ByteNoescape(p *byte) *byte {
	return (*byte)(UnsafeNoescape(unsafe.Pointer(p)))
}
