package fmt

import (
	"strconv"
	"unicode/utf8"
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/time"
)

type Formatter struct {
	Buffer []byte
	Pos    int

	Width     int
	Precision int
}

func (f *Formatter) applyWidth(n int, after bool) {
	var leftAlign bool

	if n == 0 {
		return
	}

	width := f.Width
	if width < 0 {
		leftAlign = true
		width = -width
	}

	if leftAlign == after {
		for i := 0; i < width-n; i++ {
			f.Buffer[f.Pos] = ' '
			f.Pos++
		}
		f.Width = 0
	}
}

func (f *Formatter) InitWithUnsafePointer(ptr unsafe.Pointer, n int) {
	f.Buffer = bytes.SliceFromUnsafePointer(ptr, n)
}

func (f *Formatter) InitWithBytePointer(ptr *byte, n int) {
	f.Buffer = bytes.SliceFromBytePointer(ptr, n)
}

func (f *Formatter) InitWithByteSlice(buf []byte) {
	f.Buffer = bytes.SliceFromBytePointer(&buf[0], len(buf))
}

func (f *Formatter) D(d int) *Formatter {
	buf := make([]byte, ints.Bufsize)
	n := slices.PutInt(buf, d)
	return f.S(bytes.AsString(buf[:n]))
}

func (f *Formatter) D32(d int32) *Formatter {
	buf := make([]byte, ints.Bufsize)
	n := slices.PutInt(buf, int(d))
	return f.S(bytes.AsString(buf[:n]))
}

func (f *Formatter) D64(d int64) *Formatter {
	buf := make([]byte, 0, ints.Bufsize)
	buf = strconv.AppendInt(buf, d, 10)
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) Date(t int64) *Formatter {
	buf := make([]byte, 10)
	time.PutTmDate(buf, time.ToTm(t))
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) DateTime(t int64) *Formatter {
	buf := make([]byte, 19)
	time.PutTmDateTime(buf, time.ToTm(t))
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) Err(err error) *Formatter {
	s := "<nil>"
	if err != nil {
		s = err.Error()
	}
	return f.S(s)
}

func (f *Formatter) E(e float64) *Formatter {
	return f.E64(e)
}

func (f *Formatter) E32(e float32) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, float64(e), 'e', ints.Or(f.Precision, 6), int(unsafe.Sizeof(e)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) E64(e float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, e, 'e', ints.Or(f.Precision, 6), int(unsafe.Sizeof(e)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) F(f_ float64) *Formatter {
	return f.F64(f_)
}

func (f *Formatter) F32(f_ float32) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, float64(f_), 'f', ints.Or(f.Precision, 6), int(unsafe.Sizeof(f_)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) F64(f_ float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, f_, 'f', ints.Or(f.Precision, 6), int(unsafe.Sizeof(f_)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) G(g float64) *Formatter {
	return f.G64(g)
}

func (f *Formatter) G32(g float32) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, float64(g), 'g', ints.Or(f.Precision, -1), int(unsafe.Sizeof(g)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) G64(g float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, g, 'g', ints.Or(f.Precision, -1), int(unsafe.Sizeof(g)*8))
	f.Precision = 0
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) I(i int) *Formatter {
	return f.D(i)
}

func (f *Formatter) I32(i int32) *Formatter {
	return f.D32(i)
}

func (f *Formatter) I64(i int64) *Formatter {
	return f.D64(i)
}

func (f *Formatter) Ln() *Formatter {
	return f.S(LineTerminator)
}

func (f *Formatter) P(p unsafe.Pointer) *Formatter {
	const prefix = "0x"

	buf := make([]byte, len(prefix), ints.Bufsize)
	copy(buf, prefix)
	buf = strconv.AppendInt(buf, int64(uintptr(p)), 16)

	return f.S(bytes.AsString(buf))
}

func quotedStringLength(s string) int {
	var runeTmp [utf8.UTFMax]byte
	quote := '"'
	size := 2

	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if (width == 1) && (r == utf8.RuneError) {
			size += 2 + 2
			continue
		}
		if (r == rune(quote)) || (r == '\\') {
			size += 2
			continue
		}
		if strconv.IsPrint(r) {
			size += utf8.EncodeRune(runeTmp[:], r)
			continue
		}
		switch r {
		case '\a', '\b', '\f', '\n', '\r', '\t', '\v':
			size++
		default:
			switch {
			case r < ' ':
				size += 2 + 2
			case r > utf8.MaxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				size += 2 + 4
			default:
				size += 2 + 8
			}
		}
	}

	return size
}

