package pointers

import (
	"unsafe"

	"github.com/anton2920/gofa/debug"
)

//go:nosplit
func AlignUp(ptr unsafe.Pointer, align uintptr) unsafe.Pointer {
	debug.Assert((align&(align-1)) == 0, "alignment must be a power of two")
	return unsafe.Pointer((uintptr(ptr) + (align - 1)) & ^(align - 1))
}
