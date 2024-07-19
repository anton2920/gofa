package http1

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/net/http"
	"github.com/anton2920/gofa/strings"
)

func ParseRequestsUnsafeEx(buffer []byte, consumed *int, rs []http.Request, remoteAddr string) (int, error) {
	request := unsafe.String(unsafe.SliceData(buffer), len(buffer))
	pos := *consumed

	var i int
	for i = 0; i < len(rs); i++ {
		r := &rs[i]
		r.Reset()

		/* Parsing request line. */
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

		/* Parsing headers. */
		for {
			lineEnd := strings.FindChar(request[pos:], '\r')
			if lineEnd == -1 {
				return i, nil
			} else if lineEnd == 0 {
				pos += len("\r\n")
				break
			}

			header := request[pos : pos+lineEnd]
			colon := strings.FindChar(header, ':')
			if colon == -1 {
				return i, http.BadRequest("expected HTTP header, got %q", header)
			}

			key := header[:colon]
			value := header[colon+2:]
			r.Headers.Add(key, value)

			pos += len(header) + len("\r\n")
		}

		/* Parsing body. */
		if r.Headers.Has("Content-Length") {
			contentLength, err := strconv.Atoi(r.Headers.Get("Content-Length"))
			if (err != nil) || (contentLength < 0) {
				return i, http.BadRequest("invalid Content-Length value: %q", r.Headers.Get("Content-Length"))
			}

			if len(request[pos:]) < contentLength {
				return i, nil
			}

			r.Body = unsafe.Slice(unsafe.StringData(request[pos:]), contentLength)
			pos += len(r.Body)
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
