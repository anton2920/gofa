package types

/* Must be in sync with 'reflect.StringHeader'. */
type StringHeader struct {
	Data uintptr
	Len  int
}
