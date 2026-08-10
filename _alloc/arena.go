package alloc

import (
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace/trace_"
)

type Arena struct {
	Buffer []byte
	Used   int
}

func (a *Arena) NewSlice(n int) []byte {
	t := trace_.Begin("")

	if a.Used+n >= cap(a.Buffer) {
		a.Buffer = make([]byte, ints.Max(a.Used+n*2, cap(a.Buffer)*2))
		a.Used = 0
	}
	ret := a.Buffer[a.Used : a.Used+n]
	a.Used += n

	trace_.End(t)
	return ret
}

func (a *Arena) Copy(bs []byte) []byte {
	t := trace_.Begin("")

	ret := a.NewSlice(len(bs))
	copy(ret, bs)

	trace_.End(t)
	return ret
}

func (a *Arena) CopyString(s string) string {
	t := trace_.Begin("")

	ret := a.NewSlice(len(s))
	copy(ret, s)

	trace_.End(t)
	return bytes.AsString(ret)
}

func (a *Arena) EscapeString(s string) string {
	bs := strings.AsBytes(s)
	if (uintptr(unsafe.Pointer(&bs[0])) >= uintptr(unsafe.Pointer(&a.Buffer[0]))) && (uintptr(unsafe.Pointer(&bs[0])) <= uintptr(unsafe.Pointer(&a.Buffer[a.Used-1]))) {
		return string(bs)
	}
	return s
}

func (a *Arena) Reset() {
	a.Used = 0
}
