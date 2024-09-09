package http

import (
	"github.com/anton2920/gofa/mime/multipart"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
	"github.com/anton2920/gofa/util"
)

type Request struct {
	Method string
	URL    url.URL
	Proto  string

	Headers Headers
	Body    []byte

	Form  url.Values
	Files multipart.Files
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

func (r *Request) ParseForm() error {
	t := trace.Begin("")

	if len(r.Form.Keys) != 0 {
		trace.End(t)
		return nil
	}

	query := util.Slice2String(r.Body)
	err := url.ParseQuery(&r.Form, query)

	trace.End(t)
	return err
}

func (r *Request) ParseMultipartForm() error {
	t := trace.Begin("")

	if (len(r.Form.Keys) != 0) || (len(r.Files.Keys) != 0) {
		trace.End(t)
		return nil
	}

	err := multipart.ParseFormData(r.Headers.Get("Content-Type"), &r.Form, &r.Files, r.Body)

	trace.End(t)
	return err
}

func (r *Request) Reset() {
	r.URL.Query.Reset()
	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()
	r.Files.Reset()
}
