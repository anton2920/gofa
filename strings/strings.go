package strings

import (
	"unsafe"

	"github.com/anton2920/gofa/go/types"
)

//go:nosplit
func AsBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&types.SliceHeader{Data: uintptr(unsafe.Pointer(Data(s))), Len: len(s), Cap: len(s)}))
}

//go:nosplit
func Data(s string) *byte {
	return (*byte)(unsafe.Pointer((*types.StringHeader)(unsafe.Pointer(&s)).Data))
}

func StartsEndsWith(s string, starts string, ends string) bool {
	return (len(s) >= (len(starts) + len(ends))) && (StartsWith(s, starts)) && (EndsWith(s, ends))
}
