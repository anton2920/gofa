package syscall

import (
	"reflect"
	"unsafe"
)

/* NOTE(anton2920): this is basically a Go's string type. */
/* From <sys/_iovec.h>. */
/*
 * struct iovec {
 *	void	*iov_base;
 *	size_t	iov_len;
 * };
 */
type Iovec string

var IovecZ = *(*Iovec)(unsafe.Pointer(&reflect.StringHeader{Data: 0, Len: 0}))

func IovecForByteSlice(buf []byte) Iovec {
	if buf == nil {
		return IovecZ
	}
	return *(*Iovec)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&buf[0])), Len: len(buf)}))
}
