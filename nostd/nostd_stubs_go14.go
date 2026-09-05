//go:build gofanostd && gofanostd14
// +build gofanostd,gofanostd14

package nostd

import "unsafe"

//go:nosplit
func memhash(p *byte, s, h uintptr) uintptr {
	const (
		ptrSize = unsafe.Sizeof(uintptr(0))
		c0      = uintptr((8-ptrSize)/4*2860486313 + (ptrSize-4)/4*33054211828000289)
		c1      = uintptr((8-ptrSize)/4*3267000013 + (ptrSize-4)/4*23344194077549503)
	)

	h ^= c0
	for s > 0 {
		h = (h ^ uintptr(*p)) * c1
		p = inc(p)
		s--
	}
	return h
}

//go:nosplit
func memequal(p *byte, q *byte, size uintptr) bool {
	if p == q {
		return true
	}
	for size > 0 {
		if *p != *q {
			return false
		}
		p = inc(p)
		q = inc(q)
		size--
	}
	return true
}

//go:nosplit
func memequal0(p, q *byte, size uintptr) bool {
	return true
}

//go:nosplit
func memequal8(p, q *byte, size uintptr) bool {
	return *(*int8)(unsafe.Pointer(p)) == *(*int8)(unsafe.Pointer(q))
}

//go:nosplit
func memequal16(p, q *byte, size uintptr) bool {
	return *(*int16)(unsafe.Pointer(p)) == *(*int16)(unsafe.Pointer(q))
}

//go:nosplit
func memequal32(p, q *byte, size uintptr) bool {
	return *(*int32)(unsafe.Pointer(p)) == *(*int32)(unsafe.Pointer(q))
}

//go:nosplit
func memequal64(p, q *byte, size uintptr) bool {
	return *(*int64)(unsafe.Pointer(p)) == *(*int64)(unsafe.Pointer(q))
}

//go:nosplit
func memequal128(p, q *byte, size uintptr) bool {
	return *(*[2]int64)(unsafe.Pointer(p)) == *(*[2]int64)(unsafe.Pointer(q))
}

//go:nosplit
func f32equal(p, q *byte, size uintptr) bool {
	return *(*float32)(unsafe.Pointer(p)) == *(*float32)(unsafe.Pointer(q))
}

//go:nosplit
func f64equal(p, q *byte, size uintptr) bool {
	return *(*float64)(unsafe.Pointer(p)) == *(*float64)(unsafe.Pointer(q))
}

//go:nosplit
func c64equal(p, q *byte, size uintptr) bool {
	return *(*complex64)(unsafe.Pointer(p)) == *(*complex64)(unsafe.Pointer(q))
}

//go:nosplit
func c128equal(p, q *byte, size uintptr) bool {
	return *(*complex128)(unsafe.Pointer(p)) == *(*complex128)(unsafe.Pointer(q))
}

//go:nosplit
func strequal(p, q *byte, size uintptr) bool {
	return *(*string)(unsafe.Pointer(p)) == *(*string)(unsafe.Pointer(q))
}

//go:nosplit
func strhash(a *byte, _ uintptr, h uintptr) uintptr {
	return memhash(*(**byte)(unsafe.Pointer(a)), uintptr(len(*(*string)(unsafe.Pointer(a)))), h)
}
