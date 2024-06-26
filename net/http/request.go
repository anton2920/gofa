package http

import (
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/net/url"
	"github.com/anton2920/gofa/strings"
)

type Request struct {
	Arena arena.Arena

	RemoteAddr string

	Method string
	URL    url.URL
	Proto  string

	Headers []string
	Body    []byte

	Form url.Values
}

func (r *Request) Cookie(name string) string {
	for i := 0; i < len(r.Headers); i++ {
		header := r.Headers[i]
		if strings.StartsWith(header, "Cookie: ") {
			cookie := header[len("Cookie: "):]
			if strings.StartsWith(cookie, name) {
				cookie = cookie[len(name):]
				if cookie[0] != '=' {
					return ""
				}
				return cookie[1:]
			}

		}
	}

	return ""
}

func (r *Request) ParseForm() error {
	var err error

	if len(r.Form) != 0 {
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

		keyBuffer := r.Arena.NewSlice(len(key))
		n, ok := url.QueryDecode(keyBuffer, key)
		if !ok {
			if err == nil {
				err = errors.New("invalid key")
			}
			continue
		}
		key = unsafe.String(unsafe.SliceData(keyBuffer), n)

		valueBuffer := r.Arena.NewSlice(len(value))
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

	return err
}

func (r *Request) Reset() {
	r.Headers = r.Headers[:0]
	r.Body = r.Body[:0]
	r.Form = r.Form[:0]
	r.Arena.Reset()
}
