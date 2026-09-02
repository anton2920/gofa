package http

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace/trace_"
)

type Response struct {
	Status  Status
	Headers Headers
	Body    []byte
}

/*
func (w *Response) DelCookie(name string) {
	t := trace_.Begin("")

	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.Headers.Set("Set-Cookie", bytes.AsString(cookie[:n]))

	trace_.End(t)
}

func (w *Response) SetCookie(name, value string, expiry int64) {
	t := trace_.Begin("")

	const secure = "; HttpOnly; Secure; SameSite=Strict"
	const expires = "; Expires="
	const path = "; Path=/"
	const eq = "="

	cookie := w.Arena.NewSlice(len(name) + len(eq) + len(value) + len(path) + len(expires) + time.RFC822Len + len(secure))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], eq)
	n += copy(cookie[n:], value)
	n += copy(cookie[n:], path)
	n += copy(cookie[n:], expires)
	n += time.PutTmRFC822(cookie[n:], time.ToTm(expiry))
	n += copy(cookie[n:], secure)

	w.Headers.Set("Set-Cookie", bytes.AsString(cookie[:n]))

	trace_.End(t)
}

*/ /* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */ /*
func (w *Response) SetCookieUnsafe(name, value string, expiry int64) {
	t := trace_.Begin("")

	const expires = "; Expires="
	const path = "; Path=/"
	const eq = "="

	cookie := w.Arena.NewSlice(len(name) + len(eq) + len(value) + len(path) + len(expires) + time.RFC822Len)

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], eq)
	n += copy(cookie[n:], value)
	n += copy(cookie[n:], path)
	n += copy(cookie[n:], expires)
	n += time.PutTmRFC822(cookie[n:], time.ToTm(expiry))

	w.Headers.Set("Set-Cookie", bytes.AsString(cookie[:n]))

	trace_.End(t)
}


func (w *Response) Redirect(path string, code Status) {
	t := trace_.Begin("")

	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.Headers.Set("Location", bytes.AsString(pathBuf))
	w.Body = w.Body[:0]
	w.Status = code

	trace_.End(t)
}
*/

func (w *Response) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

func (w *Response) WriteHTML(b []byte) {
	w.WriteHTMLString(bytes.AsString(b))
}

func (w *Response) WriteString(s string) (int, error) {
	w.Body = append(w.Body, s...)
	return len(s), nil
}

func (w *Response) WriteHTMLString(s string) {
	last := 0
	for i := 0; i < len(s); i++ {
		var seq string
		switch s[i] {
		case '\000':
			seq = "\uFFFD"
		case '"':
			seq = "&#34;"
		case '\'':
			seq = "&#39;"
		case '&':
			seq = "&amp;"
		case '<':
			seq = "&lt;"
		case '>':
			seq = "&gt;"
		default:
			continue
		}
		w.Body = append(w.Body, s[last:i]...)
		w.Body = append(w.Body, seq...)
		last = i + 1
	}
	w.Body = append(w.Body, s[last:]...)
}

func (w *Response) Reset() {
	w.Status = StatusOK
	w.Headers.Reset()
	w.Body = w.Body[:0]
}

func FillResponsesV1(response []byte, now int64, ws []Response) {
	t := trace_.Begin("")

	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		/* TODO(anton2920): fill in status line. */
		//response = append(response, StatusLines[c.Version][w.Status]...)

		if !w.Headers.Has("Date") {
			dateBuf := make([]byte, time.RFC822Len)
			time.PutTmRFC822(dateBuf, time.ToTm(now))
			response = append(response, "Date: "...)
			response = append(response, dateBuf...)
			response = append(response, "\r\n"...)
		}

		if !w.Headers.Has("Server") {
			response = append(response, "Server: gofa/http\r\n"...)
		}

		if !w.Headers.Has("Content-Type") {
			response = append(response, "Content-Type: text/plain; charset=\"UTF-8\"\r\n"...)
		}

		if !w.Headers.Has("Content-Length") {
			lengthBuf := make([]byte, ints.Bufsize)
			n := slices.PutInt(lengthBuf, len(w.Body))

			response = append(response, "Content-Length: "...)
			response = append(response, lengthBuf[:n]...)
			response = append(response, "\r\n"...)
		}

		for i := 0; i < len(w.Headers.Keys); i++ {
			key := w.Headers.Keys[i]
			response = append(response, key...)
			response = append(response, ": "...)
			for j := 0; j < len(w.Headers.Values[i]); j++ {
				value := w.Headers.Values[i][j]
				if j > 0 {
					response = append(response, ","...)
				}
				response = append(response, value...)
			}
			response = append(response, "\r\n"...)
		}

		response = append(response, "\r\n"...)
		response = append(response, w.Body...)

		//connection := w.Headers.Get("Connection")
		w.Reset()

		//if connection == "close" {
		//	c.CloseAfterWrite = true
		//	break
		//}
	}

	trace_.End(t)
}
