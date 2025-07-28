package json

import (
	"bytes"
	"go/token"

	"github.com/anton2920/gofa/go/lexer"
	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Deserializer struct {
	lexer.Lexer

	ExpectComma bool
}

func UnescapeJSONString(s string) string {
	t := trace.Begin("")

	var buf bytes.Buffer

	bs := strings.FindChar(s, '\\')
	if bs == -1 {
		return s
	}
	buf.WriteString(s[:bs])
	s = s[bs+1:]

	for {
		bs := strings.FindChar(s, '\\')
		if bs == -1 {
			break
		}
		buf.WriteString(s[:bs])
		s = s[bs+1:]
	}
	buf.WriteString(s)

	trace.End(t)
	return buf.String()
}

func (d *Deserializer) Init(buf []byte) {
	t := trace.Begin("")

	d.Reset()
	d.Lexer.Tokens = d.Lexer.Tokens[:0]
	d.Lexer.FileSet = token.NewFileSet()
	d.Lexer.Scanner.Init(d.Lexer.FileSet.AddFile("json", d.Lexer.FileSet.Base(), len(buf)), buf, nil, 0)

	trace.End(t)
}

func (d *Deserializer) GetObjectBegin() bool {
	t := trace.Begin("")

	ret := d.Lexer.ParseToken(token.LBRACE)

	trace.End(t)
	return ret
}

func (d *Deserializer) GetObjectEnd() bool {
	t := trace.Begin("")

	ret := d.Lexer.ParseToken(token.RBRACE)

	trace.End(t)
	return ret
}

func (d *Deserializer) GetArrayBegin() bool {
	t := trace.Begin("")

	ret := d.Lexer.ParseToken(token.LBRACK)

	trace.End(t)
	return ret
}

func (d *Deserializer) GetArrayEnd() bool {
	t := trace.Begin("")

	ret := d.Lexer.ParseToken(token.RBRACK)

	trace.End(t)
	return ret
}

func (d *Deserializer) GetSliceBegin(n *int) bool {
	t := trace.Begin("")

	var notEmpty bool
	var brackets int
	var braces int

	if !d.GetArrayBegin() {
		trace.End(t)
		return false
	}

	pos := d.Lexer.Position

	brackets = 1
	for brackets > 0 {
		switch d.Lexer.Curr().GoToken {
		default:
			notEmpty = true
		case token.LBRACE:
			braces++
		case token.LBRACK:
			brackets++
		case token.RBRACE:
			braces--
		case token.RBRACK:
			brackets--
		case token.COMMA:
			if (braces == 0) && (brackets == 1) {
				*n++
			}
		}
		d.Lexer.Next()
	}

	d.Lexer.Position = pos
	if notEmpty {
		*n++
	}

	trace.End(t)
	return true
}

func (d *Deserializer) GetSliceEnd() bool {
	return d.GetArrayEnd()
}

func (d *Deserializer) GetComma() bool {
	t := trace.Begin("")

	if d.ExpectComma {
		if !d.Lexer.ParseToken(token.COMMA) {
			trace.End(t)
			return false
		}
		d.ExpectComma = false
	}

	trace.End(t)
	return true
}

func (d *Deserializer) GetKey(key *string) bool {
	t := trace.Begin("")

	tok := d.Lexer.Curr()
	if tok.GoToken == token.RBRACE || tok.GoToken == token.RBRACK {
		trace.End(t)
		return false
	}

	if !d.GetComma() {
		trace.End(t)
		return false
	}
	if !d.Lexer.ParseStringLit(key) {
		trace.End(t)
		return false
	}
	ret := d.Lexer.ParseToken(token.COLON)

	trace.End(t)
	return ret
}

func (d *Deserializer) GetString(s *string) bool {
	t := trace.Begin("")

	if !d.GetComma() {
		trace.End(t)
		return false
	}

	var lit string
	if !d.Lexer.ParseStringLit(&lit) {
		trace.End(t)
		return false
	}
	*s = UnescapeJSONString(lit)
	d.ExpectComma = true

	trace.End(t)
	return true
}

func (d *Deserializer) GetInt(i *int) bool {
	t := trace.Begin("")

	if !d.GetComma() {
		trace.End(t)
		return false
	}

	if !d.Lexer.ParseIntLit(i) {
		trace.End(t)
		return false
	}
	d.ExpectComma = true

	trace.End(t)
	return true
}

func (d *Deserializer) GetInt32(i *int32) bool {
	t := trace.Begin("")

	if !d.GetComma() {
		trace.End(t)
		return false
	}

	var n int
	if !d.Lexer.ParseIntLit(&n) {
		trace.End(t)
		return false
	}
	*i = int32(n)
	d.ExpectComma = true

	trace.End(t)
	return true
}

func (d *Deserializer) GetUint32(i *uint32) bool {
	t := trace.Begin("")

	if !d.GetComma() {
		trace.End(t)
		return false
	}

	var n int
	if !d.Lexer.ParseIntLit(&n) {
		trace.End(t)
		return false
	}
	*i = uint32(n)
	d.ExpectComma = true

	trace.End(t)
	return true
}

func (d *Deserializer) GetInt64(i *int64) bool {
	t := trace.Begin("")

	if !d.GetComma() {
		trace.End(t)
		return false
	}

	var n int
	if !d.Lexer.ParseIntLit(&n) {
		trace.End(t)
		return false
	}
	*i = int64(n)
	d.ExpectComma = true

	trace.End(t)
	return true
}

func (d *Deserializer) Reset() {
	t := trace.Begin("")

	d.Lexer.Position = 0
	d.Lexer.Error = nil
	d.ExpectComma = false

	trace.End(t)
}
