package http

import (
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/syscall"
)

type Context struct {
	/* NOTE(anton2920): Check must be the same as the last pointer's bit, if context is in use. */
	Check int32

	Connection    int32
	ClientAddress string

	RequestBuffer *buffer.Circular

	ResponseBuffer []byte
	ResponsePos    int

	/* TODO(anton2920): I don't like this. */
	CloseAfterWrite bool
}

//go:norace
func InitContext(ctx *Context, c int32, addr syscall.SockAddrIn, rb *buffer.Circular) {
	ctx.Connection = c
	ctx.RequestBuffer = rb
	ctx.ResponseBuffer = make([]byte, 0, len(rb.Buf))

	buffer := make([]byte, 21)
	n := tcp.PutAddress(buffer, addr.Addr, addr.Port)
	ctx.ClientAddress = string(buffer[:n])
}

func GetContextFromPointer(ptr unsafe.Pointer) (*Context, bool) {
	if ptr == nil {
		return nil, false
	}

	check := uintptr(ptr) & 0x1
	ctx := (*Context)(unsafe.Pointer(uintptr(ptr) - check))

	return ctx, ctx.Check == int32(check)
}

//go:norace
func (ctx *Context) Pointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(ctx)) | uintptr(ctx.Check))
}
