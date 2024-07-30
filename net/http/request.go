package http

import (
	"unsafe"

	"github.com/anton2920/gofa/errors"
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

	var err error

	if len(r.Form.Keys) != 0 {
		prof.End(p)
		return nil
	}

	query := unsafe.String(unsafe.SliceData(r.Body), len(r.Body))
	for query != "" {
		var key string
		key, query, _ = strings.Cut(query, "&")
		if strings.FindChar(key, ';') != -1 {
			err = errors.New("invalid semicolon separator in query")
			continue
		}
		if key == "" {
			continue
		}
		key, value, _ := strings.Cut(key, "=")

		keyBuffer := r.Form.Arena.NewSlice(len(key))
		n, ok := url.QueryDecode(keyBuffer, key)
		if !ok {
			if err == nil {
				err = errors.New("invalid key")
			}
			continue
		}
		key = unsafe.String(unsafe.SliceData(keyBuffer), n)

		valueBuffer := r.Form.Arena.NewSlice(len(value))
		n, ok = url.QueryDecode(valueBuffer, value)
		if !ok {
			if err == nil {
				err = errors.New("invalid value")
			}
			continue
		}
		value = unsafe.String(unsafe.SliceData(valueBuffer), n)

		r.Form.Add(key, value)
	}

	prof.End(p)
	return err
}

//go:nosplit
func (r *Request) Reset() {
	p := prof.Begin("")

	r.Headers.Reset()
	r.Body = r.Body[:0]
	r.Form.Reset()

	prof.End(p)
}
