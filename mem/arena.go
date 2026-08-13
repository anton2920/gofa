/* Inspired by https://www.gingerbill.org/article/2019/02/08/memory-allocation-strategies-002/. */
package mem

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug/debug_"
	"github.com/anton2920/gofa/pointers"
	"github.com/anton2920/gofa/trace/trace_"
)

type Arena struct {
	Base unsafe.Pointer
	Size int

	PrevOfft int
	CurrOfft int
}

const arenaDefaultAlignment = unsafe.Alignof(uintptr(0))

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

func (a *Arena) AllocationComesFromHere(ptr unsafe.Pointer, n int) bool {
	return (uintptr(ptr) >= uintptr(a.Base)) && ((uintptr(ptr) + uintptr(n)) <= (uintptr(a.Base) + uintptr(a.CurrOfft)))
}

func (a *Arena) PushSizeWithAlignment(n int, align uintptr) unsafe.Pointer {
	t := trace_.Begin("")

	curr := pointers.Add(a.Base, uintptr(a.CurrOfft))
	a.CurrOfft += pointers.Delta(pointers.AlignUpPow2(curr, align), curr)
	if a.CurrOfft+n <= a.Size {
		curr := pointers.Add(a.Base, uintptr(a.CurrOfft))
		debug_.AssertZero(int(uintptr(curr)&(align-1)), "(*Arena).PushSizeWithAlignment tried to return unaligned pointer")

		a.PrevOfft = a.CurrOfft
		a.CurrOfft += n
		bytes.SliceClear(bytes.SliceFromUnsafePointer(a.Base, a.Size)[a.PrevOfft:a.CurrOfft])

		trace_.End(t)
		return curr
	}

	trace_.End(t)

	/* TODO(anton2920): handle shortage of memory better. */
	panic("BUY MORE RAM!")
}

func (a *Arena) PushSize(n int) unsafe.Pointer {
	return a.PushSizeWithAlignment(n, arenaDefaultAlignment)
}

func (a *Arena) RepushSizeWithAlignment(optr unsafe.Pointer, on int, nn int, align uintptr) unsafe.Pointer {
	t := trace_.Begin("")

	if on == 0 {
		nptr := a.PushSizeWithAlignment(nn, align)
		trace_.End(t)
		return nptr
	} else if a.AllocationComesFromHere(optr, on) {
		if optr == pointers.Add(a.Base, uintptr(a.PrevOfft)) {
			if nn > on {
				bytes.SliceClear(bytes.SliceFromUnsafePointer(a.Base, a.Size)[a.PrevOfft+on : a.CurrOfft])
			}
			a.CurrOfft = a.PrevOfft + nn
			trace_.End(t)
			return optr
		} else {
			nptr := a.PushSizeWithAlignment(nn, align)
			copy(bytes.SliceFromUnsafePointer(nptr, int(nn)), bytes.SliceFromUnsafePointer(optr, int(on)))
			trace_.End(t)
			return nptr
		}
	}

	trace_.End(t)
	panic("old memory does not come from this arena")
}

func (a *Arena) RepushSize(optr unsafe.Pointer, on int, nn int) unsafe.Pointer {
	return a.RepushSizeWithAlignment(optr, on, nn, arenaDefaultAlignment)
}

func (a *Arena) Reset() {
	a.PrevOfft = 0
	a.CurrOfft = 0
}
