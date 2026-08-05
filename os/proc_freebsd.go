//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/os/posix/freebsd"
)

type (
	Signal        freebsd.Signal
	SignalHandler unsafe.Pointer
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

func InstallSignalHandler(s Signal, handler SignalHandler) error {
	act := freebsd.Sigaction_t{Handler: uintptr(handler), Flags: freebsd.SA_ONSTACK | freebsd.SA_RESTART}
	return freebsd.Sigaction(int32(s), &act, nil)
}
