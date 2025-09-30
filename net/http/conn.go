package http

import (
	"unsafe"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/os"
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
	Version09: "HTTP/0.9",
	Version10: "HTTP/1.0",
	Version11: "HTTP/1.1",
}

type Conn struct {
	/* NOTE(anton2920): Check must be the same as the last pointer's bit, if context is in use. */
	Check int32
	Version

	alloc.Arena

	Socket     os.Handle
	RemoteAddr string

	RequestBuffer *buffer.Circular

	DateRFC822 []byte

	ResponseBuffer []byte
	ResponsePos    int

	CloseAfterWrite bool
	Closed          bool
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

func (c *Conn) Close() error {
	if c.Closed {
		return nil
	}

	c.RequestBuffer.Free()
	c.Arena.Reset()
	c.Closed = true

	return os.Close(c.Socket)
}

func (c *Conn) ReadRequests(rs []Request) (int, error) {
	t := trace.Begin("")

	buf := c.RequestBuffer.RemainingSlice()
	if (c.Closed) || (len(rs) == 0) {
		trace.End(t)
		return 0, nil
	} else if len(buf) == 0 {
		rs[0].Error = RequestEntityTooLarge("no space left in buffer")
		trace.End(t)
		return 1, nil
	}

	n, err := os.Read(c.Socket, buf)
	if err != nil {
		trace.End(t)
		return -1, err
	}
	c.RequestBuffer.Produce(int(n))

	trace.End(t)
	return ParseRequests(c, rs)
}

func (c *Conn) WriteResponses(ws []Response) (int, error) {
	t := trace.Begin("")

	var err error

	if c.Closed {
		panic("write to a closed connection")
	}

	FillResponses(c, ws)

	n, err := os.Write(c.Socket, c.ResponseBuffer[c.ResponsePos:])
	if err != nil {
		trace.End(t)
		return -1, err
	}
	c.ResponsePos += int(n)

	if c.ResponsePos == len(c.ResponseBuffer) {
		c.ResponseBuffer = c.ResponseBuffer[:0]
		c.ResponsePos = 0
		if c.CloseAfterWrite {
			err = c.Close()
		}
	}

	trace.End(t)
	return len(ws), err
}

//go:norace
func (c *Conn) Pointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(c)) | uintptr(c.Check))
}

func GetConnFromPointer(ptr unsafe.Pointer) (*Conn, bool) {
	if ptr == nil {
		return nil, false
	}

	check := uintptr(ptr) & 0x1
	c := (*Conn)(unsafe.Pointer(uintptr(ptr) - check))

	return c, c.Check == int32(check)
}
