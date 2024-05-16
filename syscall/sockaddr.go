package syscall

/* From <sys/socket.h>. */
type Sockaddr struct {
	Len    byte
	Family byte
	Data   [14]byte
}
