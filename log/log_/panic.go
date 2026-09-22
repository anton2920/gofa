//go:build !gofanostd
// +build !gofanostd

package log_

func Panic(msg interface{}) {
	panic(msg)
}
