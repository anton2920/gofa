package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/trace"
)

const DefaultMaxActiveClients = 16384

var CtxPool = NewContextPool(DefaultMaxActiveClients)

func Accept(l int32, bufferSize int) (*Context, error) {
	var addr syscall.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	rb, err := buffer.NewCircular(bufferSize)
	if err != nil {
		syscall.Close(c)
		return nil, fmt.Errorf("failed to create new request buffer: %w", err)
	}

	ctx, err := CtxPool.Get()
	if err != nil {
		ctx = new(Context)
		err = TooManyClients
	}
	InitContext(ctx, c, addr, rb)

	return ctx, err
}

func Read(ctx *Context) (int, error) {
	t := trace.Begin("")

	rBuf := ctx.RequestBuffer
	buf := rBuf.RemainingSlice()

	if len(buf) == 0 {
		trace.End(t)
		return 0, NoSpaceLeft
	}
	n, err := syscall.Read(ctx.Connection, buf)
	if err != nil {
		trace.End(t)
		return 0, err
	}
	rBuf.Produce(n)

	trace.End(t)
	return n, nil
}

//go:norace
func Write(ctx *Context) (int, error) {
	t := trace.Begin("")

	var written int
	if len(ctx.ResponseBuffer[ctx.ResponsePos:]) > 0 {
		n, err := syscall.Write(ctx.Connection, ctx.ResponseBuffer[ctx.ResponsePos:])
		if err != nil {
			trace.End(t)
			return 0, err
		}
		ctx.ResponsePos += n

		if ctx.ResponsePos == len(ctx.ResponseBuffer) {
			ctx.ResponseBuffer = ctx.ResponseBuffer[:0]
			ctx.ResponsePos = 0
			if ctx.CloseAfterWrite {
				Close(ctx)
			}
		}
	}

	trace.End(t)
	return written, nil
}

//go:norace
func Close(ctx *Context) error {
	ctx.ClientAddress = ""
	ctx.Check = 1 - ctx.Check
	ctx.CloseAfterWrite = false

	buffer.FreeCircular(ctx.RequestBuffer)
	ctx.RequestBuffer = nil
	ctx.ResponseBuffer = nil

	err := syscall.Close(ctx.Connection)

	CtxPool.Put(ctx)
	return err
}

//go:norace
func CloseAfterWrite(ctx *Context) {
	ctx.CloseAfterWrite = true
}
