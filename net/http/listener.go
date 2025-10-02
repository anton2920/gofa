package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/floats"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/os"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/trace"
)

type Listener struct {
	*ConnPool

	Socket os.Handle
}

type ListenerOptions struct {
	Backlog               int
	ConcurrentConnections int

	MaxVersion float32
}

func MergeListenerOptions(opts ...ListenerOptions) ListenerOptions {
	t := trace.Begin("")

	var result ListenerOptions

	for i := 0; i < len(opts); i++ {
		opt := &opts[i]

		ints.Replace(&result.Backlog, opt.Backlog)
		ints.Replace(&result.ConcurrentConnections, opt.ConcurrentConnections)

		floats.Replace32(&result.MaxVersion, opt.MaxVersion)
	}

	trace.End(t)
	return result
}

func Listen(addr string, opts ...ListenerOptions) (*Listener, error) {
	t := trace.Begin("")

	var l Listener
	var err error

	opt := MergeListenerOptions(opts...)

	l.Socket, err = tcp.Listen(addr, ints.Or(opt.Backlog, 128))
	if err != nil {
		trace.End(t)
		return nil, fmt.Errorf("failed to listen on addr %q: %v", err)
	}
	if opt.MaxVersion > 2.0 {
		trace.End(t)
		panic("HTTP/2+ is not supported")
	}

	l.ConnPool = NewConnPool(ints.Or(opt.ConcurrentConnections, 16*1024))

	trace.End(t)
	return &l, nil
}

/* TODO(anton2920): remove syscall references. */
func (l *Listener) Accept(opts ...ConnOptions) (*Conn, error) {
	t := trace.Begin("")

	var addr syscall.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	opt := MergeConnOptions(opts...)

	sock, err := syscall.Accept(int32(l.Socket), (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		trace.End(t)
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	rb, err := buffer.NewCircular(ints.Or(opt.RequestBufferSize, os.PageSize))
	if err != nil {
		syscall.Close(sock)
		trace.End(t)
		return nil, fmt.Errorf("failed to create new request buffer: %w", err)
	}

	c, err := l.ConnPool.Get()
	if err != nil {
		syscall.Close(sock)
		trace.End(t)
		panic("handle too many connections")
	}
	c.ConnPool = l.ConnPool

	c.Socket = os.Handle(sock)
	c.RequestBuffer = rb

	buffer := c.Arena.NewSlice(21)
	n := tcp.PutAddress(buffer, addr.Addr, addr.Port)
	c.RemoteAddr = string(buffer[:n])

	trace.End(t)
	return c, err
}

func (l *Listener) Close() error {
	t := trace.Begin("")

	err := os.Close(l.Socket)

	trace.End(t)
	return err
}
