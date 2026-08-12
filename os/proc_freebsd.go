//go:build freebsd
// +build freebsd

package os

import (
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type (
	Signal        freebsd.Signal
	SignalHandler uintptr
)

const (
	SignalHangup    = Signal(freebsd.SIGHUP)
	SignalInterrupt = Signal(freebsd.SIGINT)
	SignalTerminate = Signal(freebsd.SIGTERM)
)

var (
	DefaultSignalHandler = SignalHandler(freebsd.SIG_DFL)
	IgnoreSignalHandler  = SignalHandler(freebsd.SIG_IGN)
)

func (s Signal) String() string {
	return freebsd.Signal(s).String()
}

//go:nosplit
func Exit(code int) {
	freebsd.Exit(int32(code))
}

/* NOTE(anton2920): this is f**cking unsafe as hell! You may manage to get it working, but you better have 'os.Exit' at the end of your handler or pray that your program is not inside a 'Syscall[69]'. Also, may gods have mercy on your soul... */
//go:nosplit
func AsSignalHandler(fn func(Signal)) SignalHandler

func InstallSignalHandler(ctx *context.Context, s Signal, handler SignalHandler) bool {
	act := freebsd.Sigaction_t{Handler: uintptr(handler)}
	return freebsd.Sigaction(ctx, int32(s), &act, nil)
}
