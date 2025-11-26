//go:build linux || freebsd || openbsd || netbsd || darwin
// +build linux freebsd openbsd netbsd darwin

package os

import (
	"github.com/anton2920/gofa/os/unix"
	"github.com/anton2920/gofa/syscall"
)

/* TODO(anton2920): rewrite using 'unix' package. */
func Open(path string) (Handle, error) {
	f, err := syscall.Open(path, syscall.O_RDONLY, 0)
	return Handle(f), err
}

func Close(handle Handle) error {
	return unix.Close(int32(handle))
}

func Read(handle Handle, buf []byte) (int64, error) {
	return unix.Read(int32(handle), buf)
}

func Write(handle Handle, buf []byte) (int64, error) {
	return unix.Write(int32(handle), buf)
}
