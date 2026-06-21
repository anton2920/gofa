package cgo

import "unsafe"

//go:nosplit
//go:noescape
func Call0(fn unsafe.Pointer) uintptr

//go:nosplit
//go:noescape
func Call1(fn unsafe.Pointer, a0 uintptr) uintptr

//go:nosplit
//go:noescape
func Call2(fn unsafe.Pointer, a0 uintptr, a1 uintptr) uintptr
