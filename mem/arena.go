/* Inspired by https://www.gingerbill.org/article/2019/02/08/memory-allocation-strategies-002/. */
package mem

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/pointers"
	"github.com/anton2920/gofa/trace"
)

type Arena struct {
	Buffer   []byte
	PrevOfft int
	CurrOfft int
}

type ArenaSavePoint struct {
	Arena    *Arena
	PrevOfft int
	CurrOfft int
}

const arenaDefaultAlignment = unsafe.Alignof(uintptr(0))

func ArenaFromByteSlice(buf []byte) Arena {
	return Arena{Buffer: bytes.SliceFromBytePointer(&buf[0], len(buf))}
}

func ArenaFromBytePointer(ptr *byte, n int) Arena {
	return ArenaFromByteSlice(bytes.SliceFromBytePointer(ptr, n))
}

func ArenaFromUnsafePointer(ptr unsafe.Pointer, n int) Arena {
	return ArenaFromByteSlice(bytes.SliceFromUnsafePointer(ptr, n))
}

func (a *Arena) InitWithByteSlice(buf []byte) {
	a.Buffer = bytes.SliceFromBytePointer(&buf[0], len(buf))
	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) InitWithBytePointer(ptr *byte, n int) {
	a.Buffer = bytes.SliceFromBytePointer(ptr, n)
	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) InitWithUnsafePointer(ptr unsafe.Pointer, n int) {
	a.Buffer = bytes.SliceFromUnsafePointer(ptr, n)
	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) PushSize(n int) unsafe.Pointer {
	return a.PushSizeWithAlignment(n, arenaDefaultAlignment)
}

func (a *Arena) PushSizeWithAlignment(n int, align uintptr) unsafe.Pointer {
	t := trace.Begin("")

	a.CurrOfft += pointers.Delta(pointers.AlignUp(unsafe.Pointer(&a.Buffer[a.CurrOfft]), align), unsafe.Pointer(&a.Buffer[a.CurrOfft]))
	if a.CurrOfft+n <= len(a.Buffer) {
		debug.AssertZero(int(uintptr(unsafe.Pointer(&a.Buffer[a.CurrOfft]))&(align-1)), "(*Arena).PushSizeWithAlignment tried to return unaligned pointer")
		a.PrevOfft = a.CurrOfft
		a.CurrOfft += n
		buf := a.Buffer[a.PrevOfft:a.CurrOfft]

		for i := 0; i < len(buf); i++ {
			buf[i] = 0
		}

		trace.End(t)
		return pointers.UnsafeNoescape(unsafe.Pointer(&buf[0]))
	}

	debug.Printf("no more space in arena %p: requested %d, got only %d left", &a.Buffer[0], n, len(a.Buffer)-a.CurrOfft)
	trace.End(t)
	return nil
}

func (a *Arena) RepushSize(optr unsafe.Pointer, on int, nn int) unsafe.Pointer {
	return a.RepushSizeWithAlignment(optr, on, nn, arenaDefaultAlignment)
}

func (a *Arena) RepushSizeWithAlignment(optr unsafe.Pointer, on int, nn int, align uintptr) unsafe.Pointer {
	t := trace.Begin("")

	if on == 0 {
		nptr := a.PushSizeWithAlignment(nn, align)
		trace.End(t)
		return nptr
	} else if (uintptr(optr) >= uintptr(unsafe.Pointer(&a.Buffer[0]))) && ((uintptr(optr) + uintptr(on)) <= uintptr(unsafe.Pointer(&a.Buffer[a.CurrOfft-1]))) {
		if optr == unsafe.Pointer(&a.Buffer[a.PrevOfft]) {
			if nn > on {
				for i := on; i < nn; i++ {
					a.Buffer[a.PrevOfft+i] = 0
				}
			}
			a.CurrOfft = a.PrevOfft + nn
			return optr
		} else {
			nptr := a.PushSizeWithAlignment(nn, align)
			copy(bytes.SliceFromUnsafePointer(nptr, nn), bytes.SliceFromUnsafePointer(optr, on))
			return nptr
		}
	}

	trace.End(t)
	panic("old memory does not come from this arena")
}

func (a *Arena) PushInt() *int {
	return (*int)(a.PushSizeWithAlignment(int(unsafe.Sizeof(int(0))), unsafe.Alignof(int(0))))
}

func (a *Arena) PushFloat32() *float32 {
	return (*float32)(a.PushSizeWithAlignment(int(unsafe.Sizeof(float32(0))), unsafe.Alignof(float32(0))))
}

func (a *Arena) PushByteArray(n int) []byte {
	return bytes.SliceFromUnsafePointer(a.PushSizeWithAlignment(n, unsafe.Alignof(byte(0))), n)
}

func (a *Arena) RepushByteArray(old []byte, n int) []byte {
	return bytes.SliceFromUnsafePointer(a.RepushSizeWithAlignment(unsafe.Pointer(&old[0]), len(old), n, unsafe.Alignof(byte(0))), n)
}

func (a *Arena) Reset() {
	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) Begin() ArenaSavePoint {
	return ArenaSavePoint{Arena: a, PrevOfft: a.PrevOfft, CurrOfft: a.CurrOfft}
}

func (a *ArenaSavePoint) Rollback() {
	a.Arena.PrevOfft = a.PrevOfft
	a.Arena.CurrOfft = a.CurrOfft
}
