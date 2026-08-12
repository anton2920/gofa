//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os/posix/freebsd"
)

type SecondsWithNanoseconds struct {
	Seconds     int
	Nanoseconds int
}

func BlockForSpecifiedAmountOfTime(ctx *context.Context, t SecondsWithNanoseconds) bool {
	return freebsd.Nanosleep(ctx, (*freebsd.Timespec)(unsafe.Pointer(&t)), nil)
}

func GetCurrentTime(ctx *context.Context, t *SecondsWithNanoseconds) bool {
	return freebsd.ClockGettime(ctx, freebsd.CLOCK_REALTIME, (*freebsd.Timespec)(unsafe.Pointer(t)))
}
