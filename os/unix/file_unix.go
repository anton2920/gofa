//go:build linux || freebsd || openbsd || netbsd || darwin
// +build linux freebsd openbsd netbsd darwin

package unix

import "github.com/anton2920/gofa/syscall"

/* TODO(anton2920): unify 'syscall' and 'unix' packages. */
func Close(f int32) error {
	return syscall.Close(f)
}

func Read(f int32, buf []byte) (int, error) {
	return syscall.Read(f, buf)
}

func Write(f int32, buf []byte) (int, error) {
	return syscall.Write(f, buf)
}
