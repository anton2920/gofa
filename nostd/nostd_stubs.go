//go:build gofanostd
// +build gofanostd

package nostd

import "unsafe"

var g [4]uintptr
var tls [8]uintptr

func add(p *byte, inc uintptr) *byte {
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + inc))
}
func inc(p *byte) *byte         { return add(p, 1) }
func dec(p *byte) *byte         { var zero uintptr; return add(p, zero-1) }
func asuintptr(p *byte) uintptr { return uintptr(unsafe.Pointer(p)) }

func memmove(dst *byte, src *byte, length uintptr) {
	if (length == 0) || (dst == src) {
		return
	}

	if asuintptr(dst) < asuintptr(src) {
		for length > 0 {
			*dst = *src
			dst = inc(dst)
			src = inc(src)
			length--
		}
	} else {
		dst = add(dst, length-1)
		src = add(src, length-1)
		for length > 0 {
			*dst = *src
			dst = dec(dst)
			src = dec(src)
			length--
		}
	}
}

func cmpstring(s1, s2 string) int {
	l := len(s1)
	if len(s2) < l {
		l = len(s2)
	}
	for i := 0; i < l; i++ {
		c1, c2 := s1[i], s2[i]
		if c1 < c2 {
			return -1
		}
		if c1 > c2 {
			return +1
		}
	}
	if len(s1) < len(s2) {
		return -1
	}
	if len(s1) > len(s2) {
		return +1
	}
	return 0
}

func cmpbytes(s1, s2 []byte) int {
	l := len(s1)
	if len(s2) < l {
		l = len(s2)
	}
	for i := 0; i < l; i++ {
		c1, c2 := s1[i], s2[i]
		if c1 < c2 {
			return -1
		}
		if c1 > c2 {
			return +1
		}
	}
	if len(s1) < len(s2) {
		return -1
	}
	if len(s1) > len(s2) {
		return +1
	}
	return 0
}

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

func strhash(a *byte, s, h uintptr) uintptr {
	return memhash(*(**byte)(unsafe.Pointer(a)), uintptr(len(*(*string)(unsafe.Pointer(a)))), h)
}

func eqstring(s1 string, s2 string) bool {
	return cmpstring(s1, s2) == 0
}

func memequal(p *byte, q *byte, size uintptr) bool {
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

func memequal0(p, q *byte, size uintptr) bool {
	return true
}
func memequal8(p, q *byte, size uintptr) bool {
	return *(*int8)(unsafe.Pointer(p)) == *(*int8)(unsafe.Pointer(q))
}
func memequal16(p, q *byte, size uintptr) bool {
	return *(*int16)(unsafe.Pointer(p)) == *(*int16)(unsafe.Pointer(q))
}
func memequal32(p, q *byte, size uintptr) bool {
	return *(*int32)(unsafe.Pointer(p)) == *(*int32)(unsafe.Pointer(q))
}
func memequal64(p, q *byte, size uintptr) bool {
	return *(*int64)(unsafe.Pointer(p)) == *(*int64)(unsafe.Pointer(q))
}
func memequal128(p, q *byte, size uintptr) bool {
	return *(*[2]int64)(unsafe.Pointer(p)) == *(*[2]int64)(unsafe.Pointer(q))
}
func f32equal(p, q *byte, size uintptr) bool {
	return *(*float32)(unsafe.Pointer(p)) == *(*float32)(unsafe.Pointer(q))
}
func f64equal(p, q *byte, size uintptr) bool {
	return *(*float64)(unsafe.Pointer(p)) == *(*float64)(unsafe.Pointer(q))
}
func c64equal(p, q *byte, size uintptr) bool {
	return *(*complex64)(unsafe.Pointer(p)) == *(*complex64)(unsafe.Pointer(q))
}
func c128equal(p, q *byte, size uintptr) bool {
	return *(*complex128)(unsafe.Pointer(p)) == *(*complex128)(unsafe.Pointer(q))
}
func strequal(p, q *byte, size uintptr) bool {
	return *(*string)(unsafe.Pointer(p)) == *(*string)(unsafe.Pointer(q))
}

// NOTE: Really dst *unsafe.Pointer, src unsafe.Pointer,
// but if we do that, Go inserts a write barrier on *dst = src.
//
//go:nosplit
func writebarrierptr(dst *uintptr, src uintptr) {
	*dst = src
}

//go:nosplit
func writebarrierstring(dst *[2]uintptr, src [2]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
}

//go:nosplit
func writebarrierslice(dst *[3]uintptr, src [3]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
	dst[2] = src[2]
}

//go:nosplit
func writebarrieriface(dst *[2]uintptr, src [2]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
}

//go:nosplit
func writebarrierfat2(dst *[2]uintptr, _ *byte, src [2]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
}

//go:nosplit
func writebarrierfat3(dst *[3]uintptr, _ *byte, src [3]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
	dst[2] = src[2]
}

//go:nosplit
func writebarrierfat4(dst *[4]uintptr, _ *byte, src [4]uintptr) {
	dst[0] = src[0]
	dst[1] = src[1]
	dst[2] = src[2]
	dst[3] = src[3]
}

//go:nosplit
func writebarrierfat(siz *uintptr, dst *byte, src *byte) {
	memmove(dst, src, *siz)
}
