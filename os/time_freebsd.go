//go:build freebsd
// +build freebsd

package os

import (
	"unsafe"

	"github.com/anton2920/gofa/os/posix/freebsd"
)

type SecondsWithNanoseconds struct {
	Seconds     int
	Nanoseconds int
}

func BlockForSpecifiedAmount(t SecondsWithNanoseconds) error {
	return freebsd.Nanosleep((*freebsd.Timespec)(unsafe.Pointer(&t)), nil)
}

func GetCurrentTime() (SecondsWithNanoseconds, error) {
	var t SecondsWithNanoseconds
	err := freebsd.ClockGettime(freebsd.CLOCK_REALTIME, (*freebsd.Timespec)(unsafe.Pointer(&t)))
	return t, err
}
