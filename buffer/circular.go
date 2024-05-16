package buffer

import (
	"unsafe"

	"github.com/anton2920/gofa/syscall"
)

type Circular struct {
	SourceBuffer unsafe.Pointer

	Buf  []byte
	Head int
	Tail int
}

const (
	/* See <sys/mman.h>. */
	PROT_NONE  = 0x00
	PROT_READ  = 0x01
	PROT_WRITE = 0x02

	MAP_SHARED  = 0x0001
	MAP_PRIVATE = 0x0002

	MAP_FIXED = 0x0010
	MAP_ANON  = 0x1000
)

var SHM_ANON = unsafe.String((*byte)(unsafe.Pointer(uintptr(1))), 8)
var NULL = unsafe.String(nil, 0)

func NewCircular(size int) (Circular, error) {
	var c Circular

	/* NOTE(anton2920): rounding up to the page boundary. */
	size = (size + (4096 - 1)) & ^(4096 - 1)

	fd, err := syscall.ShmOpen2(SHM_ANON, syscall.O_RDWR, 0, 0, NULL)
	if err != nil {
		return c, err
	}
	defer syscall.Close(fd)

	if err := syscall.Ftruncate(fd, int64(size)); err != nil {
		return c, err
	}

	buffer, err := syscall.Mmap(nil, 2*uint64(size), PROT_NONE, MAP_PRIVATE|MAP_ANON, -1, 0)
	if err != nil {
		return c, err
	}

	if _, err := syscall.Mmap(buffer, uint64(size), PROT_READ|PROT_WRITE, MAP_SHARED|MAP_FIXED, fd, 0); err != nil {
		return c, err
	}
	if _, err := syscall.Mmap(unsafe.Add(buffer, size), uint64(size), PROT_READ|PROT_WRITE, MAP_SHARED|MAP_FIXED, fd, 0); err != nil {
		return c, err
	}

	c.SourceBuffer = buffer
	c.Buf = unsafe.Slice((*byte)(buffer), 2*size)

	/* NOTE(anton2920): sanity checks. */
	c.Buf[0] = '\x00'
	c.Buf[size-1] = '\x00'
	c.Buf[size] = '\x00'
	c.Buf[2*size-1] = '\x00'

	return c, nil
}

func (c *Circular) Consume(n int) {
	c.Head += n
	if c.Head > len(c.Buf)/2 {
		c.Head -= len(c.Buf) / 2
		c.Tail -= len(c.Buf) / 2
	}
}

func (c *Circular) Produce(n int) {
	c.Tail += n
}

func (c *Circular) RemainingSlice() []byte {
	return c.Buf[c.Tail : c.Head+len(c.Buf)/2]
}

func (c *Circular) RemainingSpace() int {
	return (len(c.Buf) / 2) - (c.Tail - c.Head)
}

func (c *Circular) Reset() {
	c.Head = 0
	c.Tail = 0
}

func (c *Circular) UnconsumedLen() int {
	return c.Tail - c.Head
}

func (c *Circular) UnconsumedSlice() []byte {
	return unsafe.Slice(&c.Buf[c.Head], c.UnconsumedLen())
}

func (c *Circular) UnconsumedString() string {
	return unsafe.String(&c.Buf[c.Head], c.UnconsumedLen())
}

func FreeCircular(c *Circular) {
	syscall.Munmap(unsafe.Pointer(unsafe.SliceData(c.Buf)), uint64(len(c.Buf)))
	syscall.Munmap(c.SourceBuffer, uint64(len(c.Buf)))
}
