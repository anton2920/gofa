//go:build gofanostd
// +build gofanostd

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

var framepointer_enabled bool

//go:linkname buildVersion runtime.buildVersion
var buildVersion string

//go:linkname modinfo runtime.modinfo
var modinfo string

//go:nosplit
func exit(code int32)

//go:linkname main_main main.main
func main_main()

//go:nosplit
func main() {
	if buildVersion == "" {
		buildVersion = "unknown"
	}
	if len(modinfo) == 1 {
		modinfo = ""
	}
	main_main()
}

//go:nosplit
//go:linkname panicslice runtime.panicslice
func panicslice() {
	exit(66)
}

//go:nosplit
//go:linkname panicSliceAlen runtime.panicSliceAlen
func panicSliceAlen() {
	exit(66)
}

//go:nosplit
//go:linkname panicSliceAcap runtime.panicSliceAcap
func panicSliceAcap() {
	exit(66)
}

//go:nosplit
//go:linkname panicSliceB runtime.panicSliceB
func panicSliceB() {
	exit(66)
}

//go:nosplit
//go:linkname panicindex runtime.panicindex
func panicindex() {
	exit(67)
}

//go:nosplit
//go:linkname panicIndex runtime.panicIndex
func panicIndex() {
	exit(67)
}

//go:nosplit
//go:linkname panicwrap runtime.panicwrap
func panicwrap() {
	exit(68)
}

//go:nosplit
//go:linkname gopanic runtime.gopanic
func gopanic() {
	exit(69)
}

//go:nosplit
//go:linkname panic runtime.panic
func panic() {
	exit(69)
}

//go:nosplit
//go:linkname growslice runtime.growslice
func growslice() {
	panic()
}

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
//go:linkname memmove runtime.memmove
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
//go:linkname cmpstring runtime.cmpstring
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
//go:linkname typedmemmove runtime.typedmemmove
func typedmemmove(psize *uintptr, dst *byte, src *byte) {
	memmove(dst, src, *psize)
}

//go:nosplit
//go:linkname typedslicecopy runtime.typedslicecopy
func typedslicecopy(psize *uintptr, dst slice, src slice) int {
	n := dst.len
	if src.len < n {
		n = src.len
	}
	memmove(dst.array, src.array, *psize*n)
	return int(n)
}
