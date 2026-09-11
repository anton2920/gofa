//go:build gofanostd
// +build gofanostd

package fmt

import "unsafe"

func (f *Formatter) D64(d int64) *Formatter {
	return f
}

func (f *Formatter) E32(e float32) *Formatter {
	return f
}

func (f *Formatter) E64(e float64) *Formatter {
	return f
}

func (f *Formatter) F32(f_ float32) *Formatter {
	return f
}

func (f *Formatter) F64(f_ float64) *Formatter {
	return f
}

func (f *Formatter) G32(g float32) *Formatter {
	return f
}

func (f *Formatter) G64(g float64) *Formatter {
	return f
}

func (f *Formatter) P(p unsafe.Pointer) *Formatter {
	return f
}
