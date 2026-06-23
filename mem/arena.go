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
	Base unsafe.Pointer
	Size uint

	PrevOfft uint
	CurrOfft uint
}

type ArenaSavePoint struct {
	Arena    *Arena
	PrevOfft uint
	CurrOfft uint
}

const arenaDefaultAlignment = unsafe.Alignof(uintptr(0))

func ArenaFromByteSlice(buf []byte) Arena {
	return ArenaFromBytePointer(&buf[0], uint(len(buf)))
}

func ArenaFromBytePointer(ptr *byte, n uint) Arena {
	return ArenaFromUnsafePointer(unsafe.Pointer(ptr), n)
}

func ArenaFromUnsafePointer(ptr unsafe.Pointer, n uint) Arena {
	return Arena{Base: ptr, Size: n}
}

func (a *Arena) InitWithByteSlice(buf []byte) {
	a.InitWithBytePointer(&buf[0], uint(len(buf)))
}

func (a *Arena) InitWithBytePointer(ptr *byte, n uint) {
	a.InitWithUnsafePointer(unsafe.Pointer(ptr), n)
}

func (a *Arena) InitWithUnsafePointer(ptr unsafe.Pointer, n uint) {
	a.Base = pointers.UnsafeNoescape(ptr)
	a.Size = n

	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) PushSize(n uint) unsafe.Pointer {
	return a.PushSizeWithAlignment(n, arenaDefaultAlignment)
}

func (a *Arena) PushSizeWithAlignment(n uint, align uintptr) unsafe.Pointer {
	t := trace.Begin("")

	curr := pointers.Add(a.Base, uintptr(a.CurrOfft))
	a.CurrOfft += uint(pointers.Delta(pointers.AlignUpPow2(curr, align), curr))
	if a.CurrOfft+n <= a.Size {
		curr := pointers.Add(a.Base, uintptr(a.CurrOfft))
		debug.AssertZero(int(uintptr(curr)&(align-1)), "(*Arena).PushSizeWithAlignment tried to return unaligned pointer")
		a.PrevOfft = a.CurrOfft
		a.CurrOfft += n

		buf := bytes.SliceFromUnsafePointer(curr, int(a.CurrOfft-a.PrevOfft))
		for i := 0; i < len(buf); i++ {
			buf[i] = 0
		}

		trace.End(t)
		return curr
	}

	debug.Printf("no more space in arena %p: requested %d, got only %d left", a.Base, n, a.Size-a.CurrOfft)
	trace.End(t)
	return nil
}

func (a *Arena) RepushSize(optr unsafe.Pointer, on uint, nn uint) unsafe.Pointer {
	return a.RepushSizeWithAlignment(optr, on, nn, arenaDefaultAlignment)
}

func (a *Arena) RepushSizeWithAlignment(optr unsafe.Pointer, on uint, nn uint, align uintptr) unsafe.Pointer {
	t := trace.Begin("")

	if on == 0 {
		nptr := a.PushSizeWithAlignment(nn, align)
		trace.End(t)
		return nptr
	} else if (uintptr(optr) >= uintptr(a.Base)) && ((uintptr(optr) + uintptr(on)) <= (uintptr(a.Base) + uintptr(a.CurrOfft))) {
		if optr == pointers.Add(a.Base, uintptr(a.PrevOfft)) {
			if nn > on {
				buf := bytes.SliceFromUnsafePointer(a.Base, int(a.Size))
				for i := on; i < nn; i++ {
					buf[a.PrevOfft+i] = 0
				}
			}
			a.CurrOfft = a.PrevOfft + nn
			trace.End(t)
			return optr
		} else {
			nptr := a.PushSizeWithAlignment(nn, align)
			copy(bytes.SliceFromUnsafePointer(nptr, int(nn)), bytes.SliceFromUnsafePointer(optr, int(on)))
			trace.End(t)
			return nptr
		}
	}

	trace.End(t)
	panic("old memory does not come from this arena")
}

func (a *Arena) PushInt() *int {
	return (*int)(a.PushSizeWithAlignment(uint(unsafe.Sizeof(int(0))), unsafe.Alignof(int(0))))
}

func (a *Arena) PushFloat32() *float32 {
	return (*float32)(a.PushSizeWithAlignment(uint(unsafe.Sizeof(float32(0))), unsafe.Alignof(float32(0))))
}

func (a *Arena) PushByteArray(n uint) []byte {
	return bytes.SliceFromUnsafePointer(a.PushSizeWithAlignment(n, unsafe.Alignof(byte(0))), int(n))
}

func (a *Arena) RepushByteArray(old []byte, n uint) []byte {
	return bytes.SliceFromUnsafePointer(a.RepushSizeWithAlignment(unsafe.Pointer(&old[0]), uint(len(old)), n, unsafe.Alignof(byte(0))), int(n))
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
