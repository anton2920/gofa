/* Inspired by https://www.gingerbill.org/article/2019/02/08/memory-allocation-strategies-002/. */
package mem

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug/debug_"
	"github.com/anton2920/gofa/pointers"
)

type Arena struct {
	Base unsafe.Pointer
	Size int

	PrevOfft uintptr
	CurrOfft uintptr
}

const arenaDefaultAlignment = 2 * unsafe.Alignof(uintptr(0))

func (a *Arena) InitWithByteSlice(buf []byte) {
	a.InitWithBytePointer(&buf[0], len(buf))
}

func (a *Arena) InitWithBytePointer(ptr *byte, n int) {
	a.InitWithUnsafePointer(unsafe.Pointer(ptr), n)
}

func (a *Arena) InitWithUnsafePointer(ptr unsafe.Pointer, n int) {
	a.Base = pointers.UnsafeNoescape(ptr)
	a.Size = n

	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *Arena) AllocationComesFromHere(ptr unsafe.Pointer, n uintptr) bool {
	return (uintptr(ptr) >= uintptr(a.Base)) && ((uintptr(ptr) + n) <= (uintptr(a.Base) + a.CurrOfft))
}

var buyMoreRamMsg interface{} = "BUY MORE RAM!"

func (a *Arena) PushSizeWithAlignment(n uintptr, align uintptr) unsafe.Pointer {
	curr := pointers.Add(a.Base, a.CurrOfft)
	a.CurrOfft += uintptr(pointers.Delta(pointers.AlignUpPow2(curr, align), curr))
	if a.CurrOfft+n <= uintptr(a.Size) {
		curr := pointers.Add(a.Base, a.CurrOfft)
		debug_.AssertZero(int(uintptr(curr)&(align-1)), "(*Arena).PushSizeWithAlignment tried to return unaligned pointer")

		a.PrevOfft = a.CurrOfft
		a.CurrOfft += n
		bytes.SliceClear(bytes.SliceFromUnsafePointer(a.Base, a.Size)[a.PrevOfft:a.CurrOfft])

		return curr
	}

	/* TODO(anton2920): handle shortage of memory better. */
	panic(buyMoreRamMsg)
}

func (a *Arena) PushSize(n uintptr) unsafe.Pointer {
	return a.PushSizeWithAlignment(n, arenaDefaultAlignment)
}

var oldMemoryIsNotFromHereMsg interface{} = "old memory does not come from this arena"

func (a *Arena) RepushSizeWithAlignment(optr unsafe.Pointer, on uintptr, nn uintptr, align uintptr) unsafe.Pointer {
	if on == 0 {
		nptr := a.PushSizeWithAlignment(nn, align)
		return nptr
	} else if a.AllocationComesFromHere(optr, on) {
		if optr == pointers.Add(a.Base, uintptr(a.PrevOfft)) {
			if nn > on {
				bytes.SliceClear(bytes.SliceFromUnsafePointer(a.Base, a.Size)[a.PrevOfft+on : a.CurrOfft])
			}
			a.CurrOfft = a.PrevOfft + nn
			return optr
		} else {
			nptr := a.PushSizeWithAlignment(nn, align)
			copy(bytes.SliceFromUnsafePointer(nptr, int(nn)), bytes.SliceFromUnsafePointer(optr, int(on)))
			return nptr
		}
	}

	panic(oldMemoryIsNotFromHereMsg)
}

func (a *Arena) RepushSize(optr unsafe.Pointer, on uintptr, nn uintptr) unsafe.Pointer {
	return a.RepushSizeWithAlignment(optr, on, nn, arenaDefaultAlignment)
}

func (a *Arena) Reset() {
	a.PrevOfft = 0
	a.CurrOfft = 0
}
