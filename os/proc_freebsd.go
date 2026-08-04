//go:build freebsd
// +build freebsd

package os

import "github.com/anton2920/gofa/os/posix/freebsd"

//go:nosplit
func Exit(code int) {
	freebsd.Exit(int32(code))
}
