package syscall

import "unsafe"

/* NOTE(anton2920): this is basically a Go's string type. */
/* From <sys/_iovec.h>. */
/*
 * struct iovec {
 *	void	*iov_base;
 *	size_t	iov_len;
 * };
 */
type Iovec string

var IovecZ = Iovec(unsafe.String(nil, 0))

func IovecForByteSlice(buf []byte) Iovec {
	return Iovec(unsafe.String(unsafe.SliceData(buf), len(buf)))
}
