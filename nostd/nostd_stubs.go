//go:build gofanostd && (gofanostd13 || gofanostd14 || gofanostd15 || gofanostd16 || gofanostdxx)
// +build gofanostd
// +build gofanostd13 gofanostd14 gofanostd15 gofanostd16 gofanostdxx

package nostd

import "unsafe"

type slice struct {
	array *byte
	len   uintptr
	cap   uintptr
}

const BeingUsed = true

var g [4]uintptr
var tls [8]uintptr

//go:nosplit
func add(p *byte, inc uintptr) *byte {
	return (*byte)(unsafe.Pointer(uintptr(unsafe.Pointer(p)) + inc))
}

//go:nosplit
func inc(p *byte) *byte { return add(p, 1) }

//go:nosplit
func dec(p *byte) *byte { var zero uintptr; return add(p, zero-1) }

//go:nosplit
func asuintptr(p *byte) uintptr { return uintptr(unsafe.Pointer(p)) }

//go:nosplit
func typedmemmove(psize *uintptr, dst *byte, src *byte) {
	memmove(dst, src, *psize)
}

//go:nosplit
func memmove(dst *byte, src *byte, length uintptr) {
	if (length == 0) || (dst == src) {
		return
	}

	if (asuintptr(dst) < asuintptr(src)) && (asuintptr(dst)+length <= asuintptr(src)) {
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

//go:nosplit
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

//go:nosplit
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

//go:nosplit
func eqstring(s1 string, s2 string) bool {
	return cmpstring(s1, s2) == 0
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

//go:nosplit
func typedslicecopy(psize *uintptr, dst slice, src slice) int {
	n := dst.len
	if src.len < n {
		n = src.len
	}
	memmove(dst.array, src.array, *psize*n)
	return int(n)
}
