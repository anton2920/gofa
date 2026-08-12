package mem_

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/os"
)

type CircularBuffer struct {
	Buffer []byte
	Head   int
	Tail   int
}

func (c *CircularBuffer) Init(ctx *context.Context, size int) bool {
	const value = 0xDEADBEEF

	ptr, ok := os.CreateCircularMemoryMapping(ctx, size, 2*size)
	if !ok {
		return false
	}
	c.Buffer = bytes.SliceFromUnsafePointer(ptr, 2*size)

	*(*uint)(unsafe.Pointer(&c.Buffer[0])) = value
	if *(*uint)(unsafe.Pointer(&c.Buffer[size])) != value {
		c.Free(ctx)
		return false
	}
	*(*uint)(unsafe.Pointer(&c.Buffer[0])) = 0

	return true
}

func (c *CircularBuffer) Consume(n int) {
	c.Head += n
	if c.Head > len(c.Buffer)/2 {
		c.Head -= len(c.Buffer) / 2
		c.Tail -= len(c.Buffer) / 2
	}
}

func (c *CircularBuffer) Produce(n int) {
	c.Tail += n
}

func (c *CircularBuffer) RemainingSlice() []byte {
	return c.Buffer[c.Tail : c.Head+len(c.Buffer)/2]
}

func (c *CircularBuffer) RemainingSpace() int {
	return (len(c.Buffer) / 2) - (c.Tail - c.Head)
}

func (c *CircularBuffer) Reset() {
	c.Head = 0
	c.Tail = 0
}

func (c *CircularBuffer) UnconsumedLen() int {
	return c.Tail - c.Head
}

func (c *CircularBuffer) UnconsumedSlice() []byte {
	//return *(*[]byte)(unsafe.Pointer(&types.SliceHeader{Data: uintptr(unsafe.Pointer(&c.Buffer[c.Head])), Len: c.UnconsumedLen(), Cap: c.UnconsumedLen()}))
	return c.Buffer[c.Head : c.Head+c.UnconsumedLen()]
}

func (c *CircularBuffer) UnconsumedString() string {
	//return *(*string)(unsafe.Pointer(&types.StringHeader{Data: uintptr(unsafe.Pointer(&c.Buffer[c.Head])), Len: c.UnconsumedLen()}))
	return bytes.AsString(c.UnconsumedSlice())
}

func (c *CircularBuffer) Free(ctx *context.Context) bool {
	return os.DeallocateVirtualMemory(ctx, unsafe.Pointer(&c.Buffer[0]), len(c.Buffer)*2)
}
