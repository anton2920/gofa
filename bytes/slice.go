package bytes

import (
	"unsafe"

	"github.com/anton2920/gofa/go/types"
)

//go:nosplit
func SliceFromBytePointer(ptr *byte, n int) []byte {
	return SliceFromUnsafePointer(unsafe.Pointer(ptr), n)
}

//go:nosplit
func SliceFromUnsafePointer(ptr unsafe.Pointer, n int) []byte {
	return *((*[]byte)(unsafe.Pointer(&types.SliceHeader{Data: uintptr(ptr), Len: n, Cap: n})))
}
