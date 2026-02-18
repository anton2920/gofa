package bplus

import "unsafe"

type Meta struct {
	PageHeader

	Magic   uint64
	Version uint64

	Root        int64
	EndSentinel int64

	LastSeq uint64

	_ [PageSize - PageHeaderSize - 5*unsafe.Sizeof(int64(0))]byte
}

func (m *Meta) Page() *Page {
	return (*Page)(unsafe.Pointer(m))
}
