package cgo

import "unsafe"

type On uintptr

//go:nosplit
//go:noescape
func Call0(fn unsafe.Pointer) uintptr

//go:nosplit
//go:noescape
func Call1(fn unsafe.Pointer, a0 uintptr) uintptr

//go:nosplit
//go:noescape
func Call2(fn unsafe.Pointer, a0 uintptr, a1 uintptr) uintptr

//go:nosplit
//go:noescape
func Call3(fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr) uintptr

//go:nosplit
//go:noescape
func Call4(fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr) uintptr

//go:nosplit
//go:noescape
func on_call0(sp On, fn unsafe.Pointer) uintptr

//go:nosplit
//go:noescape
func on_call1(sp On, fn unsafe.Pointer, a0 uintptr) uintptr

//go:nosplit
//go:noescape
func on_call2(sp On, fn unsafe.Pointer, a0 uintptr, a1 uintptr) uintptr

//go:nosplit
//go:noescape
func on_call3(sp On, fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr) uintptr

//go:nosplit
//go:noescape
func on_call4(sp On, fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr) uintptr

func (sp On) Call0(fn unsafe.Pointer) uintptr {
	return on_call0(sp, fn)
}

func (sp On) Call1(fn unsafe.Pointer, a0 uintptr) uintptr {
	return on_call1(sp, fn, a0)
}

func (sp On) Call2(fn unsafe.Pointer, a0 uintptr, a1 uintptr) uintptr {
	return on_call2(sp, fn, a0, a1)
}

func (sp On) Call3(fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr) uintptr {
	return on_call3(sp, fn, a0, a1, a2)
}

func (sp On) Call4(fn unsafe.Pointer, a0 uintptr, a1 uintptr, a2 uintptr, a3 uintptr) uintptr {
	return on_call4(sp, fn, a0, a1, a2, a3)
}
