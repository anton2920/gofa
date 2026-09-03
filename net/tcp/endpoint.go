package tcp

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

func atoi(ctx *context.Context, s string) (int, bool) {
	var n, i int

	sign := 1
	if s[0] == '-' {
		sign = -1
		i++
	}

	for ; i < len(s); i++ {
		if (s[i] < '0') || (s[i] > '9') {
			ctx.NewError().S("expected digit, got ").Q(s[i : i+1])
			return 0, false
		}
		n = n*10 + int(s[i]-'0')
	}

	return sign * n, true
}

func ParseEndpoint(ctx *context.Context, endpoint string) (uint32, uint16, bool) {
	var addr uint32

	colon := strings.FindChar(endpoint, ':')
	if colon == -1 {
		ctx.NewError().S("no port specified")
		return 0, 0, false
	}

	part, ok := atoi(ctx, endpoint[colon+1:])
	if !ok {
		ctx.NewError().S("failed to parse port value: ").S(ctx.OldError())
		return 0, 0, false
	}
	port := ints.SwapBytesInWord(uint16(part))

	endpoint = endpoint[:colon]
	dot := strings.FindChar(endpoint, '.')
	if dot == -1 {
		return 0, port, true
	}
	part, ok = atoi(ctx, endpoint[:dot])
	if !ok {
		ctx.NewError().S("failed to parse first endpoint octet: ").S(ctx.OldError())
		return 0, 0, false
	}
	addr |= uint32(part)

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		ctx.NewError().S("expected second endpoint octet, found nothing")
		return 0, 0, false
	}
	part, ok = atoi(ctx, endpoint[:dot])
	if !ok {
		ctx.NewError().S("failed to parse second endpoint octet: ").S(ctx.OldError())
		return 0, 0, false
	}
	addr |= uint32(part) << 8

	endpoint = endpoint[dot+1:]
	dot = strings.FindChar(endpoint, '.')
	if dot == -1 {
		ctx.NewError().S("expected third endpoint octet, found nothing")
		return 0, 0, false
	}
	part, ok = atoi(ctx, endpoint[:dot])
	if !ok {
		ctx.NewError().S("failed to parse third endpoint octet: ").S(ctx.OldError())
		return 0, 0, false
	}
	addr |= uint32(part) << 16

	endpoint = endpoint[dot+1:]
	part, ok = atoi(ctx, endpoint)
	if !ok {
		ctx.NewError().S("failed to parse fourth endpoint octet: ").S(ctx.OldError())
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
