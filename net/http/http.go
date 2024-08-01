package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/syscall"
)

func Accept(l int32, ctx *Context, bufferSize int) error {
	var addr syscall.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	rb, err := buffer.NewCircular(bufferSize)
	if err != nil {
		syscall.Close(c)
		return fmt.Errorf("failed to create new request buffer: %w", err)
	}

	InitContext(ctx, c, addr, rb)
	return nil
}

func Read(ctx *Context) (int, error) {
	p := prof.Begin("")

	rBuf := ctx.RequestBuffer
	buf := rBuf.RemainingSlice()

	if len(buf) == 0 {
		prof.End(p)
		return 0, NoSpaceLeft
	}
	n, err := syscall.Read(ctx.Connection, buf)
	if err != nil {
		prof.End(p)
		return 0, err
	}
	rBuf.Produce(int(n))

	prof.End(p)
	return n, nil
}

func Write(ctx *Context) (int, error) {
	p := prof.Begin("")

	var written int
	if len(ctx.ResponseBuffer[ctx.ResponsePos:]) > 0 {
		n, err := syscall.Write(ctx.Connection, ctx.ResponseBuffer[ctx.ResponsePos:])
		if err != nil {
			prof.End(p)
			return 0, err
		}
		ctx.ResponsePos += int(n)

		if ctx.ResponsePos == len(ctx.ResponseBuffer) {
			ctx.ResponseBuffer = ctx.ResponseBuffer[:0]
			ctx.ResponsePos = 0
			if ctx.CloseAfterWrite {
				Close(ctx)
			}
		}
	}

	prof.End(p)
	return written, nil
}

func Close(ctx *Context) error {
	ctx.Check = 1 - ctx.Check
	ctx.CloseAfterWrite = false
	buffer.FreeCircular(ctx.RequestBuffer)
	return syscall.Close(ctx.Connection)
}

func CloseAfterWrite(ctx *Context) {
	ctx.CloseAfterWrite = true
}
