package syscall

import (
	"reflect"
	"unsafe"
)

const (
	/* From <stdio.h>. */
	SEEK_SET = 0
	SEEK_END = 2
)

var NULL = *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: 0, Len: 0}))
