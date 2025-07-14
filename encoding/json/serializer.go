package json

import (
	"github.com/anton2920/gofa/slices"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/util"
)

type Serializer struct {
	Buf       []byte
	Pos       int
	NeedComma bool
}

func (s *Serializer) PutObjectBegin() {
	s.Buf[s.Pos] = '{'
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) PutObjectEnd() {
	s.Buf[s.Pos] = '}'
	s.Pos++
}

func (s *Serializer) PutArrayBegin() {
	s.Buf[s.Pos] = '['
	s.Pos++
	s.NeedComma = false
}

func (s *Serializer) PutArrayEnd() {
	s.Buf[s.Pos] = ']'
	s.Pos++
}

func (s *Serializer) PutInt(x int) {
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buf[s.Pos:], x)
}

func (s *Serializer) PutInt32(x int32) {
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buf[s.Pos:], int(x))
}

func (s *Serializer) PutUint32(x uint32) {
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buf[s.Pos:], int(x))
}

/* TODO(anton2920): this is incorrect on i386. */
func (s *Serializer) PutInt64(x int64) {
	s.NeedComma = true
	s.Pos += slices.PutInt(s.Buf[s.Pos:], int(x))
}

func (s *Serializer) PutString(str string) {
	s.NeedComma = true

	s.Buf[s.Pos] = '"'
	s.Pos++

	for {
		quote := strings.FindChar(str, '"')
		if quote == -1 {
			s.Pos += copy(s.Buf[s.Pos:], str)
			break
		}
		s.Pos += copy(s.Buf[s.Pos:], str[:quote])
		s.Pos += copy(s.Buf[s.Pos:], `\"`)
		if quote == len(str)-1 {
			break
		}
		str = str[quote+1:]
	}

	s.Buf[s.Pos] = '"'
	s.Pos++
}

func (s *Serializer) PutKey(key string) {
	if s.NeedComma {
		s.NeedComma = false
		s.Buf[s.Pos] = ','
		s.Pos++
	}

	s.PutString(key)

	s.Buf[s.Pos] = ':'
	s.Pos++
}

func (s *Serializer) Reset() {
	s.Pos = 0
	s.NeedComma = false
}

func (s *Serializer) Bytes() []byte {
	return s.Buf[:s.Pos]
}

func (s Serializer) String() string {
	return util.Slice2String(s.Buf[:s.Pos])
}
