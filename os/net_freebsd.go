//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

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

func (ia *InternetAddress) AsNetworkAddress() *NetworkAddress {
	return (*NetworkAddress)(unsafe.Pointer(ia))
}

func CreateNetworkSocket(pf ProtocolFamily, typ SocketType, proto Protocol) (Handle, error) {
	s, err := freebsd.Socket(int32(pf), int32(typ), int32(proto))
	return Handle(s), err
}

func BindSocketToAddress(s Handle, addr *NetworkAddress, addrLen uint32) error {
	return freebsd.Bind(int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
}

func ListenForIncomingConnections(s Handle, backlog int) error {
	return freebsd.Listen(int32(s), int32(backlog))
}

func AcceptIncomingConnection(s Handle, addr *NetworkAddress, addrLen *uint32) (Handle, error) {
	c, err := freebsd.Accept(int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
	return Handle(c), err
}

func ConnectToAddress(s Handle, addr *NetworkAddress, addrLen uint32) error {
	return freebsd.Connect(int32(s), (*freebsd.Sockaddr)(unsafe.Pointer(addr)), addrLen)
}

func SetSocketOption(s Handle, level SocketOptionLevel, name SocketOptionName, val unsafe.Pointer, valLen uint32) error {
	return freebsd.Setsockopt(int32(s), int32(level), int32(name), val, valLen)
}
