package tcp

import (
	"errors"
	"fmt"
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/syscall"
)

func SwapBytesInWord(x uint16) uint16 {
	return ((x << 8) & 0xFF00) | (x >> 8)
}

func ParseAddress(address string) (uint32, uint16, error) {
	var addr uint32

	colon := strings.FindChar(address, ':')
	if colon == -1 {
		return 0, 0, errors.New("no port specified")
	}

	part, err := strconv.Atoi(address[colon+1:])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse port value: %w", err)
	}
	port := SwapBytesInWord(uint16(part))

	address = address[:colon]
	dot := strings.FindChar(address, '.')
	if dot == -1 {
		return syscall.INADDR_ANY, port, nil
	}
	part, err = strconv.Atoi(address[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse first address octet: %w", err)
	}
	addr |= uint32(part)

	address = address[dot+1:]
	dot = strings.FindChar(address, '.')
	if dot == -1 {
		return 0, 0, fmt.Errorf("expected second address octet, found nothing")
	}
	part, err = strconv.Atoi(address[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse second address octet: %w", err)
	}
	addr |= uint32(part) << 8

	address = address[dot+1:]
	dot = strings.FindChar(address, '.')
	if dot == -1 {
		return 0, 0, fmt.Errorf("expected third address octet, found nothing")
	}
	part, err = strconv.Atoi(address[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse third address octet: %w", err)
	}
	addr |= uint32(part) << 16

	address = address[dot+1:]
	part, err = strconv.Atoi(address)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse fourth address octet: %w", err)
	}
	addr |= uint32(part) << 24

	return addr, port, nil
}

func PutAddress(buffer []byte, addr uint32, port uint16) int {
	var n int

	n += slices.PutInt(buffer[n:], int((addr&0x000000FF)>>0))
	buffer[n] = ':'
	n++

	n += slices.PutInt(buffer[n:], int((addr&0x0000FF00)>>8))
	buffer[n] = '.'
	n++

	n += slices.PutInt(buffer[n:], int((addr&0x00FF0000)>>16))
	buffer[n] = '.'
	n++

	n += slices.PutInt(buffer[n:], int((addr&0xFF000000)>>24))
	buffer[n] = '.'
	n++

	n += slices.PutInt(buffer[n:], int(SwapBytesInWord(port)))

	return n
}

/* Listen creates TCP/IPv4 socket and starts listening on a specified address. */
func Listen(address string, backlog int) (os.Handle, error) {
	addr, port, err := ParseAddress(address)
	if err != nil {
		return -1, fmt.Errorf("failed to parse address string: %w", err)
	}

	l, err := syscall.Socket(syscall.PF_INET, syscall.SOCK_STREAM, 0)
	if err != nil {
		return -1, fmt.Errorf("failed to create new socket: %w", err)
	}

	var enable int32 = 1
	if err := syscall.Setsockopt(l, syscall.SOL_SOCKET, syscall.SO_REUSEPORT_LB, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable))); err != nil {
		return -1, fmt.Errorf("failed to apply options to socket: %w", err)
	}

	if err := syscall.Setsockopt(l, syscall.IPPROTO_TCP, syscall.TCP_NODELAY, unsafe.Pointer(&enable), uint32(unsafe.Sizeof(enable))); err != nil {
		return -1, fmt.Errorf("failed to apply options to socket: %w", err)
	}

	sin := syscall.SockAddrIn{Family: syscall.AF_INET, Addr: addr, Port: port}
	if err := syscall.Bind(l, (*syscall.Sockaddr)(unsafe.Pointer(&sin)), uint32(unsafe.Sizeof(sin))); err != nil {
		return -1, fmt.Errorf("failed to bind socket to address: %w", err)
	}

	if err := syscall.Listen(l, int32(backlog)); err != nil {
		return -1, fmt.Errorf("failed to listen for incoming connections: %w", err)
	}

	return os.Handle(l), nil
}
