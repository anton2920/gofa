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

func (ia *InternetAddress) Length() uint32 {
	return uint32(unsafe.Sizeof(*ia))
}

func ListenForIncomingConnections(proto string, endpoint string) (Handle, error) {
	var buffer [freebsd.SOCK_MAXADDRLEN]byte
	var pf, semantics int32
	var addrLen uint32

	switch proto {
	case "tcp/ip":
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

		addrLen = iaddr.Length()
	default:
		return -1, errors.ErrNotImplemented
	}

	s, err := freebsd.Socket(pf, semantics, 0)
	if err != nil {
		return -1, err
	}

	if err := freebsd.Bind(s, (*freebsd.Sockaddr)(unsafe.Pointer(&buffer)), addrLen); err != nil {
		return -1, err
	}

	/* TODO(anton2920): set socket options. */

	return Handle(s), nil
}

func ConnectToAddress() (Handle, error) {
	return -1, errors.ErrNotImplemented
}

func AcceptIncomingConnection(s Handle, addr *NetworkAddress, addrLen *uint32) (Handle, error) {
	return -1, errors.ErrNotImplemented
}
