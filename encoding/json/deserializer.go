package json

import (
	"bytes"

	"github.com/anton2920/gofa/strings"
	"github.com/anton2920/gofa/trace"
)

type Deserializer struct {
	Error error
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

func (d *Deserializer) Begin() bool {
	return true
}

func (d *Deserializer) End() bool {
	return true
}

func (d *Deserializer) ObjectBegin() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Key(key *string) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) ObjectEnd() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) ArrayBegin() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Next() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) ArrayEnd() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Bool(b *bool) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Int32(i *int32) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Uint32(i *uint32) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Int64(i *int64) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) String(s *string) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}
