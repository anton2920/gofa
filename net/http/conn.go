package http

import (
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/syscall"
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

type Conn struct {
	*Listener
	Version

	/* NOTE(anton2920): Check must be the same as the last pointer's bit, if context is in use. */
	Check int32

	Socket  os.Handle
	Address string

	RequestBuffer *buffer.Circular

	ResponseBuffer []byte
	ResponsePos    int

	Closed bool
}

func (c *Conn) Close() error {
	if c.Closed {
		return nil
	}

	c.RequestBuffer.Free()
	c.Closed = true

	return os.Close(c.Socket)
}

func (c *Conn) Read(rs []Request) (int, error) {
	rBuf := c.RequestBuffer

	if len(buf) == 0 {
		rs[0].Error = NoSpaceLeft
		return 1, nil
	}

	n, err := os.Read(c.Socket, rBuf.RemainingSlice())
	if err != nil {
		return -1, err
	}

	return ParseRequests(rBuf, rs)
}

func (c *Conn) Write(ws []Response) (int, error) {
	wBuf := &c.ResponseBuffer

	FillResponses(wBuf, ws)

	n, err := os.Write(c.Socket, wBuf)
	if err != nil {
		return -1, err
	}

}

//go:norace
func (c *Conn) Pointer() unsafe.Pointer {
	return unsafe.Pointer(uintptr(unsafe.Pointer(c)) | uintptr(c.Check))
}

//go:norace
func InitConn(c *Conn, sock int32, addr syscall.SockAddrIn, rb *buffer.Circular) {
	c.Socket = sock
	c.RequestBuffer = rb

	buffer := make([]byte, 21)
	n := tcp.PutAddress(buffer, addr.Addr, addr.Port)
	c.ClientAddress = string(buffer[:n])
}

func GetConnFromPointer(ptr unsafe.Pointer) (*Conn, bool) {
	if ptr == nil {
		return nil, false
	}

	check := uintptr(ptr) & 0x1
	c := (*Context)(unsafe.Pointer(uintptr(ptr) - check))

	return c, c.Check == int32(check)
}
