package context

import (
	"unsafe"

	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/mem"
	"github.com/anton2920/gofa/pointers"
)

type Context struct {
	Arena mem.Arena

	Fmt fmt.Formatter
	Log log.Formatter

	ErrFmt  [2]fmt.Formatter
	ErrCode [2]int
	ErrCurr int
}

func (ctx *Context) InitWithFourByteSlices(arenaBuf []byte, fmtBuf []byte, logBuf []byte, errBuf []byte) {
	ctx.Arena.InitWithByteSlice(arenaBuf)

	ctx.Fmt.InitWithByteSlice(fmtBuf)
	ctx.Log.InitWithByteSlice(logBuf)

	size := len(errBuf) / len(ctx.ErrFmt)
	for i := 0; i < len(ctx.ErrFmt); i++ {
		ctx.ErrFmt[i].InitWithByteSlice(errBuf[i*size : (i+1)*size])
	}
}

func (ctx *Context) InitWithEvenlySplitByteSlice(buf []byte) {
	size := len(buf) / 4
	ctx.InitWithFourByteSlices(buf[0:size], buf[size:2*size], buf[2*size:3*size], buf[3*size:])
}

func (ctx *Context) Error() string {
	return ctx.ErrFmt[ctx.ErrCurr].String()
}

func (ctx *Context) ErrorCode() int {
	return ctx.ErrCode[ctx.ErrCurr]
}

func (ctx *Context) OldError() string {
	return ctx.ErrFmt[1-ctx.ErrCurr].String()
}

func (ctx *Context) OldErrorCode() int {
	return ctx.ErrCode[1-ctx.ErrCurr]
}

func (ctx *Context) NewErrorWithCode(code int) *fmt.Formatter {
	ctx.ErrCurr = 1 - ctx.ErrCurr
	ctx.ErrCode[ctx.ErrCurr] = code
	return ctx.ErrFmt[ctx.ErrCurr].Reset()
}

func (ctx *Context) NewError() *fmt.Formatter {
	return ctx.NewErrorWithCode(0)
}

func (ctx *Context) ResetError() {
	ctx.ErrFmt[ctx.ErrCurr].Reset()
	ctx.ErrCode[ctx.ErrCurr] = 0
}

func (ctx *Context) OK() bool {
	return len(ctx.Error()) == 0
}

func (ctx *Context) Noescape() *Context {
	return (*Context)(pointers.UnsafeNoescape(unsafe.Pointer(ctx)))
}
