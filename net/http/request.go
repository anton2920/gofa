package http

import (
	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Request struct {
	session.Session

	Method string
	URL    url.URL

	RemoteAddr string

	Proto      string
	ProtoMajor int
	ProtoMinor int

	Headers Headers
	Body    []byte

	Form  url.Values
	Files multipart.Files

	Arena alloc.Arena
	Error error
}

func (r *Request) Cookie(name string) string {
	t := trace.Begin("")

	cookies := r.Headers.GetMany("Cookie")
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if strings.StartsWith(cookie, name) {
			cookie = cookie[len(name):]
			if cookie[0] != '=' {
				trace.End(t)
				return ""
			}
			trace.End(t)
			return cookie[1:]
		}
	}

	trace.End(t)
	return ""
}

func (r *Request) Reset() {
	//t := trace.Begin("")

	r.URL.Query.Reset()
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
	r.Files.Reset()
	r.Arena.Reset()
	r.Error = nil

	//trace.End(t)
}

func ParseRequestsV1(c *Conn, rs []Request) int {
	var consumed int
	var pos int
	var i int

	rBuf := &c.RequestBuffer
	remoteAddr := c.RemoteAddr()
	requestBytes := rBuf.UnconsumedSlice()
	request := bytes.AsString(requestBytes)

forRequests:
	for i = 0; i < len(rs); i++ {
		r := &rs[i]
		r.Reset()
		r.RemoteAddr = remoteAddr

		/* Parsing request line. */
		lineEnd := strings.FindChar(request[pos:], '\r')
		if lineEnd == -1 {
			break
		}

		sp := strings.FindChar(request[pos:pos+lineEnd], ' ')
		if sp == -1 {
			r.Error = BadRequest("expected method, found %q", request[pos:])
			rBuf.Reset()
			return i + 1
		}
		r.Method = r.Arena.CopyString(request[pos : pos+sp])
		pos += len(r.Method) + 1
		lineEnd -= len(r.Method) + 1

		uriEnd := strings.FindChar(request[pos:pos+lineEnd], ' ')
		if uriEnd == -1 {
			r.Error = BadRequest("expected space after URI, found %q", request[pos:pos+lineEnd])
			rBuf.Reset()
			return i + 1
		}

		queryBegin := strings.FindChar(request[pos:pos+uriEnd], '?')
		if queryBegin != -1 {
			r.URL.Path = url.Path(r.Arena.CopyString(request[pos : pos+queryBegin]))
			r.URL.RawQuery = r.Arena.CopyString(request[pos+queryBegin+1 : pos+uriEnd])
			pos += len(r.URL.Path) + len(r.URL.RawQuery) + 2
			lineEnd -= len(r.URL.Path) + len(r.URL.RawQuery) + 2
		} else {
			r.URL.Path = url.Path(r.Arena.CopyString(request[pos : pos+uriEnd]))
			r.URL.RawQuery = ""
			pos += len(r.URL.Path) + 1
			lineEnd -= len(r.URL.Path) + 1
		}

		const versionPrefix = "HTTP/"
		if request[pos:pos+len(versionPrefix)] != versionPrefix {
			r.Error = BadRequest("expected protocol, found %q", request[pos:pos+lineEnd])
			rBuf.Reset()
			return i + 1
		}
		r.Proto = request[pos : pos+lineEnd]
		/*
			switch request[pos+len(versionPrefix) : pos+lineEnd] {
			case "1.1":
				r.Proto = "HTTP/1.1"
				r.ProtoMajor = 1
				r.ProtoMinor = 1
				c.Version = Version11
			case "1.0":
				r.Proto = "HTTP/1.0"
				r.ProtoMajor = 1
				r.ProtoMinor = 0
				c.Version = Version10
				c.CloseAfterWrite = true
			case "0.9":
				r.Proto = "HTTP/0.9"
				r.ProtoMajor = 0
				r.ProtoMinor = 9
				c.Version = Version09
				c.CloseAfterWrite = true
			default:
				r.Error = BadRequest("invalid protocol %q", request[pos:pos+lineEnd])
				rBuf.Reset()
				return i + 1
			}
		*/
		pos += len(r.Proto) + len("\r\n")

		/* Parsing headers. */
		for {
			lineEnd := strings.FindChar(request[pos:], '\r')
			if lineEnd == -1 {
				break forRequests
			} else if lineEnd == 0 {
				pos += len("\r\n")
				break
			}

			header := request[pos : pos+lineEnd]
			colon := strings.FindChar(header, ':')
			if colon == -1 {
				r.Error = BadRequest("expected HTTP header, got %q", header)
				rBuf.Reset()
				return i + 1
			}

			key := r.Arena.CopyString(header[:colon])
			value := r.Arena.CopyString(header[colon+2:])
			r.Headers.Add(key, value)

			pos += len(header) + len("\r\n")
		}

		/* Parsing body. */
		/* TODO(anton2920): add support for 'Transfer-Encoding: chunked'. */
		if r.Headers.Has("Content-Length") {
			contentLength, err := r.Headers.GetInt("Content-Length")
			if (err != nil) || (contentLength < 0) {
				r.Error = BadRequest("invalid Content-Length value: %q", r.Headers.Get("Content-Length"))
				rBuf.Reset()
				return i + 1
			}

			if len(request[pos:]) < contentLength {
				break
			}

			r.Body = r.Arena.Copy(requestBytes[pos : pos+contentLength])
			pos += len(r.Body)
		}

		consumed = pos
	}

	rBuf.Consume(consumed)
	return i
}

func ParseRequests(c *Conn, rs []Request) int {
	t := trace.Begin("")

	var n int

	if (c.Error != nil) && (len(rs) > 0) {
		c.RequestBuffer.Reset()
		rs[0].Error = c.Error
		trace.End(t)
		return 1
	}

	/* TODO(anton2920): uncomment once support for HTTP >=2 is added. */
	//switch c.Version {
	//case VersionNone, Version09, Version10, Version11:
	n = ParseRequestsV1(c, rs)
	//default:
	//	trace.End(t)
	//	panic("unsupported version")
	//}

	trace.End(t)
	return n
}
