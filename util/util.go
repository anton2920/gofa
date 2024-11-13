package util

import (
	"reflect"
	"unsafe"
)

func AlignUp(x int, quantum int) int {
	return (x + (quantum - 1)) & ^(quantum - 1)
}

func Bool2Int(b bool) int {
	if b {
		return 1
	}
	return 0
}

/* Clamp returns number clamped into a range from l inclusive to r exclusive. */
func Clamp(x int, l int, r int) int {
	if x > r-1 {
		x = r - 1
	}
	if x < l {
		x = l
	}
	return x
}

/* GetCallerPC returns a value of %IP register that is going to be used by RET instruction. arg0 is the address of the first agrument function of interest accepts. */
func GetCallerPC(arg0 unsafe.Pointer) uintptr {
	return *(*uintptr)(PtrAdd(arg0, -int(unsafe.Sizeof(arg0))))
}

func Max(a, b int) int {
	if a > b {
		return a
	} else {
		return b
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	} else {
		return b
	}
}

func MoveIntDown(vs []int, i int) {
	if (i >= 0) && (i < len(vs)-1) {
		vs[i], vs[i+1] = vs[i+1], vs[i]
	}
}

func MoveStringDown(vs []string, i int) {
	if (i >= 0) && (i < len(vs)-1) {
		vs[i], vs[i+1] = vs[i+1], vs[i]
	}
}

func MoveIntUp(vs []int, i int) {
	if (i > 0) && (i <= len(vs)-1) {
		vs[i-1], vs[i] = vs[i], vs[i-1]
	}
}

func MoveStringUp(vs []string, i int) {
	if (i > 0) && (i <= len(vs)-1) {
		vs[i-1], vs[i] = vs[i], vs[i-1]
	}
}

/* Noescape hides a pointer from escape analysis. Noescape is the identity function but escape analysis doesn't think the output depends on the input. Noescape is inlined and currently compiles down to zero instructions. */
//go:nosplit
func Noescape(p unsafe.Pointer) unsafe.Pointer {
	x := uintptr(p)
	return unsafe.Pointer(x ^ 0)
}

func PtrAdd(ptr unsafe.Pointer, x int) unsafe.Pointer {
	return unsafe.Pointer(uintptr(ptr) + uintptr(x))
}

func Slice2String(s []byte) string {
	return *(*string)(unsafe.Pointer(&s))
}

func String2Slice(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: (*reflect.StringHeader)(unsafe.Pointer(&s)).Data, Len: len(s), Cap: len(s)}))
}

func StringData(s string) *byte {
	return (*byte)(unsafe.Pointer((*reflect.StringHeader)(unsafe.Pointer(&s)).Data))
}

func SwapBytesInWord(x uint16) uint16 {
	return ((x << 8) & 0xFF00) | (x >> 8)
}

func RemoveIntAtIndex(vs []int, i int) []int {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}

func RemoveInt32AtIndex(vs []int32, i int) []int32 {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}

func RemoveStringAtIndex(vs []string, i int) []string {
	if (len(vs) == 0) || (i < 0) || (i >= len(vs)) {
		return vs
	}
	if i < len(vs)-1 {
		copy(vs[i:], vs[i+1:])
	}
	return vs[:len(vs)-1]
}
