package http

import (
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
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
	r.URL.Query.Reset()
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
	r.Files.Reset()
	r.Error = nil
}
