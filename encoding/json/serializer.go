package json

import (
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
)

type Serializer struct {
	Buffer    []byte
	Pos       int
	NeedComma bool
}

func (s *Serializer) PutComma() {
	if s.NeedComma {
		s.NeedComma = false
		s.Buffer[s.Pos] = ','
		s.Pos++
	}
}

func (s *Serializer) PutObjectBegin() {
	s.PutComma()
	s.Buffer[s.Pos] = '{'
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) PutObjectEnd() {
	s.Buffer[s.Pos] = '}'
	s.Pos++
}

func (s *Serializer) PutArrayBegin() {
	s.PutComma()
	s.Buffer[s.Pos] = '['
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) PutArrayEnd() {
	s.Buffer[s.Pos] = ']'
	s.Pos++
}

func (s *Serializer) PutInt(x int) {
	s.PutComma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], x)
}

func (s *Serializer) PutInt32(x int32) {
	s.PutComma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

func (s *Serializer) PutUint32(x uint32) {
	s.PutComma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

/* TODO(anton2920): this is incorrect on i386. */
func (s *Serializer) PutInt64(x int64) {
	s.PutComma()
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buffer[s.Pos:], int(x))
}

func (s *Serializer) PutString(str string) {
	s.PutComma()
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

func (s *Serializer) PutKey(key string) {
	s.PutComma()
	s.PutString(key)

	s.Buffer[s.Pos] = ':'
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) Reset() {
	s.Pos = 0
	s.NeedComma = false
}

func (s *Serializer) Bytes() []byte {
	return s.Buffer[:s.Pos]
}

func (s Serializer) String() string {
	return bytes.AsString(s.Buffer[:s.Pos])
}
