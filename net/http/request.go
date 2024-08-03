package http

import (
	"unsafe"

	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/strings"
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
	p := prof.Begin("")

	cookies := r.Headers.GetMany("Cookie")
	for i := 0; i < len(cookies); i++ {
		cookie := cookies[i]
		if strings.StartsWith(cookie, name) {
			cookie = cookie[len(name):]
			if cookie[0] != '=' {
				prof.End(p)
				return ""
			}
			return cookie[1:]
		}
	}

	prof.End(p)
	return ""
}

func (r *Request) ParseForm() error {
	p := prof.Begin("")

	if len(r.Form.Keys) != 0 {
		prof.End(p)
		return nil
	}

	query := unsafe.String(unsafe.SliceData(r.Body), len(r.Body))
	err := url.ParseQuery(&r.Form, query)

	prof.End(p)
	return err
}

func (r *Request) Reset() {
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
}
