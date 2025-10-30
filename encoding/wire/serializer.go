package wire

import (
	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/trace"
)

type Serializer struct {
	Buffer []byte
}

func (s *Serializer) Begin(version byte) {
	t := trace.Begin("")

	s.PutByte(version)

	trace.End(t)
}

func (s *Serializer) PutType(typ ValueType) {
	t := trace.Begin("")

	s.Buffer = append(s.Buffer, byte(typ))

	trace.End(t)
}

func (s *Serializer) PutByte(b byte) {
	t := trace.Begin("")

	s.PutType(ValueTypeByte)
	s.Buffer = append(s.Buffer, b)

	trace.End(t)
}

func (s *Serializer) PutInt8(i int8) {
	t := trace.Begin("")

	s.PutByte(byte(i))

	trace.End(t)
}

func (s *Serializer) PutInt32(i int32) {
	t := trace.Begin("")

	s.PutType(ValueTypeInt32)
	s.Buffer = append(s.Buffer, byte((i>>0)&0xFF), byte((i>>8)&0xFF), byte((i>>16)&0xFF), byte((i>>24)&0xFF))

	trace.End(t)
}

func (s *Serializer) PutFlags(f bits.Flags) {
	t := trace.Begin("")

	s.PutInt32(int32(f))

	trace.End(t)
}

func (s *Serializer) PutString(str string) {
	t := trace.Begin("")

	s.PutType(ValueTypeString)
	s.PutInt32(int32(len(str)))
	s.Buffer = append(s.Buffer, str...)

	trace.End(t)
}

func (s *Serializer) PutSliceBegin(l int) {
	t := trace.Begin("")

	s.PutType(ValueTypeSlice)
	s.PutInt32(int32(l))

	trace.End(t)
}

func (s *Serializer) End() {
	/* TODO(anton2920): calculate CRC32 or something... */
}
