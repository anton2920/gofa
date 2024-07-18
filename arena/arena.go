package arena

type Arena struct {
}

func (a *Arena) NewSlice(n int) []byte {
	return make([]byte, n)
}

func (a *Arena) Reset() {
}
