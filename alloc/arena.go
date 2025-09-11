package alloc

import "github.com/anton2920/gofa/ints"

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

func (a *Arena) Reset() {
	a.Used = 0
}
