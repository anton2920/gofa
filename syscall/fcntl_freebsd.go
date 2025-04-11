package syscall

const (
	/* From <fcntl.h>. */
	O_RDONLY   = 0x0000
	O_WRONLY   = 0x0001
	O_RDWR     = 0x0002
	O_NONBLOCK = 0x0004
	O_CREAT    = 0x0200

	F_SETFL = 4
)
