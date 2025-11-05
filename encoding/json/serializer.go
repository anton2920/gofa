package json

import (
	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/encoding"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

type Serializer struct {
	Buffer    []byte
	Pos       int
	NeedComma bool
}

var _ encoding.Serializer = new(Serializer)

func (s *Serializer) Comma() {
	if s.NeedComma {
		s.NeedComma = false
		s.Buffer[s.Pos] = ','
		s.Pos++
	}
}

func (s *Serializer) Begin() {

}

func (s *Serializer) End() {

}

func (s *Serializer) ObjectBegin() {
	s.Comma()
	s.Buffer[s.Pos] = '{'
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) ObjectEnd() {
	s.Buffer[s.Pos] = '}'
	s.Pos++
}

func (s *Serializer) ArrayBegin() {
	s.Comma()
	s.Buffer[s.Pos] = '['
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) ArrayEnd() {
	s.Buffer[s.Pos] = ']'
	s.Pos++
}

func (s *Serializer) Bool(b bool) {
	values := []string{"false", "true"}

	s.Comma()
	s.Pos += copy(s.Buffer[s.Pos:], values[bools.ToInt(b)])
}

func (s *Serializer) Int32(x int32) {
	s.Comma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

func (s *Serializer) Uint32(x uint32) {
	s.Comma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

/* TODO(anton2920): this is incorrect on i386. */
func (s *Serializer) Int64(x int64) {
	s.Comma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

func (s *Serializer) String(str string) {
	s.Comma()
	s.NeedComma = true

	s.Buffer[s.Pos] = '"'
	s.Pos++

	for {
		quote := strings.FindChar(str, '"')
		if quote == -1 {
			s.Pos += copy(s.Buffer[s.Pos:], str)
			break
		}
		s.Pos += copy(s.Buffer[s.Pos:], str[:quote])
		s.Pos += copy(s.Buffer[s.Pos:], `\"`)
		if quote == len(str)-1 {
			break
		}
		str = str[quote+1:]
	}

	s.Buffer[s.Pos] = '"'
	s.Pos++
}

func (s *Serializer) Key(key string) {
	s.Comma()
	s.String(key)

	s.Buffer[s.Pos] = ':'
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) Bytes() []byte {
	return s.Buffer[:s.Pos]
}

func (s *Serializer) Reset() {
	s.Pos = 0
	s.NeedComma = false
}
