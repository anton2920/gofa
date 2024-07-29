package http

import (
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
	"github.com/anton2920/gofa/trace"
)

type Response struct {
	Arena arena.Arena

	StatusCode Status
	Headers    Headers
	Body       []byte
}

func (w *Response) DelCookie(name string) {
	defer trace.End(trace.Start(""))

	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.Headers.Set("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

func (w *Response) SetCookie(name, value string, expiry int) {
	defer trace.End(trace.Start(""))

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

	w.Headers.Set("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int) {
	defer trace.End(trace.Start(""))

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

	w.Headers.Set("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

func (w *Response) Redirect(path string, code Status) {
	defer trace.End(trace.Start(""))

	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.Headers.Set("Location", unsafe.String(unsafe.SliceData(pathBuf), len(pathBuf)))
	w.Body = w.Body[:0]
	w.StatusCode = code
}

func (w *Response) RedirectID(prefix string, id database.ID, code Status) {
	defer trace.End(trace.Start(""))

	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], int(id))

	w.Headers.Set("Location", unsafe.String(unsafe.SliceData(buffer), n))
	w.Body = w.Body[:0]
	w.StatusCode = code
}

func (w *Response) Write(b []byte) (int, error) {
	defer trace.End(trace.Start(""))

	w.Body = append(w.Body, b...)
	return len(b), nil
}

/* WriteHTML writes to w the escaped html. equivalent of the plain text data b. */
func (w *Response) WriteHTML(b []byte) {
	defer trace.End(trace.Start(""))

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
		w.Write(b[last:i])
		w.WriteString(seq)
		last = i + 1
	}
	w.Write(b[last:])
}

func (w *Response) WriteInt(i int) (int, error) {
	defer trace.End(trace.Start(""))

	buffer := make([]byte, 20)
	n := slices.PutInt(buffer, i)
	w.Write(buffer[:n])
	return n, nil
}

func (w *Response) WriteID(id database.ID) (int, error) {
	defer trace.End(trace.Start(""))

	return w.WriteInt(int(id))
}

func (w *Response) WriteString(s string) (int, error) {
	defer trace.End(trace.Start(""))
	return w.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (w *Response) WriteHTMLString(s string) {
	defer trace.End(trace.Start(""))

	w.WriteHTML(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (w *Response) Reset() {
	defer trace.End(trace.Start(""))

	w.StatusCode = StatusOK
	w.Headers.Reset()
	w.Body = w.Body[:0]
	w.Arena.Reset()
}
