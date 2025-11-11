package wire

import "github.com/anton2920/gofa/trace"

type Serializer struct {
	Buffer []byte
}

func (s *Serializer) Begin(version int) {
	t := trace.Begin("")

	s.Int32(int32(version))

	trace.End(t)
}

func (s *Serializer) Int32(n int32) {
	s.Uint32(uint32(n))
}

func (s *Serializer) Uint32(n uint32) {
	t := trace.Begin("")

	s.Buffer = append(s.Buffer, byte((n>>0)&0xFF), byte((n>>8)&0xFF), byte((n>>16)&0xFF), byte((n>>24)&0xFF))

	trace.End(t)
}

func (s *Serializer) String(str string) {
	t := trace.Begin("")

	s.Int32(int32(len(str)))
	s.Buffer = append(s.Buffer, str...)

	trace.End(t)
}

func (s *Serializer) End() {
	/* TODO(anton2920): calculate CRC32 or something... */
}

func (s *Serializer) Reset() {
	s.Buffer = s.Buffer[:0]
}
