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
)

type Listener struct {
	*ConnPool

	SocketTCP os.Handle
	SocketUDP os.Handle
}

type ListenerParams struct {
	Backlog               int
	ConcurrentConnections int

	MaxVersion float32
}

func MergeParams(params ...ListenerParams) ListenerParams {
	var result ListenerParams

	for i := 0; i < len(params); i++ {
		param := &params[i]

		ints.Replace(&result.Backlog, param.Backlog)
		ints.Replace(&result.ConcurrentConnections, param.ConcurrentConnections)

		floats.Replace32(&result.MaxVersion, param.MaxVersion)
	}

	return result
}

func Listen(addr string, params ...ListenerParams) (*Listener, error) {
	var l Listener
	var err error

	result := MergeParams(params...)

	l.SocketTCP, err = tcp.Listen(addr, ints.Or(result.Backlog, 128))
	if err != nil {
		return nil, fmt.Errorf("failed to listen on addr %q: %v", err)
	}
	if result.MaxVersion == 3.0 {
		/* TODO(anton2920): listen for HTTP/3. */
	}

	l.ConnPool = NewConnPool(ints.Or(result.ConcurrentConnections, 16*1024))

	return &l, nil
}

/* TODO(anton2920): remove syscall references. */
func (l *Listener) Accept() (*Conn, error) {
	var addr syscall.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	rb, err := buffer.NewCircular(bufferSize)
	if err != nil {
		syscall.Close(c)
		return nil, fmt.Errorf("failed to create new request buffer: %w", err)
	}

	ctx, err := CtxPool.Get()
	if err != nil {
		ctx = new(Context)
		err = TooManyClients
	}
	InitContext(ctx, c, addr, rb)

	return ctx, err
}
