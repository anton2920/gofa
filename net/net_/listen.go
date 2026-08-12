package net_

import (
	"unsafe"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/os"
)

func Listen(ctx *context.Context, proto string, endpoint string) (os.Handle, bool) {
	var addrBuf [unsafe.Sizeof(os.NetworkAddress{})]byte
	var pf os.ProtocolFamily
	var typ os.SocketType
	var prot os.Protocol
	var addrLen uint32

	switch proto {
	case "tcp", "tcp/ip", "tcp/ip4", "tcp/ipv4":
		pf = os.ProtocolFamilyInternet
		typ = os.SocketTypeStream

		addr, port, ok := tcp.ParseEndpoint(ctx, endpoint)
		if !ok {
			return -1, false
		}

		iaddr := (*os.InternetAddress)(unsafe.Pointer(&addrBuf))
		iaddr.Family = os.AddressFamilyInternet
		iaddr.Address = addr
		iaddr.Port = port

		addrLen = uint32(unsafe.Sizeof(*iaddr))
	default:
		panic("protocol is not supported")
	}

	s, ok := os.CreateNetworkSocket(ctx, pf, typ, prot)
	if !ok {
		return -1, false
	}

	paddr := (*os.NetworkAddress)(unsafe.Pointer(&addrBuf))
	if !os.BindSocketToAddress(ctx, s, paddr, addrLen) {
		os.CloseHandle(ctx, s)
		return -1, false
	}

	if pf == os.ProtocolFamilyInternet {
		if !os.SetSocketBooleanOption(ctx, s, os.SocketOptionReuseLocalAddressAndPortWithLoadBalancing, true) {
			os.CloseHandle(ctx, s)
			return -1, false
		}
		if typ == os.SocketTypeStream {
			if !os.SetSocketBooleanOption(ctx, s, os.SocketOptionTCPNoDelay, true) {
				os.CloseHandle(ctx, s)
				return -1, false
			}
		}
	}

	if !os.ListenForIncomingConnections(ctx, s, 128) {
		os.CloseHandle(ctx, s)
		return -1, false
	}

	return s, true
}
