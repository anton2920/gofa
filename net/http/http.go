package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/syscall"
)

func Accept(l int32, bufferSize int) (*Context, error) {
	var addr tcp.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	ctx, err := NewContext(c, addr, bufferSize)
	if err != nil {
		syscall.Close(c)
		return nil, fmt.Errorf("failed to create new  context: %w", err)
	}

	return ctx, nil
}

//go:norace
func AddClientToQueue(q *event.Queue, ctx *Context, request event.Request, trigger event.Trigger) error {
	/* TODO(anton2920): switch to pinning inside platform methods. */
	q.Pinner.Pin(ctx)
	return q.AddSocket(ctx.Connection, request, trigger, ctx.Pointer())
}

func ContentTypeHTML(bodies []syscall.Iovec) bool {
	return (len(bodies) > 0) && (bodies[0] == html.Header)
}

func Read(ctx *Context) (int, error) {
	rBuf := &ctx.RequestBuffer
	buf := rBuf.RemainingSlice()

	if len(buf) == 0 {
		return 0, NoSpaceLeft
	}
	n, err := syscall.Read(ctx.Connection, buf)
	if err != nil {
		return 0, err
	}
	rBuf.Produce(int(n))

	return n, nil
}

func Write(ctx *Context) (int, error) {
	var written int

	/* NOTE(anton2920): IOV_MAX is 1024, so F**CK ME for not sending large pipelines with one syscall!!! */
	for len(ctx.ResponseIovs[ctx.ResponsePos:]) > 0 {
		end := min(len(ctx.ResponseIovs[ctx.ResponsePos:]), syscall.IOV_MAX)
		n, err := syscall.Writev(ctx.Connection, ctx.ResponseIovs[ctx.ResponsePos:ctx.ResponsePos+end])
		if err != nil {
			return 0, err
		}
		written += int(n)

		prevPos := ctx.ResponsePos
		for (ctx.ResponsePos < len(ctx.ResponseIovs)) && (n >= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))) {
			n -= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))
			ctx.ResponsePos++
		}
		if ctx.ResponsePos == len(ctx.ResponseIovs) {
			ctx.ResponseIovs = ctx.ResponseIovs[:0]
			ctx.ResponsePos = 0
		} else if ctx.ResponsePos-prevPos < end {
			log.Panicf("Written %d iovs out of %d", ctx.ResponsePos-prevPos, end)
			ctx.ResponseIovs[ctx.ResponsePos] = ctx.ResponseIovs[ctx.ResponsePos][n:]
			break

			/* TODO(anton2920): as an option, gather buffers manually into some local arena. */
		}
	}

	if ctx.CloseAfterWrite {
		Close(ctx)
	}
	return written, nil
}

func Close(ctx *Context) error {
	ctx.Reset()
	buffer.FreeCircular(&ctx.RequestBuffer)
	return syscall.Close(ctx.Connection)
}

func CloseAfterWrite(ctx *Context) {
	ctx.CloseAfterWrite = true
}
