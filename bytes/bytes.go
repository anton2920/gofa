package bytes

import "unsafe"

//go:nosplit
func AsString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
