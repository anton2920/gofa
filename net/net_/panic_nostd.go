//go:build gofanostd
// +build gofanostd

package net_

func Panic(_ interface{}) {
	panic(nil)
}
