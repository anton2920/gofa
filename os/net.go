package os

import (
	"unsafe"

	"github.com/anton2920/gofa/errors"
)

func BindSocketToAddressAndListenForIncomingConnections(s Handle, addr *NetworkAddress, addrLen uint32, backlog int) error {
	return errors.Or(BindSocketToAddress(s, addr, addrLen), ListenForIncomingConnections(s, backlog))
}

func SetSocketBooleanOption(s Handle, name SocketOptionName, enabled bool) error {
	var enable uint32
	if enabled {
		enable = 1
	}

	level, ok := SocketOptionName2SocketOptionLevel[name]
	if !ok {
		panic("no mapping between socket option name and socket option level")
	}

	return SetSocketOption(s, level, name, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable)))
}
