package http1

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/log"
	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/strings"
)

type State int

const (
	StateRequestLine State = iota
	StateHeader
	StateBody

	StateUnknown
	StateDone
)

func ParseRequestsUnsafeEx(buffer []byte, consumed *int, rs []http.Request, remoteAddr string) (int, error) {
	var err error

	request := unsafe.String(unsafe.SliceData(buffer), len(buffer))
	pos := *consumed

	var i int
	for i = 0; i < len(rs); i++ {
		var contentLength int
		var state State

		r := &rs[i]
		r.RemoteAddr = remoteAddr
		r.Reset()

		for state != StateDone {
			switch state {
			default:
				log.Panicf("Unknown parser state %d", state)
			case StateUnknown:
				if len(request[pos:]) < 2 {
					return i, nil
				}
				if request[pos:pos+2] == "\r\n" {
					pos += len("\r\n")

					if contentLength != 0 {
						state = StateBody
					} else {
						state = StateDone
					}
				} else {
					state = StateHeader
				}
			case StateRequestLine:
				lineEnd := strings.FindChar(request[pos:], '\r')
				if lineEnd == -1 {
					return i, nil
				}

				sp := strings.FindChar(request[pos:pos+lineEnd], ' ')
				if sp == -1 {
					return i, http.BadRequest("expected method, found %q", request[pos:])
				}
				r.Method = request[pos : pos+sp]
				pos += len(r.Method) + 1
				lineEnd -= len(r.Method) + 1

				uriEnd := strings.FindChar(request[pos:pos+lineEnd], ' ')
				if uriEnd == -1 {
					return i, http.BadRequest("expected space after URI, found %q", request[pos:pos+lineEnd])
				}

				queryStart := strings.FindChar(request[pos:pos+uriEnd], '?')
				if queryStart != -1 {
					r.URL.Path = request[pos : pos+queryStart]
					r.URL.Query = request[pos+queryStart+1 : pos+uriEnd]
				} else {
					r.URL.Path = request[pos : pos+uriEnd]
					r.URL.Query = ""
				}
				pos += len(r.URL.Path) + len(r.URL.Query) + 1
				lineEnd -= len(r.URL.Path) + len(r.URL.Query) + 1

				if request[pos:pos+len("HTTP/")] != "HTTP/" {
					return i, http.BadRequest("expected version prefix, found %q", request[pos:pos+lineEnd])
				}
				r.Proto = request[pos : pos+lineEnd]

				pos += len(r.Proto) + len("\r\n")
				state = StateUnknown
			case StateHeader:
				lineEnd := strings.FindChar(request[pos:], '\r')
				if lineEnd == -1 {
					return i, nil
				}

				header := request[pos : pos+lineEnd]
				colon := strings.FindChar(header, ':')
				if colon == -1 {
					return i, http.BadRequest("expected HTTP header, got %q", header)
				}

				key := header[:colon]
				value := header[colon+2:]
				r.Headers.Add(key, value)

				if key == "Content-Length" {
					contentLength, err = strconv.Atoi(value)
					if err != nil {
						return i, http.BadRequest("failed to parse Content-Length value: %v", err)
					}
				}

				pos += len(header) + len("\r\n")
				state = StateUnknown
			case StateBody:
				if len(request[pos:]) < contentLength {
					return i, nil
				}

				r.Body = unsafe.Slice(unsafe.StringData(request[pos:]), contentLength)
				pos += len(r.Body)
				state = StateDone
			}
		}

		*consumed = pos
	}

	return i, nil
}

/* ParseRequestsUnsafe fills slice of requests with data from (*http.Context).RequestBuffer. Data in buffer must live for as long as requests are needed. */
func ParseRequestsUnsafe(ctx *http.Context, rs []http.Request) (int, error) {
	rBuf := &ctx.RequestBuffer
	var pos int

	n, err := ParseRequestsUnsafeEx(rBuf.UnconsumedSlice(), &pos, rs, ctx.ClientAddress)
	rBuf.Consume(pos)
	return n, err
}
