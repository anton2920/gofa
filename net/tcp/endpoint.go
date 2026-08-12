package tcp

import (
	"strconv"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

func ParseEndpoint(ctx *context.Context, endpoint string) (uint32, uint16, bool) {
	var addr uint32

	colon := strings.FindChar(endpoint, ':')
	if colon == -1 {
		ctx.NewError().S("no port specified")
		return 0, 0, false
	}

	/* TODO(anton2920): replace this with my own 'strings.ToInt'. */
	part, err := strconv.Atoi(endpoint[colon+1:])
	if err != nil {
		ctx.NewError().S("failed to parse port value: ").Err(err)
		return 0, 0, false
	}
	port := ints.SwapBytesInWord(uint16(part))

	endpoint = endpoint[:colon]
	dot := strings.FindChar(endpoint, '.')
	if dot == -1 {
		return 0, port, true
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		ctx.NewError().S("failed to parse first endpoint octet: ").Err(err)
		return 0, 0, false
	}
	addr |= uint32(part)

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		ctx.NewError().S("expected second endpoint octet, found nothing")
		return 0, 0, false
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		ctx.NewError().S("failed to parse second endpoint octet: ").Err(err)
		return 0, 0, false
	}
	addr |= uint32(part) << 8

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		ctx.NewError().S("expected third endpoint octet, found nothing")
		return 0, 0, false
	}
	part, err = strconv.Atoi(endpoint[:dot])
	if err != nil {
		ctx.NewError().S("failed to parse third endpoint octet: ").Err(err)
		return 0, 0, false
	}
	addr |= uint32(part) << 16

	endpoint = endpoint[dot+1:]
	part, err = strconv.Atoi(endpoint)
	if err != nil {
		ctx.NewError().S("failed to parse fourth endpoint octet: ").Err(err)
		return 0, 0, false
	}
	addr |= uint32(part) << 24

	return addr, port, true
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
