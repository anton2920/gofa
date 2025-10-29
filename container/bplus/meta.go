package bplus

import "unsafe"

type Meta struct {
	PageHeader

	Magic   int64
	Version int64

	Root        int64
	EndSentinel int64

	_ [PageSize - PageHeaderSize - 4*unsafe.Sizeof(int64(0))]byte
}

func (m *Meta) Page() *Page {
	return (*Page)(unsafe.Pointer(m))
}
