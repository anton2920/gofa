package syscall

const (
	/* From <fcntl.h>. */
	O_RDONLY   = 00000
	O_WRONLY   = 00001
	O_RDWR     = 00002
	O_CREAT    = 00100
	O_NONBLOCK = 04000

	F_SETFL = 4
)
