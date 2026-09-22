//go:build !gofanostd
// +build !gofanostd

package net_

func Panic(msg interface{}) {
	panic(msg)
}
