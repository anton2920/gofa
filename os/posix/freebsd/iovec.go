package freebsd

import "github.com/anton2920/gofa/bytes"

/* NOTE(anton2920): this is basically a Go's string type. */
/* From <sys/_iovec.h>. */
/*
 * struct iovec {
 *	void	*iov_base;
 *	size_t	iov_len;
 * };
 */
type Iovec string

var IovecZ = Iovec("")

func IovecForByteSlice(buf []byte) Iovec {
	if buf == nil {
		return IovecZ
	}
	return Iovec(bytes.AsString(buf))
}
