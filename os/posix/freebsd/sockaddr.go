package freebsd

/* From <sys/socket.h>. */
type Sockaddr struct {
	Len    uint8
	Family uint8
	Data   [14]byte
}

/* From <netinet/in.h>. */
type SockaddrIn struct {
	Len     uint8
	Family  uint8
	Port    uint16
	Address uint32
	_       [8]byte
}

const SOCK_MAXADDRLEN = 255 /* longest possible addresses */
