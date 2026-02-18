package syscall

const (
	/* From <fcntl.h>. */
	O_RDONLY   = 0x0000 /* open for reading only */
	O_WRONLY   = 0x0001 /* open for writing only */
	O_RDWR     = 0x0002 /* open for reading and writing */
	O_NONBLOCK = 0x0004 /* no delay */
	O_APPEND   = 0x0008 /* set append mode */
	O_CREAT    = 0x0200 /* create if nonexistent */
	O_TRUNC    = 0x0400 /* truncate to zero length */

	F_GETFL = 3 /* get file status flags */
	F_SETFL = 4 /* set file status flags */
)
