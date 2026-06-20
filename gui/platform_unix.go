//go:build aix || darwin || dragonfly || freebsd || illumos || linux || netbsd || openbsd || solaris || zos
// +build aix darwin dragonfly freebsd illumos linux netbsd openbsd solaris zos

package gui

import (
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

func platformSleep(duration int64) error {
	tp := syscall.Timespec{Sec: duration / time.Second, Nsec: duration % time.Second}
	return syscall.Nanosleep(&tp, nil)
}
