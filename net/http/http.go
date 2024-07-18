package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/event"
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

	n, err := syscall.Write(ctx.Connection, ctx.ResponseBuffer[ctx.ResponsePos:])
	if err != nil {
		return 0, err
	}
	ctx.ResponsePos += int(n)

	if ctx.ResponsePos == len(ctx.ResponseBuffer) {
		ctx.ResponseBuffer = ctx.ResponseBuffer[:0]
		ctx.ResponsePos = 0
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
