package http

import (
	"strconv"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type Response struct {
	Arena alloc.Arena

	StatusCode Status
	Headers    Headers
	Body       []byte
}

func (w *Response) DelCookie(name string) {
	t := trace.Begin("")

	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.Headers.Set("Set-Cookie", bytes.AsString(cookie[:n]))

	trace.End(t)
}

func (w *Response) SetCookie(name, value string, expiry int64) {
	t := trace.Begin("")

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

	trace.End(t)
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int64) {
	t := trace.Begin("")

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

	trace.End(t)
}

func (w *Response) Redirect(path string, code Status) {
	t := trace.Begin("")

	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.Headers.Set("Location", bytes.AsString(pathBuf))
	w.Body = w.Body[:0]
	w.StatusCode = code

	trace.End(t)
}

func (w *Response) RedirectID(prefix string, id database.ID, code Status) {
	t := trace.Begin("")

	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], int(id))

	w.Headers.Set("Location", bytes.AsString(buffer[:n]))
	w.Body = w.Body[:0]
	w.StatusCode = code

	trace.End(t)
}

func (w *Response) Write(b []byte) (int, error) {
	w.Body = append(w.Body, b...)
	return len(b), nil
}

/* WriteHTML writes to w the escaped html. equivalent of the plain text data b. */
func (w *Response) WriteHTML(b []byte) {
	w.WriteHTMLString(bytes.AsString(b))
}

func (w *Response) WriteFloat64(f float64) {
	w.Body = strconv.AppendFloat(w.Body, f, 'f', -1, 64)
}

func (w *Response) WriteInt(i int) {
	buffer := make([]byte, 20)
	n := slices.PutInt(buffer, i)
	w.Body = append(w.Body, buffer[:n]...)
}

func (w *Response) WriteID(id database.ID) {
	buffer := make([]byte, 20)
	n := slices.PutInt(buffer, int(id))
	w.Body = append(w.Body, buffer[:n]...)
}

func (w *Response) WriteString(s string) (int, error) {
	w.Body = append(w.Body, s...)
	return len(s), nil
}

func (w *Response) WriteHTMLString(s string) {
	t := trace.Begin("")

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

	trace.End(t)
}

func (w *Response) Reset() {
	w.StatusCode = StatusOK
	w.Headers.Reset()
	w.Body = w.Body[:0]
	w.Arena.Reset()
}

func FillResponses(c *Conn, ws []Response) {
	t := trace.Begin("")

	for i := 0; i < len(ws); i++ {
		w := &ws[i]

		c.ResponseBuffer = append(c.ResponseBuffer, StatusLines[c.Version][w.StatusCode]...)

		if !w.Headers.Has("Date") {
			dateBuf := c.DateRFC822
			if dateBuf == nil {
				dateBuf = make([]byte, time.RFC822Len)
				time.PutTmRFC822(dateBuf, time.ToTm(time.Now()))
			}
			c.ResponseBuffer = append(c.ResponseBuffer, "Date: "...)
			c.ResponseBuffer = append(c.ResponseBuffer, dateBuf...)
			c.ResponseBuffer = append(c.ResponseBuffer, "\r\n"...)
		}

		if !w.Headers.Has("Server") {
			c.ResponseBuffer = append(c.ResponseBuffer, "Server: gofa/http\r\n"...)
		}

		if !w.Headers.Has("Content-Type") {
			c.ResponseBuffer = append(c.ResponseBuffer, "Content-Type: text/plain; charset=\"UTF-8\"\r\n"...)
		}

		if !w.Headers.Has("Content-Length") {
			lengthBuf := make([]byte, ints.Bufsize)
			n := slices.PutInt(lengthBuf, len(w.Body))

			c.ResponseBuffer = append(c.ResponseBuffer, "Content-Length: "...)
			c.ResponseBuffer = append(c.ResponseBuffer, lengthBuf[:n]...)
			c.ResponseBuffer = append(c.ResponseBuffer, "\r\n"...)
		}

		for i := 0; i < len(w.Headers.Keys); i++ {
			key := w.Headers.Keys[i]
			c.ResponseBuffer = append(c.ResponseBuffer, key...)
			c.ResponseBuffer = append(c.ResponseBuffer, ": "...)
			for j := 0; j < len(w.Headers.Values[i]); j++ {
				value := w.Headers.Values[i][j]
				if j > 0 {
					c.ResponseBuffer = append(c.ResponseBuffer, ","...)
				}
				c.ResponseBuffer = append(c.ResponseBuffer, value...)
			}
			c.ResponseBuffer = append(c.ResponseBuffer, "\r\n"...)
		}

		c.ResponseBuffer = append(c.ResponseBuffer, "\r\n"...)
		c.ResponseBuffer = append(c.ResponseBuffer, w.Body...)
		w.Reset()
	}

	trace.End(t)
}
