package os

import (
	"unsafe"

	"github.com/anton2920/gofa/context"
)

func BindSocketToAddressAndListenForIncomingConnections(ctx *context.Context, s Handle, addr *NetworkAddress, addrLen uint32, backlog int) bool {
	return BindSocketToAddress(ctx, s, addr, addrLen) || ListenForIncomingConnections(ctx, s, backlog)
}

func SetSocketBooleanOption(ctx *context.Context, s Handle, name SocketOptionName, enabled bool) bool {
	var enable uint32
	if enabled {
		enable = 1
	}

	level, ok := SocketOptionName2SocketOptionLevel[name]
	if !ok {
		panic("no mapping between socket option name and socket option level")
	}

	return SetSocketOption(ctx, s, level, name, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable)))
}
