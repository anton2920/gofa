package fmt

import (
	"unsafe"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/ints"
)

type Scanner struct {
	Buffer []byte
	Pos    int
	MaxPos int

	Requested int
	Matched   int
}

func isspace(ch byte) bool {
	return (ch == ' ') || (ch == '\t') || (ch == '\n') || (ch == '\r')
}

func isdigit(ch byte) bool {
	return (ch >= '0') && (ch <= '9')
}

func (s *Scanner) InitWithUnsafePointer(ptr unsafe.Pointer, n int) *Scanner {
	s.Buffer = bytes.SliceFromUnsafePointer(ptr, n)
	return s.Reset()
}

func (s *Scanner) InitWithBytePointer(ptr *byte, n int) *Scanner {
	return s.InitWithUnsafePointer(unsafe.Pointer(ptr), n)
}

func (s *Scanner) InitWithByteSlice(buf []byte) *Scanner {
	return s.InitWithBytePointer(&buf[0], len(buf))
}

func (s *Scanner) SkipSpaces() *Scanner {
	for ; (s.Pos < len(s.Buffer)) && (isspace(s.Buffer[s.Pos])); s.Pos++ {
	}
	return s
}

func (s *Scanner) S(ps *[]byte) *Scanner {
	s.Requested++
	s.SkipSpaces()

	start := s.Pos
	for ; (s.Pos < s.MaxPos) && (!isspace(s.Buffer[s.Pos])); s.Pos++ {
	}

	if (s.Pos < s.MaxPos) && (isspace(s.Buffer[s.Pos])) {
		n := copy(*ps, bytes.AsString(s.Buffer[start:s.Pos]))
		*ps = (*ps)[:n]
		s.Matched++
	}

	return s
}

func (s *Scanner) D(pd *int) *Scanner {
	s.Requested++

	if s.MaxPos-s.Pos > 0 {
		var num int

		sign := 1
		switch s.Buffer[s.Pos] {
		case '-':
			sign = -1
			fallthrough
		case '+':
			s.Pos++
		}

		for ; (s.Pos < s.MaxPos) && (isdigit(s.Buffer[s.Pos])); s.Pos++ {
			num = num*10 + int(s.Buffer[s.Pos]-'0')
		}

		*pd = sign * num
		s.Matched++
	}

	return s
}

func (s *Scanner) Match(part string) *Scanner {
	s.Requested++

	minLen := ints.Min(len(part), s.MaxPos-s.Pos)
	s.Matched += bools.ToInt(part == bytes.AsString(s.Buffer[s.Pos:s.Pos+minLen]))
	s.Pos += minLen

	return s
}

func (s *Scanner) OK() bool {
	return s.Requested == s.Matched
}

func (s *Scanner) Reset() *Scanner {
	s.Pos = 0
	s.MaxPos = len(s.Buffer)

	s.Requested = 0
	s.Matched = 0

	return s
}
