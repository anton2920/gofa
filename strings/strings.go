package strings

import (
	"reflect"
	"unsafe"
)

func AsBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data, Len: len(s), Cap: len(s)}))
}

func Data(s string) *byte {
	return (*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))
}
