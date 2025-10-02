package alloc

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
	"github.com/anton2920/gofa/trace"
)

type Arena struct {
	Buffer []byte
	Used   int
}

func (a *Arena) NewSlice(n int) []byte {
	t := trace.Begin("")

	if a.Used+n >= cap(a.Buffer) {
		a.Buffer = make([]byte, ints.Max(a.Used+n*2, cap(a.Buffer)*2))
		a.Used = 0
	}
	ret := a.Buffer[a.Used : a.Used+n]
	a.Used += n

	trace.End(t)
	return ret
}

func (a *Arena) Copy(bs []byte) []byte {
	t := trace.Begin("")

	ret := a.NewSlice(len(bs))
	copy(ret, bs)

	trace.End(t)
	return ret
}

func (a *Arena) CopyString(s string) string {
	t := trace.Begin("")

	ret := a.NewSlice(len(s))
	copy(ret, s)

	trace.End(t)
	return bytes.AsString(ret)
}

func (a *Arena) Reset() {
	a.Used = 0
}
