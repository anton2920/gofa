package tcp

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

/* TODO(anton2920): remove memory allocation on error. */
func ParseEndpoint(endpoint string) (uint32, uint16, error) {
	var addr uint32

	colon := strings.FindChar(endpoint, ':')
	if colon == -1 {
		return 0, 0, errors.New("no port specified")
	}

	part, err := strconv.Atoi(endpoint[colon+1:])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse port value: %v", err)
	}
	port := ints.SwapBytesInWord(uint16(part))

	endpoint = endpoint[:colon]
	dot := strings.FindChar(endpoint, '.')
	if dot == -1 {
		return 0, port, nil
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse first endpoint octet: %v", err)
	}
	addr |= uint32(part)

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		return 0, 0, fmt.Errorf("expected second endpoint octet, found nothing")
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse second endpoint octet: %v", err)
	}
	addr |= uint32(part) << 8

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		return 0, 0, fmt.Errorf("expected third endpoint octet, found nothing")
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse third endpoint octet: %v", err)
	}
	addr |= uint32(part) << 16

	endpoint = endpoint[dot+1:]
	part, err = strconv.Atoi(endpoint)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse fourth endpoint octet: %v", err)
	}
	addr |= uint32(part) << 24

	return addr, port, nil
}

func PutEndpoint(buffer []byte, addr uint32, port uint16) int {
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

	n += slices.PutInt(buffer[n:], int(ints.SwapBytesInWord(port)))

	return n
}
