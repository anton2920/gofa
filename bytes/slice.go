package bytes

import (
	"unsafe"

	"github.com/anton2920/gofa/pointers"
)

type slice struct {
	Data uintptr
	Len  int
	Cap  int
}

//go:nosplit
func SliceFromBytePointer(ptr *byte, n int) []byte {
	return SliceFromUnsafePointer(unsafe.Pointer(ptr), n)
}

//go:nosplit
func SliceFromUnsafePointer(ptr unsafe.Pointer, n int) []byte {
	return *((*[]byte)(unsafe.Pointer(&slice{Data: uintptr(pointers.UnsafeNoescape(ptr)), Len: n, Cap: n})))
}
