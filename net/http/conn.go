package http

import (
	"reflect"

	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/pointers"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/trace"
)

type Version int32

const (
	VersionNone = Version(iota)
	Version09
	Version10
	Version11
	Version20
	Version30
	VersionCount
)

var Version2String = [...]string{
	VersionNone: "HTTP/1.1",
	Version09:   "HTTP/0.9",
	Version10:   "HTTP/1.0",
	Version11:   "HTTP/1.1",
}

type Conn struct {
	ConnPool *ConnPool

	Version Version
	Socket  os.Handle

	RequestBuffer buffer.Circular

	ResponseBuffer []byte
	ResponsePos    int64

	Error error

	remoteAddr      [21]byte
	CloseAfterWrite bool
	Closed          bool
	Check           uint8 /* NOTE(anton2920): Check must be the same as the last pointer's bit, if context is in use. */
}

type ConnOptions struct {
	RequestBufferSize int
}

func MergeConnOptions(opts ...ConnOptions) ConnOptions {
	var result ConnOptions

	for i := 0; i < len(opts); i++ {
		opt := &opts[i]

		ints.Replace(&result.RequestBufferSize, opt.RequestBufferSize)
	}

	return result
}

func (c *Conn) RemoteAddr() string {
	return *(*string)(unsafe.Pointer(&reflect.StringHeader{Data: uintptr(unsafe.Pointer(&c.remoteAddr)), Len: len(c.remoteAddr)}))
}

func (c *Conn) Close() error {
	if c.Closed {
		return nil
	}

	c.CloseAfterWrite = false
	c.ResponseBuffer = nil
	c.RequestBuffer.Free()
	c.ResponsePos = 0
	c.Version = 0
	c.Error = nil

	c.Closed = true
	c.Check = 1 - c.Check
	err := os.Close(c.Socket)
	c.ConnPool.Put(c)

	return err
}

func (c *Conn) ReadRequestData() (int64, error) {
	t := trace.Begin("")

	buf := c.RequestBuffer.RemainingSlice()
	if len(buf) == 0 {
		c.Error = RequestEntityTooLarge("no space left in buffer")
		trace.End(t)
		return 0, nil
	}

	n, err := syscall.Read(int32(c.Socket), buf)
	if err != nil {
		trace.End(t)
		return -1, err
	}
	c.RequestBuffer.Produce(int(n))

	trace.End(t)
	return n, nil
}

func (c *Conn) WriteResponseData() (int64, error) {
	t := trace.Begin("")

	var err error
	var n int64

	if len(c.ResponseBuffer[c.ResponsePos:]) > 0 {
		n, err = syscall.Write(int32(c.Socket), c.ResponseBuffer[c.ResponsePos:])
		if err != nil {
			trace.End(t)
			return -1, err
		}
		c.ResponsePos += n

		if c.ResponsePos == int64(len(c.ResponseBuffer)) {
			c.ResponseBuffer = c.ResponseBuffer[:0]
			c.ResponsePos = 0
			if c.CloseAfterWrite {
				err = c.Close()
			}
		}
	}

	trace.End(t)
	return n, err
}

func (c *Conn) Pointer() unsafe.Pointer {
	return pointers.Add(unsafe.Pointer(c), int(c.Check))
}

func ConnFromPointer(ptr unsafe.Pointer) (*Conn, bool) {
	if ptr == nil {
		return nil, false
	}

	check := uintptr(ptr) & 0x1
	c := (*Conn)(unsafe.Pointer(uintptr(ptr) - check))

	return c, c.Check == uint8(check)
}

func RequestBufferSize(size int) ConnOptions {
	return ConnOptions{RequestBufferSize: size}
}
