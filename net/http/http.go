package http

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/buffer"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/event"
	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/net/tcp"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

func Accept(l int32, bufferSize int) (*Context, error) {
	var addr tcp.SockAddrIn
	var addrLen uint32 = uint32(unsafe.Sizeof(addr))

	c, err := syscall.Accept(l, (*syscall.Sockaddr)(unsafe.Pointer(&addr)), &addrLen)
	if err != nil {
		return nil, fmt.Errorf("failed to accept incoming connection: %w", err)
	}

	ctx, err := NewContext(c, addr, bufferSize)
	if err != nil {
		syscall.Close(c)
		return nil, fmt.Errorf("failed to create new  context: %w", err)
	}

	return ctx, nil
}

//go:norace
func AddClientToQueue(q *event.Queue, ctx *Context, request event.Request, trigger event.Trigger) error {
	/* TODO(anton2920): switch to pinning inside platform methods. */
	q.Pinner.Pin(ctx)
	ctx.EventQueue = q
	return q.AddSocket(ctx.Connection, request, trigger, ctx.Pointer())
}

/* ReadRequests reads data from socket and parses  requests. Returns the number of requests parsed. */
func ReadRequests(ctx *Context, rs []Request) (int, error) {
	usesQ := ctx.EventQueue != nil
	parser := &ctx.RequestParser
	rBuf := &ctx.RequestBuffer

	if (!usesQ) || (ctx.RequestPendingBytes > 0) {
		if rBuf.RemainingSpace() == 0 {
			return 0, errors.New("no space left in the buffer")
		}
		n, err := syscall.Read(ctx.Connection, rBuf.RemainingSlice())
		if err != nil {
			return 0, err
		}
		rBuf.Produce(int(n))
		ctx.RequestPendingBytes = max(0, ctx.RequestPendingBytes-int(n))
	}

	var i int
	for i = 0; i < len(rs); i++ {
		r := &rs[i]

		r.RemoteAddr = ctx.ClientAddress
		r.Reset()

		n, err := parser.Parse(rBuf.UnconsumedString(), r)
		if err != nil {
			return i, err
		}
		if n == 0 {
			break
		}
		rBuf.Consume(n)
	}

	if (usesQ) && ((ctx.RequestPendingBytes > 0) || (i == len(rs))) {
		ctx.EventQueue.AppendEvent(event.Event{Type: event.Read, Identifier: ctx.Connection, Available: ctx.RequestPendingBytes, UserData: unsafe.Pointer(ctx)})
	}

	return i, nil
}

func ContentTypeHTML(bodies []syscall.Iovec) bool {
	if len(bodies) == 0 {
		return false
	}
	return bodies[0] == html.Header
}

/* WriteResponses generates  responses and writes them on wire. Returns the number of processed responses. */
func WriteResponses(ctx *Context, ws []Response) (int, error) {
	/* TODO(anton2920): remove this. */
	ctx.DateRFC822 = []byte("Thu, 09 May 2024 16:30:39 +0300")
	dateBuf := ctx.DateRFC822

	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("HTTP/1.1"), syscall.Iovec(" "), syscall.Iovec(Status2String[w.StatusCode]), syscall.Iovec(" "), syscall.Iovec(Status2Reason[w.StatusCode]), syscall.Iovec("\r\n"))

		if !w.Headers.OmitDate {
			if dateBuf == nil {
				dateBuf := make([]byte, 31)

				var tp syscall.Timespec
				syscall.ClockGettime(syscall.CLOCK_REALTIME, &tp)
				time.PutTmRFC822(dateBuf, time.ToTm(int(tp.Sec)))
			}

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Date: "), syscall.IovecForByteSlice(dateBuf), syscall.Iovec("\r\n"))
		}

		if !w.Headers.OmitServer {
			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Server: gofa/http\r\n"))
		}

		if !w.Headers.OmitContentType {
			if ContentTypeHTML(w.Bodies) {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/html; charset=\"UTF-8\"\r\n"))
			} else {
				ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Type: text/plain; charset=\"UTF-8\"\r\n"))
			}
		}

		if !w.Headers.OmitContentLength {
			var length int
			for i := 0; i < len(w.Bodies); i++ {
				length += int(len(w.Bodies[i]))
			}

			lengthBuf := w.Arena.NewSlice(20)
			n := slices.PutInt(lengthBuf, length)

			ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("Content-Length: "), syscall.IovecForByteSlice(lengthBuf[:n]), syscall.Iovec("\r\n"))
		}

		ctx.ResponseIovs = append(ctx.ResponseIovs, w.Headers.Values...)
		ctx.ResponseIovs = append(ctx.ResponseIovs, syscall.Iovec("\r\n"))
		ctx.ResponseIovs = append(ctx.ResponseIovs, w.Bodies...)
		w.Reset()
	}

	/* NOTE(anton2920): IOV_MAX is 1024, so F**CK ME for not sending large pipelines with one syscall!!! */
	for len(ctx.ResponseIovs[ctx.ResponsePos:]) > 0 {
		end := min(len(ctx.ResponseIovs[ctx.ResponsePos:]), syscall.IOV_MAX)
		n, err := syscall.Writev(ctx.Connection, ctx.ResponseIovs[ctx.ResponsePos:ctx.ResponsePos+end])
		if err != nil {
			return 0, err
		}

		prevPos := ctx.ResponsePos
		for (ctx.ResponsePos < len(ctx.ResponseIovs)) && (n >= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))) {
			n -= int64(len(ctx.ResponseIovs[ctx.ResponsePos]))
			ctx.ResponsePos++
		}
		if ctx.ResponsePos == len(ctx.ResponseIovs) {
			ctx.ResponseIovs = ctx.ResponseIovs[:0]
			ctx.ResponsePos = 0
		} else if ctx.ResponsePos-prevPos < end {
			log.Panicf("Written %d iovs out of %d", ctx.ResponsePos-prevPos, end)
			ctx.ResponseIovs[ctx.ResponsePos] = ctx.ResponseIovs[ctx.ResponsePos][n:]
			break

			/* TODO(anton2920): as an option, gather buffers manually into some local arena. */
		}
	}

	return len(ws), nil
}

func Close(ctx *Context) error {
	ctx.Reset()
	buffer.FreeCircular(&ctx.RequestBuffer)
	return syscall.Close(ctx.Connection)
}
