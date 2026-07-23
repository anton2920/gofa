package bytes

import (
	"unsafe"

	"github.com/anton2920/gofa/go/types"
	"github.com/anton2920/gofa/pointers"
)

var zeroSlice = make([]byte, 4096)

//go:nosplit
func SliceClear(buf []byte) {
	var n int
	for n < len(buf) {
		n += copy(buf[n:], zeroSlice)
	}
}

//go:nosplit
func SliceFromBytePointer(ptr *byte, n int) []byte {
	return SliceFromUnsafePointer(unsafe.Pointer(ptr), n)
}

//go:nosplit
func SliceFromUnsafePointer(ptr unsafe.Pointer, n int) []byte {
	return *((*[]byte)(unsafe.Pointer(&types.SliceHeader{Data: uintptr(pointers.UnsafeNoescape(ptr)), Len: n, Cap: n})))
}
