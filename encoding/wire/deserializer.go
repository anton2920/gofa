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

func (d *Deserializer) Int8() int8 {
	return int8(d.Uint8())
}

func (d *Deserializer) Uint8() uint8 {
	t := trace.Begin("")

	var n uint8
	n = uint8(d.Buffer[0]) << 0
	d.Buffer = d.Buffer[unsafe.Sizeof(n):]

	trace.End(t)
	return n
}

func (d *Deserializer) Int16() int16 {
	return int16(d.Uint16())
}

func (d *Deserializer) Uint16() uint16 {
	t := trace.Begin("")

	var n uint16
	n = uint16(d.Buffer[0]) << 0
	n |= uint16(d.Buffer[1]) << 8
	d.Buffer = d.Buffer[unsafe.Sizeof(n):]

	trace.End(t)
	return n
}

func (d *Deserializer) Int32() int32 {
	return int32(d.Uint32())
}

func (d *Deserializer) Uint32() uint32 {
	t := trace.Begin("")

	var n uint32
	n = uint32(d.Buffer[0]) << 0
	n |= uint32(d.Buffer[1]) << 8
	n |= uint32(d.Buffer[2]) << 16
	n |= uint32(d.Buffer[3]) << 24
	d.Buffer = d.Buffer[unsafe.Sizeof(n):]

	trace.End(t)
	return n
}

func (d *Deserializer) Int64() int64 {
	return int64(d.Uint64())
}

func (d *Deserializer) Uint64() uint64 {
	t := trace.Begin("")

	var n uint64
	n = uint64(d.Buffer[0]) << 0
	n |= uint64(d.Buffer[1]) << 8
	n |= uint64(d.Buffer[2]) << 16
	n |= uint64(d.Buffer[3]) << 24
	n |= uint64(d.Buffer[4]) << 32
	n |= uint64(d.Buffer[5]) << 40
	n |= uint64(d.Buffer[6]) << 48
	n |= uint64(d.Buffer[7]) << 56
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

func (d *Deserializer) Bytes() []byte {
	t := trace.Begin("")

	l := int(d.Int32())
	bs := d.Buffer[:l]
	d.Buffer = d.Buffer[l:]

	trace.End(t)
	return bs
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
