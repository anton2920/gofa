package syscall

/* From <sys/socket.h>. */
type Sockaddr struct {
	Len    byte
	Family byte
	Data   [14]byte
}

/* From <netinet/in.h>. */
type SockAddrIn struct {
	Len    uint8
	Family uint8
	Port   uint16
	Addr   uint32
	_      [8]byte
}
