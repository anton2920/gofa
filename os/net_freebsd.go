//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type (
	NetworkAddress  freebsd.Sockaddr
	InternetAddress freebsd.SockaddrIn
)

func (ia *InternetAddress) AsNetworkAddress() *NetworkAddress {
	return (*NetworkAddress)(unsafe.Pointer(ia))
}

func ListenForIncomingConnections(proto string, endpoint string) (Handle, error) {
	var buffer [freebsd.SOCK_MAXADDRLEN]byte
	var pf, semantics int32
	var addrLen uint32

	switch proto {
	case "tcp", "tcp/ip", "tcp/ip4", "tcp/ipv4":
		pf = freebsd.PF_INET
		semantics = freebsd.SOCK_STREAM

		addr, port, err := tcp.ParseEndpoint(endpoint)
		if err != nil {
			return -1, err
		}

		iaddr := (*InternetAddress)(unsafe.Pointer(&buffer))
		iaddr.Family = freebsd.AF_INET
		iaddr.Address = addr
		iaddr.Port = port

		addrLen = uint32(unsafe.Sizeof(*iaddr))
	default:
		return -1, errors.ErrNotImplemented
	}

	s, err := freebsd.Socket(pf, semantics, 0)
	if err != nil {
		return -1, err
	}

	if err := freebsd.Bind(s, (*freebsd.Sockaddr)(unsafe.Pointer(&buffer)), addrLen); err != nil {
		freebsd.Close(s)
		return -1, err
	}

	if (pf == freebsd.PF_INET) && (semantics == freebsd.SOCK_STREAM) {
		var enable int32 = 1
		if err := freebsd.Setsockopt(s, freebsd.SOL_SOCKET, freebsd.SO_REUSEPORT_LB, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable))); err != nil {
			freebsd.Close(s)
			return -1, err
		}

		if err := freebsd.Setsockopt(s, freebsd.IPPROTO_TCP, freebsd.TCP_NODELAY, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable))); err != nil {
			freebsd.Close(s)
			return -1, err
		}
	}

	if err := freebsd.Listen(s, 128); err != nil {
		freebsd.Close(s)
		return -1, err
	}

	return Handle(s), nil
}

func ConnectToAddress() (Handle, error) {
	return -1, errors.ErrNotImplemented
}

func AcceptIncomingConnection(s Handle, addr *NetworkAddress, addrLen *uint32) (Handle, error) {
	c, err := freebsd.Accept(int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
	return Handle(c), err
}
