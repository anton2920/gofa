//go:build gofanostd
// +build gofanostd

package log_

func Panic(_ interface{}) {
	panic(nil)
}
