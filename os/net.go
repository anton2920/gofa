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
	return SetSocketOption(ctx, s, SocketOptionName2SocketOptionLevel(name), name, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable)))
}
