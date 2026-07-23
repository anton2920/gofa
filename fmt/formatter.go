package fmt

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/slices"
)

type Formatter struct {
	Buffer []byte
	Pos    int

	Width     int
	Precision int
}

func (f *Formatter) applyWidth(n int, after bool) {
	var leftAlign bool

	debug.Assert(n > 0, "field length must be >0")

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

func (f *Formatter) InitWithByteSlice(buf []byte) {
	f.Buffer = bytes.SliceFromBytePointer(&buf[0], len(buf))
}

func (f *Formatter) D(d int) *Formatter {
	buf := make([]byte, ints.Bufsize)
	n := slices.PutInt(buf, d)
	return f.S(bytes.AsString(buf[:n]))
}

func (f *Formatter) D64(d int64) *Formatter {
	buf := make([]byte, ints.Bufsize)
	n := slices.PutInt(buf, int(d))
	return f.S(bytes.AsString(buf[:n]))
}

func (f *Formatter) E(e float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, e, 'e', ints.Or(f.Precision, 6), int(unsafe.Sizeof(e)*8))
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) Err(err error) *Formatter {
	s := "<nil>"
	if err != nil {
		s = err.Error()
	}
	return f.S(s)
}

func (f *Formatter) F(f_ float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, f_, 'f', ints.Or(f.Precision, 6), int(unsafe.Sizeof(f_)*8))
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) G(g float64) *Formatter {
	buf := make([]byte, 0, 128)
	buf = strconv.AppendFloat(buf, g, 'g', ints.Or(f.Precision, -1), int(unsafe.Sizeof(g)*8))
	return f.S(bytes.AsString(buf))
}

func (f *Formatter) I(i int) *Formatter {
	return f.D(i)
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
	return f
}
