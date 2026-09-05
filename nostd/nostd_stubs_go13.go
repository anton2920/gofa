//go:build gofanostd && gofanostd13
// +build gofanostd,gofanostd13

package nostd

import "unsafe"

//go:nosplit
func memhash(h *uintptr, s uintptr, p *byte) {
	const (
		ptrSize = unsafe.Sizeof(uintptr(0))
		c0      = uintptr((8-ptrSize)/4*2860486313 + (ptrSize-4)/4*33054211828000289)
		c1      = uintptr((8-ptrSize)/4*3267000013 + (ptrSize-4)/4*23344194077549503)
	)

	*h ^= c0
	for s > 0 {
		*h = (*h ^ uintptr(*p)) * c1
		p = inc(p)
		s--
	}
	return
}

//go:nosplit
func memequal(eq *bool, s uintptr, a *byte, b *byte) {
	if a == b {
		*eq = true
		return
	}
	for s > 0 {
		if *a != *b {
			return
		}
		a = inc(a)
		b = inc(b)
		s--
	}
	return
}

//go:nosplit
func memequal0(eq *bool, s uintptr, _ unsafe.Pointer, _ unsafe.Pointer) {
	*eq = true
}

//go:nosplit
func memequal8(eq *bool, s uintptr, a *uint8, b *uint8) {
	*eq = (*a == *b)
}

//go:nosplit
func memequal16(eq *bool, s uintptr, a *uint16, b *uint16) {
	*eq = (*a == *b)
}

//go:nosplit
func memequal32(eq *bool, s uintptr, a *uint32, b *uint32) {
	*eq = (*a == *b)
}

//go:nosplit
func memequal64(eq *bool, s uintptr, a *uint64, b *uint64) {
	*eq = (*a == *b)
}

//go:nosplit
func memequal128(eq *bool, s uintptr, a *[2]uint64, b *[2]uint64) {
	*eq = (a[0] == b[0]) && (a[1] == b[1])
}

//go:nosplit
func memcopy(length uintptr, dst *byte, src *byte) {
	memmove(dst, src, length)
}

//go:nosplit
func memcopy8(_ uintptr, dst *uint8, src *uint8) {
	if src == nil {
		*dst = 0
		return
	}
	*dst = *src
}

//go:nosplit
func memcopy16(_ uintptr, dst *uint16, src *uint16) {
	if src == nil {
		*dst = 0
		return
	}
	*dst = *src
}

//go:nosplit
func memcopy32(_ uintptr, dst *uint32, src *uint32) {
	if src == nil {
		*dst = 0
		return
	}
	*dst = *src
}

//go:nosplit
func memcopy64(_ uintptr, dst *uint64, src *uint64) {
	if src == nil {
		*dst = 0
		return
	}
	*dst = *src
}

//go:nosplit
func memcopy128(_ uintptr, dst *[2]uint64, src *[2]uint64) {
	if src == nil {
		dst[0] = 0
		dst[1] = 0
		return
	}
	dst[0] = src[0]
	dst[1] = src[1]
}

//go:nosplit
func strhash(h *uintptr, _ uintptr, a *byte) {
	memhash(h, uintptr(len(*(*string)(unsafe.Pointer(a)))), *(**byte)(unsafe.Pointer(a)))
}
