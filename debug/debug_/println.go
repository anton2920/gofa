//go:build gofadebug
// +build gofadebug

package debug_

import "unsafe"

func Println(xs ...string) {
	var n int

	buf := make([]byte, 4096)
	for i := 0; i < len(xs); i++ {
		n += copy(buf[n:], xs[i])
	}
	if buf[n-1] != '\n' {
		buf[n] = '\n'
		n++
	}

	print((*(*string)(unsafe.Pointer(&buf)))[:n])
}
