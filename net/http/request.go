package http

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/context"
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/session"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace/trace_"
)

type Request struct {
	Method string
	URL    url.URL

	Proto      string
	ProtoMajor int
	ProtoMinor int

	Headers Headers
	Body    []byte

	Form  url.Values
	Files multipart.Files

	session.Session
}

func (r *Request) Cookie(name string) string {
	t := trace_.Begin("")

	cookies := r.Headers.GetMany("Cookie")
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if strings.StartsWith(cookie, name) {
			cookie = cookie[len(name):]
			if !strings.StartsWith(cookie, "=") {
				trace_.End(t)
				return ""
			}
			trace_.End(t)
			return cookie[1:]
		}
	}

	trace_.End(t)
	return ""
}

func (r *Request) Reset() {
	r.URL.Query.Reset()
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
	r.Files.Reset()
}

func ParseRequestsV1(ctx *context.Context, request []byte, rs []Request) (n int, nbytes int, ok bool) {
	var consumed int
	var pos int
	var i int

forRequests:
	for i = 0; i < len(rs); i++ {
		r := &rs[i]
		r.Reset()

		/* Parsing request line. */
		/* TODO(anton2920): search for line ending properly. */
		lineEnd := bytes.FindByte(request[pos:], '\r')
		if lineEnd == -1 {
			break
		}

		sp := bytes.FindByte(request[pos:pos+lineEnd], ' ')
		if sp == -1 {
			ctx.NewError().S("expected method, found ").Q(bytes.AsString(request[pos:]))
			return i + 1, consumed, false
		}
		r.Method = bytes.AsString(request[pos : pos+sp])
		pos += len(r.Method) + 1
		lineEnd -= len(r.Method) + 1

		uriEnd := bytes.FindByte(request[pos:pos+lineEnd], ' ')
		if uriEnd == -1 {
			ctx.NewError().S("expected space after URI, found ").Q(bytes.AsString(request[pos : pos+lineEnd]))
			return i + 1, consumed, false
		}

		queryBegin := bytes.FindByte(request[pos:pos+uriEnd], '?')
		if queryBegin != -1 {
			r.URL.Path = url.Path(bytes.AsString(request[pos : pos+queryBegin]))
			r.URL.RawQuery = bytes.AsString(request[pos+queryBegin+1 : pos+uriEnd])
			pos += len(r.URL.Path) + len(r.URL.RawQuery) + 2
			lineEnd -= len(r.URL.Path) + len(r.URL.RawQuery) + 2
		} else {
			r.URL.Path = url.Path(bytes.AsString(request[pos : pos+uriEnd]))
			r.URL.RawQuery = ""
			pos += len(r.URL.Path) + 1
			lineEnd -= len(r.URL.Path) + 1
		}

		const versionPrefix = "HTTP/"
		if bytes.AsString(request[pos:pos+len(versionPrefix)]) != versionPrefix {
			ctx.NewError().S("expected protocol, found ").Q(bytes.AsString(request[pos : pos+lineEnd]))
			return i + 1, consumed, false
		}
		r.Proto = bytes.AsString(request[pos : pos+lineEnd])
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
			lineEnd := bytes.FindByte(request[pos:], '\r')
			if lineEnd == -1 {
				break forRequests
			} else if lineEnd == 0 {
				pos += len("\r\n")
				break
			}

			header := bytes.AsString(request[pos : pos+lineEnd])
			colon := strings.FindChar(header, ':')
			if colon == -1 {
				ctx.NewError().S("expected HTTP header, got ").Q(header)
				return i + 1, consumed, false
			}

			key := header[:colon]
			value := header[colon+2:]
			r.Headers.Add(key, value)

			pos += len(header) + len("\r\n")
		}

		/* Parsing body. */
		/* TODO(anton2920): add support for 'Transfer-Encoding: chunked'. */
		if r.Headers.Has("Content-Length") {
			contentLength, err := r.Headers.GetInt("Content-Length")
			if (err != nil) || (contentLength < 0) {
				ctx.NewError().S("invalid Content-Length value: ").Q(r.Headers.Get("Content-Length"))
				return i + 1, consumed, false
			}

			if len(request[pos:]) < contentLength {
				break
			}

			r.Body = request[pos : pos+contentLength]
			pos += len(r.Body)
		}

		consumed = pos
	}

	return i, consumed, true
}
