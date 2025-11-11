package wire

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/bytes"
	"github.com/anton2920/gofa/trace"
)

type Deserializer struct {
	Buffer []byte
}

func (d *Deserializer) Begin() int {
	t := trace.Begin("")

	version := int(d.Int32())

	trace.End(t)
	return version
}

func (d *Deserializer) Int32() int32 {
	return int32(d.Uint32())
}

func (d *Deserializer) Uint32() uint32 {
	t := trace.Begin("")

	var n uint32
	n |= uint32(d.Buffer[0]) << 0
	n |= uint32(d.Buffer[1]) << 8
	n |= uint32(d.Buffer[2]) << 16
	n |= uint32(d.Buffer[3]) << 24
	d.Buffer = d.Buffer[unsafe.Sizeof(n):]

	trace.End(t)
	return n
}

func (d *Deserializer) String() string {
	t := trace.Begin("")

	l := int(d.Int32())
	s := bytes.AsString(d.Buffer[:l])
	d.Buffer = d.Buffer[l:]

	trace.End(t)
	return s
}

func (d *Deserializer) End() error {
	t := trace.Begin("")

	/* TODO(anton2920): check for CRC. */
	if len(d.Buffer) > 0 {
		return fmt.Errorf("%d bytes are left unconsumed", len(d.Buffer))
	}

	trace.End(t)
	return nil
}
