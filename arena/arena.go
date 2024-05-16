package arena

type Arena struct {
	Buffer []byte
	Used   int
}

func (a *Arena) NewSlice(n int) []byte {
	if a.Used+n >= cap(a.Buffer) {
		buffer := make([]byte, max(a.Used+n*2, cap(a.Buffer)*2))
		copy(buffer, a.Buffer)
		a.Buffer = buffer
	}
	ret := a.Buffer[a.Used : a.Used+n]
	a.Used += n
	return ret
}

func (a *Arena) Reset() {
	a.Used = 0
}
