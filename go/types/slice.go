package types

/* Must be in sync with the 'reflect.SliceHeader'. */
type SliceHeader struct {
	Data uintptr
	Len  int
	Cap  int
}
