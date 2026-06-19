package pointers

import "unsafe"

//go:nosplit
func Add(ptr unsafe.Pointer, x uintptr) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + x)
}

//go:nosplit
func Delta(p unsafe.Pointer, q unsafe.Pointer) int {
	return int(uintptr(p) - uintptr(q))
}

//go:nosplit
func Sub(ptr unsafe.Pointer, x uintptr) unsafe.Pointer {
	return Add(ptr, ^x+1)
}
