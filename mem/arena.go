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

type ArenaSave struct {
	Arena    *Arena
	PrevOfft int
	CurrOfft int
}

const arenaDefaultAlignment = unsafe.Alignof(uintptr(0))

func NewArena(buffer []byte) Arena {
	return Arena{Buffer: buffer}
}

func NewArenaFromBytePointer(ptr *byte, n int) Arena {
	return NewArena(bytes.SliceFromBytePointer(ptr, n))
}

func NewArenaFromUnsafePointer(ptr unsafe.Pointer, n int) Arena {
	return NewArena(bytes.SliceFromUnsafePointer(ptr, n))
}

func (a *Arena) Allocate(n int) []byte {
	return a.AllocateWithAlignment(n, arenaDefaultAlignment)
}

func (a *Arena) AllocateWithAlignment(n int, align uintptr) []byte {
	t := trace.Begin("")

	a.CurrOfft += pointers.Delta(pointers.AlignUp(unsafe.Pointer(&a.Buffer[a.CurrOfft]), align), unsafe.Pointer(&a.Buffer[a.CurrOfft]))
	if a.CurrOfft+n <= cap(a.Buffer) {
		debug.AssertZero(int(uintptr(unsafe.Pointer(&a.Buffer[a.CurrOfft]))&(align-1)), "(*Arena).AllocateWithAlignment tried to return unaligned pointer")
		a.PrevOfft = a.CurrOfft
		a.CurrOfft += n
		buf := a.Buffer[a.PrevOfft:a.CurrOfft]

		for i := 0; i < len(buf); i++ {
			buf[i] = 0
		}

		trace.End(t)
		return buf
	}

	debug.Printf("no more space in arena %p: requested %d, got only %d left", &a.Buffer[0], n, cap(a.Buffer)-a.CurrOfft)
	trace.End(t)
	return nil
}

func (a *Arena) Begin() ArenaSave {
	return ArenaSave{Arena: a, PrevOfft: a.PrevOfft, CurrOfft: a.CurrOfft}
}

func (a *Arena) Reallocate(old []byte, n int) []byte {
	return a.ReallocateWithAlignment(old, n, arenaDefaultAlignment)
}

func (a *Arena) ReallocateWithAlignment(old []byte, n int, align uintptr) []byte {
	t := trace.Begin("")

	if len(old) == 0 {
		buf := a.AllocateWithAlignment(n, align)
		trace.End(t)
		return buf
	} else if (uintptr(unsafe.Pointer(&old[0])) >= uintptr(unsafe.Pointer(&a.Buffer[0]))) && (uintptr(unsafe.Pointer(&old[len(old)-1])) <= uintptr(unsafe.Pointer(&a.Buffer[a.CurrOfft-1]))) {
		if &old[0] == &a.Buffer[a.PrevOfft] {
			if n > len(old) {
				for i := len(old); i < n; i++ {
					a.Buffer[a.PrevOfft+i] = 0
				}
			}
			a.CurrOfft = a.PrevOfft + n
			return a.Buffer[a.PrevOfft:a.CurrOfft]
		} else {
			buf := a.AllocateWithAlignment(n, align)
			copy(buf, old)
			return buf
		}
	}

	trace.End(t)
	panic("old memory does not come from this arena")
}

func (a *Arena) PushInt() *int {
	buf := a.AllocateWithAlignment(int(unsafe.Sizeof(int(0))), unsafe.Alignof(int(0)))
	return (*int)(unsafe.Pointer(&buf[0]))
}

func (a *Arena) PushFloat32() *float32 {
	buf := a.AllocateWithAlignment(int(unsafe.Sizeof(float32(0))), unsafe.Alignof(float32(0)))
	return (*float32)(unsafe.Pointer(&buf[0]))
}

func (a *Arena) Reset() {
	a.PrevOfft = 0
	a.CurrOfft = 0
}

func (a *ArenaSave) Rollback() {
	a.Arena.PrevOfft = a.PrevOfft
	a.Arena.CurrOfft = a.CurrOfft
}
