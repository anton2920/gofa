package html

import (
	"github.com/anton2920/gofa/net/http"
)

type HTML struct {
	*http.Response
	*http.Request
}

func New(w *http.Response, r *http.Request) HTML {
	return HTML{Response: w, Request: r}
}

func (h *HTML) Backspace() *HTML {
	h.Response.Body = h.Response.Body[:len(h.Response.Body)-1]
	return h
}

func (h *HTML) Bytes(bs []byte) *HTML {
	h.Write(bs)
	return h
}

/*
func (h *HTML) Date(d int64) *HTML {
	return h.String(h.Dtoa(d))
}
*/

func (h *HTML) Int(n int) *HTML {
	//h.WriteInt(n)
	return h
}

func (h *HTML) HString(s string) *HTML {
	h.WriteHTMLString(s)
	return h
}

func (h *HTML) LString(s string) *HTML {
	return h.HString(h.L(s))
}

func (h *HTML) LStringColon(s string) *HTML {
	h.LString(s)
	h.String(": ")
	return h
}

/*
func (h *HTML) LStringPlural(s string, n int) *HTML {
	const suffix = "ies"

	if !strings.EndsWith(s, suffix) {
		h.LString(s[:len(s)-bools.ToInt(n == 1)])
	} else {
		buf := h.Response.Arena.NewSlice(len(s))
		copy(buf, s)

		if n == 1 {
			buf = buf[:len(s)-len(suffix)+1]
			buf[len(buf)-1] = 'y'
		}

		h.LString(bytes.AsString(buf))
	}

	return h
}
*/

func (h *HTML) String(s string) *HTML {
	h.WriteString(s)
	return h
}

/*
func (h *HTML) Time(t int64) *HTML {
	h.String(h.Ttoa(t))
	return h
}

func (h *HTML) TString(s string) *HTML {
	if s := h.L(s); len(s) > 1 {
		h.String(string(unicode.ToUpper(rune(s[0]))))
		h.String(s[1:])
	}
	return h
}

func (h *HTML) Dtoa(d int64) string {
	v := d + int64(h.Timezone)*time.Hour

	buf := h.Response.Arena.NewSlice(len("2006-01-02"))
	n := time.PutTmDate(buf, time.ToTm(v))

	return bytes.AsString(buf[:n])
}

func (h *HTML) Dtoa1(d int64) string {
	if d == 0 {
		return ""
	}
	return h.Dtoa(d)
}

func (h *HTML) Itoa(x int) string {
	buf := h.Response.Arena.NewSlice(ints.Bufsize)
	n := slices.PutInt(buf, x)
	return bytes.AsString(buf[:n])
}


func (h *HTML) Itoa1(x int) string {
	if x == 0 {
		return ""
	}
	return h.Itoa(x)
}

func (h *HTML) Ttoa(t int64) string {
	v := t + int64(h.Timezone)*time.Hour

	buf := h.Response.Arena.NewSlice(len("2006-01-02 15:04:05"))
	n := time.PutTmDateTime(buf, time.ToTm(v))

	return bytes.AsString(buf[:n])
}

func (h *HTML) IndexedName(name string, indicies ...int) string {
	var n int

	buf := h.Response.Arena.NewSlice(len(name) + ints.Bufsize*len(indicies) + 2*len(indicies))
	n += copy(buf[n:], name)
	for i := 0; i < len(indicies); i++ {
		index := indicies[i]
		if i > 0 {
			n += copy(buf[n:], ".")
		}
		n += slices.PutInt(buf[n:], index)
	}

	return bytes.AsString(buf[:n])
}
*/
