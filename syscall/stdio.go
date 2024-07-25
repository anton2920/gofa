package syscall

import "unsafe"

const (
	/* From <stdio.h>. */
	SEEK_SET = 0
	SEEK_END = 2
)

var NULL = unsafe.String(nil, 0)
