package bytes

import "unsafe"

func AsString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
