package buffer

import (
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/pointers"
	"github.com/anton2920/gofa/syscall"
)

type Circular struct {
	Buffer []byte
	Head   int
	Tail   int
}

func NewCircular(size int) (Circular, error) {
	var c Circular

	const pageSize = 4096
	size = ints.AlignUp(size, pageSize)

	/* NOTE(anton2920): first argument is SHM_ANON, cannot have that as a variable as Go's checkptr doesn't like it. */
	fd, err := syscall.ShmOpen2(*(*string)(unsafe.Pointer(&reflect.StringHeader{Data: 1, Len: 8})), syscall.O_RDWR, 0, 0, syscall.NULL)
	if err != nil {
		return c, err
	}

	defer syscall.Close(fd)

	if err := syscall.Ftruncate(fd, int64(size)); err != nil {
		return c, err
	}

	buffer, err := syscall.Mmap(nil, 2*uint64(size), syscall.PROT_NONE, syscall.MAP_PRIVATE|syscall.MAP_ANON, -1, 0)
	if err != nil {
		return c, err
	}

	if _, err := syscall.Mmap(buffer, uint64(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, fd, 0); err != nil {
		return c, err
	}
	if _, err := syscall.Mmap(pointers.Add(buffer, size), uint64(size), syscall.PROT_READ|syscall.PROT_WRITE, syscall.MAP_SHARED|syscall.MAP_FIXED, fd, 0); err != nil {
		return c, err
	}
	c.Buffer = *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(buffer), Len: 2 * size, Cap: 2 * size}))

	/* NOTE(anton2920): sanity checks. */
	c.Buffer[0] = '\x00'
	c.Buffer[size-1] = '\x00'
	c.Buffer[size] = '\x00'
	c.Buffer[2*size-1] = '\x00'

	return c, nil
}

func (c *Circular) Consume(n int) {
	c.Head += n
	if c.Head > len(c.Buffer)/2 {
		c.Head -= len(c.Buffer) / 2
		c.Tail -= len(c.Buffer) / 2
	}
}

func (c *Circular) Produce(n int) {
	c.Tail += n
}

func (c *Circular) RemainingSlice() []byte {
	return c.Buffer[c.Tail : c.Head+len(c.Buffer)/2]
}

func (c *Circular) RemainingSpace() int {
	return (len(c.Buffer) / 2) - (c.Tail - c.Head)
}

func (c *Circular) Reset() {
	c.Head = 0
	c.Tail = 0
}

func (c *Circular) UnconsumedLen() int {
	return c.Tail - c.Head
}

func (c *Circular) UnconsumedSlice() []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&c.Buffer[c.Head])), Len: c.UnconsumedLen(), Cap: c.UnconsumedLen()}))
}

func (c *Circular) UnconsumedString() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&c.Buffer[c.Head])), Len: c.UnconsumedLen()}))
}

func (c *Circular) Free() error {
	err := syscall.Munmap(unsafe.Pointer(&c.Buffer[0]), uint64(len(c.Buffer)))
	c.Buffer = nil
	return err
}
