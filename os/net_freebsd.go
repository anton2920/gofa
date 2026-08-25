//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type AddressFamily uint8

const (
	AddressFamilyInternet = AddressFamily(freebsd.AF_INET)
)

type NetworkAddress struct {
	Len    uint8
	Family AddressFamily
	Data   [14]byte
}

type InternetAddress struct {
	Len     uint8
	Family  AddressFamily
	Port    uint16
	Address uint32
	_       [8]byte
}

type ProtocolFamily int32

const (
	ProtocolFamilyInternet = ProtocolFamily(freebsd.PF_INET)
)

type SocketType int32

const (
	SocketTypeStream = SocketType(freebsd.SOCK_STREAM)
)

type Protocol int32

type SocketOptionLevel int32

const (
	SocketOptionLevelSocket = SocketOptionLevel(freebsd.SOL_SOCKET)
	SocketOptionLevelTCP    = SocketOptionLevel(freebsd.IPPROTO_TCP)
)

type SocketOptionName int32

const (
	SocketOptionReuseLocalAddress                         = SocketOptionName(freebsd.SO_REUSEADDR)
	SocketOptionReuseLocalAddressAndPort                  = SocketOptionName(freebsd.SO_REUSEPORT)
	SocketOptionReuseLocalAddressAndPortWithLoadBalancing = SocketOptionName(freebsd.SO_REUSEPORT_LB)

	SocketOptionTCPNoDelay = SocketOptionName(freebsd.TCP_NODELAY)
)

var SocketOptionName2SocketOptionLevel = map[SocketOptionName]SocketOptionLevel{
	SocketOptionReuseLocalAddress:                         SocketOptionLevelSocket,
	SocketOptionReuseLocalAddressAndPort:                  SocketOptionLevelSocket,
	SocketOptionReuseLocalAddressAndPortWithLoadBalancing: SocketOptionLevelSocket,

	SocketOptionTCPNoDelay: SocketOptionLevelTCP,
}

func (ia *InternetAddress) AsNetworkAddressPointer() *NetworkAddress {
	return (*NetworkAddress)(unsafe.Pointer(ia))
}

func CreateNetworkSocket(ctx *context.Context, pf ProtocolFamily, typ SocketType, proto Protocol) (Handle, bool) {
	s, ok := freebsd.Socket(ctx, int32(pf), int32(typ), int32(proto))
	return Handle(s), ok
}

func BindSocketToAddress(ctx *context.Context, s Handle, addr *NetworkAddress, addrLen uint32) bool {
	return freebsd.Bind(ctx, int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
}

func ListenForIncomingConnections(ctx *context.Context, s Handle, backlog int) bool {
	return freebsd.Listen(ctx, int32(s), int32(backlog))
}

func AcceptIncomingConnection(ctx *context.Context, s Handle, addr *NetworkAddress, addrLen *uint32) (Handle, bool) {
	c, ok := freebsd.Accept(ctx, int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
	return Handle(c), ok
}

func ConnectToAddress(ctx *context.Context, s Handle, addr *NetworkAddress, addrLen uint32) bool {
	return freebsd.Connect(ctx, int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
}

func SetSocketOption(ctx *context.Context, s Handle, level SocketOptionLevel, name SocketOptionName, val unsafe.Pointer, valLen uint32) bool {
	return freebsd.Setsockopt(ctx, int32(s), int32(level), int32(name), val, valLen)
}
