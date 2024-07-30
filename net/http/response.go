package http

import (
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/prof"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
)

type Response struct {
	Arena arena.Arena

	StatusCode Status
	Headers    Headers
	Body       []byte
}

func (w *Response) DelCookie(name string) {
	p := prof.Begin("")

	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.Headers.Set("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))

	prof.End(p)
}

func (w *Response) SetCookie(name, value string, expiry int) {
	p := prof.Begin("")

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

	prof.End(p)
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int) {
	p := prof.Begin("")

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

	prof.End(p)
}

func (w *Response) Redirect(path string, code Status) {
	p := prof.Begin("")

	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.Headers.Set("Location", unsafe.String(unsafe.SliceData(pathBuf), len(pathBuf)))
	w.Body = w.Body[:0]
	w.StatusCode = code

	prof.End(p)
}

func (w *Response) RedirectID(prefix string, id database.ID, code Status) {
	p := prof.Begin("")

	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], int(id))

	w.Headers.Set("Location", unsafe.String(unsafe.SliceData(buffer), n))
	w.Body = w.Body[:0]
	w.StatusCode = code

	prof.End(p)
}

//go:nosplit
func (w *Response) Write(b []byte) (int, error) {
	p := prof.Begin("")

	w.Body = append(w.Body, b...)

	prof.End(p)
	return len(b), nil
}

/* WriteHTML writes to w the escaped html. equivalent of the plain text data b. */
//go:nosplit
func (w *Response) WriteHTML(b []byte) {
	p := prof.Begin("")

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

	prof.End(p)
}

func (w *Response) WriteInt(i int) {
	p := prof.Begin("")

	buffer := make([]byte, 20)
	n := slices.PutInt(buffer, i)
	w.Write(buffer[:n])

	prof.End(p)
}

func (w *Response) WriteID(id database.ID) {
	p := prof.Begin("")

	w.WriteInt(int(id))

	prof.End(p)
}

//go:nosplit
func (w *Response) WriteString(s string) (int, error) {
	p := prof.Begin("")

	w.Body = append(w.Body, s...)

	prof.End(p)
	return len(s), nil
}

func (w *Response) WriteHTMLString(s string) {
	p := prof.Begin("")

	w.WriteHTML(unsafe.Slice(unsafe.StringData(s), len(s)))

	prof.End(p)
}

//go:nosplit
func (w *Response) Reset() {
	p := prof.Begin("")

	w.StatusCode = StatusOK
	w.Headers.Reset()
	w.Body = w.Body[:0]
	w.Arena.Reset()

	prof.End(p)
}
