package wire

import (
	"unsafe"

	"github.com/anton2920/gofa/trace/trace_"
)

type Serializer struct {
	Buffer []byte
}

func (s *Serializer) Begin(version int) {
	t := trace_.Begin("")

	s.Int32(int32(version))

	trace_.End(t)
}

func (s *Serializer) Int8(n int8) {
	s.Uint8(uint8(n))
}

func (s *Serializer) Uint8(n uint8) {
	t := trace_.Begin("")

	s.Buffer = append(s.Buffer, byte((n>>0)&0xFF))

	trace_.End(t)
}

func (s *Serializer) Int16(n int16) {
	s.Uint16(uint16(n))
}

func (s *Serializer) Uint16(n uint16) {
	t := trace_.Begin("")

	s.Buffer = append(s.Buffer, byte((n>>0)&0xFF), byte((n>>8)&0xFF))

	trace_.End(t)
}

func (s *Serializer) Int32(n int32) {
	s.Uint32(uint32(n))
}

func (s *Serializer) Uint32(n uint32) {
	t := trace_.Begin("")

	s.Buffer = append(s.Buffer, byte((n>>0)&0xFF), byte((n>>8)&0xFF), byte((n>>16)&0xFF), byte((n>>24)&0xFF))

	trace_.End(t)
}

func (s *Serializer) Int64(n int64) {
	s.Uint64(uint64(n))
}

func (s *Serializer) Uint64(n uint64) {
	t := trace_.Begin("")

	s.Buffer = append(s.Buffer, byte((n>>0)&0xFF), byte((n>>8)&0xFF), byte((n>>16)&0xFF), byte((n>>24)&0xFF), byte((n>>32)&0xFF), byte((n>>40)&0xFF), byte((n>>48)&0xFF), byte((n>>56)&0xFF))

	trace_.End(t)
}

func (s *Serializer) Int(n int) {
	/* NOTE(anton2920): for maximum compatibility, assuming that 'int' is 'int64'. */
	s.Int64(int64(n))
}

func (s *Serializer) Float64(f float64) {
	s.Uint64(*(*uint64)(unsafe.Pointer(&f)))
}

func (s *Serializer) String(str string) {
	t := trace_.Begin("")

	s.Uint32(uint32(len(str)))
	s.Buffer = append(s.Buffer, str...)

	trace_.End(t)
}

func (s *Serializer) Bytes(bytes []byte) {
	t := trace_.Begin("")

	s.Uint32(uint32(len(bytes)))
	s.Buffer = append(s.Buffer, bytes...)

	trace_.End(t)
}

func (s *Serializer) End() {
	/* TODO(anton2920): calculate CRC32 or something... */
}
