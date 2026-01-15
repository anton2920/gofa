package strings

import (
	"reflect"
	"unsafe"
)

//go:nosplit
func AsBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data, Len: len(s), Cap: len(s)}))
}

//go:nosplit
func Data(s string) *byte {
	return (*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))
}

func StartsEndsWith(s string, starts string, ends string) bool {
	return StartsWith(s, starts) && EndsWith(s, ends)
}
