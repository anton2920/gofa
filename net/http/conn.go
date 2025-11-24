package http

import (
	"io"
	"unsafe"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/pointers"
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
	*ConnPool
	alloc.Arena

	Version
	Socket     os.Handle
	RemoteAddr string

	RequestBuffer *buffer.Circular

	ResponseBuffer []byte
	ResponsePos    int

	CloseAfterWrite bool
	Closed          bool

	/* NOTE(anton2920): Check must be the same as the last pointer's bit, if context is in use. */
	Check bool
}

type ConnOptions struct {
	RequestBufferSize int
}

func MergeConnOptions(opts ...ConnOptions) ConnOptions {
	t := trace.Begin("")

	var result ConnOptions

	for i := 0; i < len(opts); i++ {
		opt := &opts[i]

		ints.Replace(&result.RequestBufferSize, opt.RequestBufferSize)
	}

	trace.End(t)
	return result
}

func (c *Conn) Close() error {
	t := trace.Begin("")

	if c.Closed {
		trace.End(t)
		return nil
	}

	c.RequestBuffer.Free()
	c.Arena.Reset()
	c.Closed = true

	err := os.Close(c.Socket)
	c.Check = !c.Check
	c.ConnPool.Put(c)

	trace.End(t)
	return err
}

func (c *Conn) ReadRequests(rs []Request) (int, error) {
	t := trace.Begin("")

	buf := c.RequestBuffer.RemainingSlice()
	if (c.Closed) || (len(rs) == 0) {
		trace.End(t)
		return 0, nil
	} else if (len(buf) == 0) && (len(rs) > 0) {
		rs[0].Error = RequestEntityTooLarge("no space left in buffer")
		trace.End(t)
		return 1, nil
	}

	n, err := os.Read(c.Socket, buf)
	if err != nil {
		trace.End(t)
		return -1, err
	} else if n == 0 {
		trace.End(t)
		return 0, io.EOF
	}
	c.RequestBuffer.Produce(int(n))

	nrs := ParseRequests(c, rs)

	trace.End(t)
	return nrs, nil
}

func (c *Conn) WriteResponses(ws []Response) (int, error) {
	t := trace.Begin("")

	var err error

	if c.Closed {
		trace.End(t)
		panic("write to a closed connection")
	}
	FillResponses(c, ws)

	if len(c.ResponseBuffer[c.ResponsePos:]) > 0 {
		n, err := os.Write(c.Socket, c.ResponseBuffer[c.ResponsePos:])
		if err != nil {
			trace.End(t)
			return int(n), err
		}
		c.ResponsePos += int(n)

		if c.ResponsePos == len(c.ResponseBuffer) {
			c.ResponseBuffer = c.ResponseBuffer[:0]
			c.ResponsePos = 0
			if c.CloseAfterWrite {
				err = c.Close()
			}
		}
	}

	trace.End(t)
	return len(ws), err
}

func (c *Conn) Pointer() unsafe.Pointer {
	return pointers.Add(unsafe.Pointer(c), bools.ToInt(c.Check))
}

func ConnFromPointer(ptr unsafe.Pointer) (*Conn, bool) {
	if ptr == nil {
		return nil, false
	}

	check := uintptr(ptr) & 0x1
	c := (*Conn)(unsafe.Pointer(uintptr(ptr) - check))

	return c, c.Check == ints.ToBool(int(check))
}

func RequestBufferSize(size int) ConnOptions {
	return ConnOptions{RequestBufferSize: size}
}
