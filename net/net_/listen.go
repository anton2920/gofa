package net_

import (
	"unsafe"

	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/os"
)

func Listen(proto string, endpoint string) (os.Handle, error) {
	var addrBuf [unsafe.Sizeof(os.NetworkAddress{})]byte
	var pf os.ProtocolFamily
	var typ os.SocketType
	var prot os.Protocol
	var addrLen uint32

	switch proto {
	case "tcp", "tcp/ip", "tcp/ip4", "tcp/ipv4":
		pf = os.ProtocolFamilyInternet
		typ = os.SocketTypeStream

		addr, port, err := tcp.ParseEndpoint(endpoint)
		if err != nil {
			return -1, err
		}

		iaddr := (*os.InternetAddress)(unsafe.Pointer(&addrBuf))
		iaddr.Family = os.AddressFamilyInternet
		iaddr.Address = addr
		iaddr.Port = port

		addrLen = uint32(unsafe.Sizeof(*iaddr))
	default:
		panic("protocol is not supported")
	}

	s, err := os.CreateNetworkSocket(pf, typ, prot)
	if err != nil {
		return -1, err
	}

	paddr := (*os.NetworkAddress)(unsafe.Pointer(&addrBuf))
	if err := os.BindSocketToAddress(s, paddr, addrLen); err != nil {
		os.CloseHandle(s)
		return -1, err
	}

	if pf == os.ProtocolFamilyInternet {
		if err := os.SetSocketBooleanOption(s, os.SocketOptionReuseLocalAddressAndPortWithLoadBalancing, true); err != nil {
			os.CloseHandle(s)
			return -1, err
		}
		if typ == os.SocketTypeStream {
			if err := os.SetSocketBooleanOption(s, os.SocketOptionTCPNoDelay, true); err != nil {
				os.CloseHandle(s)
				return -1, err
			}
		}
	}

	if err := os.ListenForIncomingConnections(s, 128); err != nil {
		os.CloseHandle(s)
		return -1, err
	}

	return s, nil
}
