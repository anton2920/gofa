//go:build unix

package gui

import (
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

func platformSleep(duration int64) error {
	tp := syscall.Timespec{Sec: duration / time.NsecPerSec, Nsec: duration % time.NsecPerSec}
	return syscall.Nanosleep(&tp, nil)
}
