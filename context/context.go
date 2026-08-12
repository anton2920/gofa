package context

import (
	"github.com/anton2920/gofa/fmt"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/mem"
)

type Context struct {
	Arena mem.Arena

	Fmt fmt.Formatter
	Log log.Formatter

	ErrFmt  [2]fmt.Formatter
	ErrCode [2]int
	ErrCurr int
}

func (ctx *Context) InitWithMinimumLogLevelAndFourByteSlices(level log.Level, arenaBuf []byte, fmtBuf []byte, logBuf []byte, errBuf []byte) {
	ctx.Arena.InitWithByteSlice(arenaBuf)

	ctx.Fmt.InitWithByteSlice(fmtBuf)
	ctx.Log.InitWithMinimumLevelAndByteSlice(level, logBuf)

	size := len(errBuf) / len(ctx.ErrFmt)
	for i := 0; i < len(ctx.ErrFmt); i++ {
		ctx.ErrFmt[i].InitWithByteSlice(errBuf[i*size : (i+1)*size])
	}
}

func (ctx *Context) InitWithMinimumLogLevelAndEvenlySplitByteSlice(level log.Level, buf []byte) {
	size := len(buf) / 4
	ctx.InitWithMinimumLogLevelAndFourByteSlices(level, buf[0:size], buf[size:2*size], buf[2*size:3*size], buf[3*size:])
}

func (ctx *Context) Init(level log.Level, buf []byte) {
	ctx.InitWithMinimumLogLevelAndEvenlySplitByteSlice(level, buf)
}

func (ctx *Context) Error() string {
	return ctx.ErrFmt[ctx.ErrCurr].String()
}

func (ctx *Context) ErrorCode() int {
	return ctx.ErrCode[ctx.ErrCurr]
}

func (ctx *Context) NewErrorWithCode(code int) *fmt.Formatter {
	ctx.ErrCurr = 1 - ctx.ErrCurr
	ctx.ErrCode[ctx.ErrCurr] = code
	return &ctx.ErrFmt[ctx.ErrCurr].Reset()
}

func (ctx *Context) NewError() *fmt.Formatter {
	return ctx.NewErrorWithCode(0)
}