var lowerhex = []string{"0", "1", "2", "3", "4", "5", "6", "7", "8", "9", "a", "b", "c", "d", "e", "f"}

func (f *Formatter) quoteString(s string) {
	var runeTmp [utf8.UTFMax]byte
	quote := '"'

	f.Pos += copy(f.Buffer[f.Pos:], "\"")
	for width := 0; len(s) > 0; s = s[width:] {
		r := rune(s[0])
		width = 1
		if r >= utf8.RuneSelf {
			r, width = utf8.DecodeRuneInString(s)
		}
		if width == 1 && r == utf8.RuneError {
			f.Pos += copy(f.Buffer[f.Pos:], `\x`)
			f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]>>4])
			f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]&0xF])
			continue
		}
		if r == rune(quote) || r == '\\' {
			f.Pos += copy(f.Buffer[f.Pos:], "\\")
			f.Pos += copy(f.Buffer[f.Pos:], s[:1])
			continue
		}
		if strconv.IsPrint(r) {
			n := utf8.EncodeRune(runeTmp[:], r)
			f.Pos += copy(f.Buffer[f.Pos:], runeTmp[:n])
			continue
		}
		switch r {
		case '\a':
			f.Pos += copy(f.Buffer[f.Pos:], `\a`)
		case '\b':
			f.Pos += copy(f.Buffer[f.Pos:], `\b`)
		case '\f':
			f.Pos += copy(f.Buffer[f.Pos:], `\f`)
		case '\n':
			f.Pos += copy(f.Buffer[f.Pos:], `\n`)
		case '\r':
			f.Pos += copy(f.Buffer[f.Pos:], `\r`)
		case '\t':
			f.Pos += copy(f.Buffer[f.Pos:], `\t`)
		case '\v':
			f.Pos += copy(f.Buffer[f.Pos:], `\v`)
		default:
			switch {
			case r < ' ':
				f.Pos += copy(f.Buffer[f.Pos:], `\x`)
				f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]>>4])
				f.Pos += copy(f.Buffer[f.Pos:], lowerhex[s[0]&0xF])
			case r > utf8.MaxRune:
				r = 0xFFFD
				fallthrough
			case r < 0x10000:
				f.Pos += copy(f.Buffer[f.Pos:], `\u`)
				for s := 12; s >= 0; s -= 4 {
					f.Pos += copy(f.Buffer[f.Pos:], lowerhex[r>>uint(s)&0xF])
				}
			default:
				f.Pos += copy(f.Buffer[f.Pos:], `\U`)
				for s := 28; s >= 0; s -= 4 {
					f.Pos += copy(f.Buffer[f.Pos:], lowerhex[r>>uint(s)&0xF])
				}
			}
		}
	}
	f.Pos += copy(f.Buffer[f.Pos:], "\"")
}

func (f *Formatter) Q(q string) *Formatter {
	size := quotedStringLength(q)

	f.applyWidth(size, false)
	f.quoteString(q)
	f.applyWidth(size, true)

	return f
}

func (f *Formatter) S(s string) *Formatter {
	f.applyWidth(len(s), false)
	f.Pos += copy(f.Buffer[f.Pos:], s)
	f.applyWidth(len(s), true)
	return f
}

func (f *Formatter) W(width int) *Formatter {
	f.Width = width
	return f
}

func (f *Formatter) Prec(prec int) *Formatter {
	f.Precision = prec
	return f
}

func (f *Formatter) Bytes() []byte {
	return f.Buffer[:f.Pos]
}

func (f *Formatter) String() string {
	return bytes.AsString(f.Bytes())
}

func (f *Formatter) Reset() *Formatter {
	f.Pos = 0
	f.Width = 0
	f.Precision = 0
	return f
}
