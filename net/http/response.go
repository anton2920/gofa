package http

import (
	"strconv"

	"github.com/anton2920/gofa/alloc"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
	"github.com/anton2920/gofa/util"
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

	w.Headers.Set("Set-Cookie", util.Slice2String(cookie[:n]))

	trace.End(t)
}

func (w *Response) SetCookie(name, value string, expiry int) {
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

	w.Headers.Set("Set-Cookie", util.Slice2String(cookie[:n]))

	trace.End(t)
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int) {
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

	w.Headers.Set("Set-Cookie", util.Slice2String(cookie[:n]))

	trace.End(t)
}

func (w *Response) Redirect(path string, code Status) {
	t := trace.Begin("")

	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.Headers.Set("Location", util.Slice2String(pathBuf))
	w.Body = w.Body[:0]
	w.StatusCode = code

	trace.End(t)
}

func (w *Response) RedirectID(prefix string, id database.ID, code Status) {
	t := trace.Begin("")

	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], int(id))

	w.Headers.Set("Location", util.Slice2String(buffer[:n]))
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
	last := 0
	for i, c := range b {
		var seq string
		switch c {
		case '\000':
			seq = html.Null
		case '"':
			seq = html.Quot
		case '\'':
			seq = html.Apos
		case '&':
			seq = html.Amp
		case '<':
			seq = html.Lt
		case '>':
			seq = html.Gt
		default:
			continue
		}
		w.Body = append(w.Body, b[last:i]...)
		w.Body = append(w.Body, seq...)
		last = i + 1
	}
	w.Body = append(w.Body, b[last:]...)
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
	last := 0
	for i, c := range s {
		var seq string
		switch c {
		case '\000':
			seq = html.Null
		case '"':
			seq = html.Quot
		case '\'':
			seq = html.Apos
		case '&':
			seq = html.Amp
		case '<':
			seq = html.Lt
		case '>':
			seq = html.Gt
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
	w.StatusCode = StatusOK
	w.Headers.Reset()
	w.Body = w.Body[:0]
	w.Arena.Reset()
}
