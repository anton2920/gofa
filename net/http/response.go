package http

import (
	"unsafe"

	"github.com/anton2920/gofa/arena"
	"github.com/anton2920/gofa/database"
	"github.com/anton2920/gofa/net/html"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/syscall"
	"github.com/anton2920/gofa/time"
)

type Response struct {
	Arena arena.Arena

	StatusCode Status
	Headers    Headers
	Bodies     []syscall.Iovec
}

func (w *Response) Append(b []byte) {
	w.Bodies = append(w.Bodies, syscall.IovecForByteSlice(b))
}

func (w *Response) AppendString(s string) {
	w.Bodies = append(w.Bodies, syscall.Iovec(s))
}

func (w *Response) DelCookie(name string) {
	const finisher = "=; Path=/; Max-Age=0; HttpOnly; Secure; SameSite=Strict"

	cookie := w.Arena.NewSlice(len(name) + len(finisher))

	var n int
	n += copy(cookie[n:], name)
	n += copy(cookie[n:], finisher)

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

func (w *Response) SetCookie(name, value string, expiry int64) {
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
	n += time.PutTmRFC822(cookie[n:], time.ToTm(int(expiry)))
	n += copy(cookie[n:], secure)

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

/* SetCookieUnsafe is useful for debugging purposes. It's also more compatible with older browsers. */
func (w *Response) SetCookieUnsafe(name, value string, expiry int64) {
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
	n += time.PutTmRFC822(cookie[n:], time.ToTm(int(expiry)))

	w.SetHeaderUnsafe("Set-Cookie", unsafe.String(unsafe.SliceData(cookie), n))
}

/* SetHeaderUnsafe sets new 'value' for 'header' relying on that memory lives long enough. */
func (w *Response) SetHeaderUnsafe(header string, value string) {
	switch header {
	case "Date":
		w.Headers.OmitDate = true
	case "Server":
		w.Headers.OmitServer = true
	case "ContentType":
		w.Headers.OmitContentType = true
	case "ContentLength":
		w.Headers.OmitContentLength = true
	}

	for i := 0; i < len(w.Headers.Values); i += 4 {
		key := w.Headers.Values[i]
		if header == string(key) {
			w.Headers.Values[i+2] = syscall.Iovec(value)
			return
		}
	}

	w.Headers.Values = append(w.Headers.Values, syscall.Iovec(header), syscall.Iovec(": "), syscall.Iovec(value), syscall.Iovec("\r\n"))
}

func (w *Response) Redirect(path string, code Status) {
	pathBuf := w.Arena.NewSlice(len(path))
	copy(pathBuf, path)

	w.SetHeaderUnsafe("Location", unsafe.String(unsafe.SliceData(pathBuf), len(pathBuf)))
	w.Bodies = w.Bodies[:0]
	w.StatusCode = code
}

func (w *Response) RedirectID(prefix string, id database.ID, code Status) {
	buffer := w.Arena.NewSlice(len(prefix) + 20)
	n := copy(buffer, prefix)
	n += slices.PutInt(buffer[n:], int(id))

	w.SetHeaderUnsafe("Location", unsafe.String(unsafe.SliceData(buffer), n))
	w.Bodies = w.Bodies[:0]
	w.StatusCode = code
}

func (w *Response) Write(b []byte) (int, error) {
	if len(b) == 0 {
		return 0, nil
	}

	buffer := w.Arena.NewSlice(len(b))
	copy(buffer, b)
	w.Append(buffer)

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
		w.Write(b[last:i])
		w.AppendString(seq)
		last = i + 1
	}
	w.Write(b[last:])
}

func (w *Response) WriteInt(i int) (int, error) {
	buffer := w.Arena.NewSlice(20)
	n := slices.PutInt(buffer, i)
	w.Append(buffer[:n])
	return n, nil
}

func (w *Response) WriteID(id database.ID) (int, error) {
	return w.WriteInt(int(id))
}

func (w *Response) WriteString(s string) (int, error) {
	return w.Write(unsafe.Slice(unsafe.StringData(s), len(s)))
}

func (w *Response) WriteHTMLString(s string) {
	w.WriteHTML(unsafe.Slice(unsafe.StringData(s), len(s)))
}
