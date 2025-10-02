//go:build linux || freebsd || openbsd || netbsd || darwin
// +build linux freebsd openbsd netbsd darwin

package os

import (
	"github.com/anton2920/gofa/os/unix"
	"github.com/anton2920/gofa/trace"
)

func Close(handle Handle) error {
	return unix.Close(int32(handle))
}

func Read(handle Handle, buf []byte) (int64, error) {
	t := trace.Begin("")

	n, err := unix.Read(int32(handle), buf)

	trace.End(t)
	return n, err
}

func Write(handle Handle, buf []byte) (int64, error) {
	t := trace.Begin("")

	n, err := unix.Write(int32(handle), buf)

	trace.End(t)
	return n, err
}
