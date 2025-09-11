package pointers

import "unsafe"

/* Add is the alternative to 'unsafe.Add' for Go versions which don't have it. */
func Add(ptr unsafe.Pointer, x int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(x))
}
