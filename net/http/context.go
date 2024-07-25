package http

import (
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/syscall"
)

type Context struct {
	/* Check must be the same as the last pointer's bit, if context is in use. */
	Check int32

	Connection    int32
	ClientAddress string

	RequestBuffer *buffer.Circular

	ResponseBuffer []byte
	ResponsePos    int

	/* TODO(anton2920): I don't like this. */
	CloseAfterWrite bool
}

func NewContext(c int32, addr syscall.SockAddrIn, bufferSize int) (*Context, error) {
	rb, err := buffer.NewCircular(bufferSize)
	if err != nil {
		return nil, err
	}

	ctx := new(Context)
	ctx.Connection = c
	ctx.RequestBuffer = rb
	ctx.ResponseBuffer = make([]byte, 0, 1024*1024)

	buffer := make([]byte, 21)
	n := tcp.PutAddress(buffer, addr.Addr, addr.Port)
	ctx.ClientAddress = string(buffer[:n])

	return ctx, nil
}

/* TODO(anton2920): maybe also test for some magic? */
func GetContextFromPointer(ptr unsafe.Pointer) (*Context, bool) {
	if ptr == nil {
		return nil, false
	}
	uptr := uintptr(ptr)

	check := uptr & 0x1
	ctx := (*Context)(unsafe.Pointer(uptr - check))

	return ctx, ctx.Check == int32(check)
}

func (ctx *Context) Pointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(ctx)) | uintptr(ctx.Check))
}

/* TODO(anton2920): maybe remove this altogether? */
func (ctx *Context) Reset() {
	ctx.Check = 1 - ctx.Check
	ctx.CloseAfterWrite = false
	ctx.RequestBuffer.Reset()
	ctx.ResponsePos = 0
	ctx.ResponseBuffer = ctx.ResponseBuffer[:0]
}
