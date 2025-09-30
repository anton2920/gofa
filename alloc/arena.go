package alloc

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
)

type Arena struct {
	Buffer []byte
	Used   int
}

func (a *Arena) NewSlice(n int) []byte {
	if a.Used+n >= cap(a.Buffer) {
		a.Buffer = make([]byte, ints.Max(a.Used+n*2, cap(a.Buffer)*2))
		a.Used = 0
	}
	ret := a.Buffer[a.Used : a.Used+n]
	a.Used += n
	return ret
}

func (a *Arena) Copy(bs []byte) []byte {
	ret := a.NewSlice(len(bs))
	copy(ret, bs)
	return ret
}

func (a *Arena) CopyString(s string) string {
	ret := a.NewSlice(len(s))
	copy(ret, s)
	return bytes.AsString(ret)
}

func (a *Arena) Reset() {
	a.Used = 0
}
