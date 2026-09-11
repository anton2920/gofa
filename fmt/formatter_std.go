//go:build !gofanostd
// +build !gofanostd

package fmt

import (
	"strconv"
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
)

func (f *Formatter) D64(d int64) *Formatter {
	buf := make([]byte, 0, ints.Bufsize)
	buf = strconv.AppendInt(buf, d, 10)
	return f.S(bytes.AsString(buf))
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

func (f *Formatter) P(p unsafe.Pointer) *Formatter {
	const prefix = "0x"

	buf := make([]byte, len(prefix), ints.Bufsize)
	copy(buf, prefix)
	buf = strconv.AppendInt(buf, int64(uintptr(p)), 16)

	return f.S(bytes.AsString(buf))
}
