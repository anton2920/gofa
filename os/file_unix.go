//go:build linux || freebsd || openbsd || netbsd || darwin
// +build linux freebsd openbsd netbsd darwin

package os

import "github.com/anton2920/gofa/os/unix"

func Close(handle Handle) error {
	return unix.Close(int32(handle))
}

func Read(handle Handle, buf []byte) (int64, error) {
	return unix.Read(int32(handle), buf)
}

func Write(handle Handle, buf []byte) (int64, error) {
	return unix.Write(int32(handle), buf)
}
