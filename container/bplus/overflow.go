package bplus

import "unsafe"

type Overflow struct {
	PageHeader

	Next int64

	Data [PageSize - PageHeaderSize - unsafe.Sizeof(int64(0))]byte
}

func (o *Overflow) SetValue(value []byte) []byte {
	o.Head = uint16(copy(o.Data[:], value[len(value)-len(o.Data):]))
	return value[:len(value)-int(o.Head)]
}

func (o *Overflow) GetValue() []byte {
	return o.Data[:o.Head]
}
