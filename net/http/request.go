package http

import (
	"unsafe"

	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Request struct {
	Method string
	URL    url.URL
	Proto  string

	Headers Headers
	Body    []byte

	Form url.Values
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
			return cookie[1:]
		}
	}

	trace.End(t)
	return ""
}

func (r *Request) ParseForm() error {
	t := trace.Begin("")

	if len(r.Form.Keys) != 0 {
		trace.End(t)
		return nil
	}

	query := unsafe.String(unsafe.SliceData(r.Body), len(r.Body))
	err := url.ParseQuery(&r.Form, query)

	trace.End(t)
	return err
}

func (r *Request) Reset() {
	r.URL.Query.Reset()
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
}
