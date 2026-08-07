package freebsd

/* From <netinet/in.h>. */
type InAddr struct {
	Addr uint32
}

type SockaddrIn struct {
	Len    uint8
	Family uint8
	Port   uint16
	Addr   InAddr
	_      [8]byte
}

/* Protocols common to RFC 1700, POSIX, and X/Open. */
const (
	IPPROTO_IP   = 0  /* dummy for IP */
	IPPROTO_ICMP = 1  /* control message protocol */
	IPPROTO_TCP  = 6  /* tcp */
	IPPROTO_UDP  = 17 /* user datagram protocol */

	INADDR_ANY       = uint32(0x00000000)
	INADDR_BROADCAST = uint32(0xffffffff) /* must be masked */
)
