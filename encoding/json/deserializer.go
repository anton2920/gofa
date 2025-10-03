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

func (d *Deserializer) Init(buf []byte) {
	t := trace.Begin("")

	trace.End(t)
}

func (d *Deserializer) GetObjectBegin() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetObjectEnd() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetArrayBegin() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetArrayEnd() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetSliceBegin(n *int) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetSliceEnd() bool {
	return d.GetArrayEnd()
}

func (d *Deserializer) GetComma() bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetKey(key *string) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetString(s *string) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetInt(i *int) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetInt32(i *int32) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetUint32(i *uint32) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) GetInt64(i *int64) bool {
	t := trace.Begin("")

	trace.End(t)
	return false
}

func (d *Deserializer) Reset() {
	t := trace.Begin("")

	trace.End(t)
}
