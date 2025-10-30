package wire

import (
	"fmt"
	"unsafe"

	"github.com/anton2920/gofa/bits"
	"github.com/anton2920/gofa/trace"
)

type Deserializer struct {
	Buffer []byte
	Pos    int

	Error error
}

func init() {
	if unsafe.Sizeof(int32(0)) != unsafe.Sizeof(bits.Flags(0)) {
		panic("int32 != bits.Flags")
	}
}

func (d *Deserializer) Begin(expectedVersion byte) bool {
	t := trace.Begin("")

	var actualVersion byte
	if d.GetByte(&actualVersion) {
		if expectedVersion != actualVersion {
			d.Error = fmt.Errorf("expected version 0x%X, got 0x%X", expectedVersion, actualVersion)
		}
	}

	trace.End(t)
	return d.Error == nil
}

func (d *Deserializer) GetType(expectedType ValueType) bool {
	t := trace.Begin("")

	if d.Error != nil {
		trace.End(t)
		return false
	}

	actualType := ValueType(d.Buffer[d.Pos])
	if expectedType != actualType {
		d.Error = fmt.Errorf("expected value of type %q, got %q (at pos %d)", SerialType2String[expectedType], SerialType2String[actualType], d.Pos)
	}
	d.Pos += int(unsafe.Sizeof(ValueType(0)))

	trace.End(t)
	return d.Error == nil
}

func (d *Deserializer) GetByte(b *byte) bool {
	t := trace.Begin("")

	if d.GetType(ValueTypeByte) {
		*b = d.Buffer[d.Pos]
		d.Pos += int(unsafe.Sizeof(*b))

		trace.End(t)
		return true
	}

	trace.End(t)
	return false
}

func (d *Deserializer) GetInt8(i *int8) bool {
	t := trace.Begin("")

	ok := d.GetByte((*byte)(unsafe.Pointer(i)))

	trace.End(t)
	return ok
}

func (d *Deserializer) GetInt32(i *int32) bool {
	t := trace.Begin("")

	if d.GetType(ValueTypeInt32) {
		*i = int32(d.Buffer[d.Pos+0]) << 0
		*i |= int32(d.Buffer[d.Pos+1]) << 8
		*i |= int32(d.Buffer[d.Pos+2]) << 16
		*i |= int32(d.Buffer[d.Pos+3]) << 24
		d.Pos += int(unsafe.Sizeof(*i))

		trace.End(t)
		return true
	}

	trace.End(t)
	return false
}

func (d *Deserializer) GetFlags(f *bits.Flags) bool {
	t := trace.Begin("")

	ok := d.GetInt32((*int32)(unsafe.Pointer(f)))

	trace.End(t)
	return ok
}

func (d *Deserializer) GetString(s *string) bool {
	t := trace.Begin("")

	if d.GetType(ValueTypeString) {
		var l int32
		if d.GetInt32(&l) {
			/* TODO(anton2920): this allocates memory! */
			*s = string(d.Buffer[d.Pos : d.Pos+int(l)])
			d.Pos += int(l)

			trace.End(t)
			return true
		}
	}

	trace.End(t)
	return false
}

func (d *Deserializer) GetSliceBegin(l *int32) bool {
	t := trace.Begin("")

	if d.GetType(ValueTypeSlice) {
		ok := d.GetInt32(l)

		trace.End(t)
		return ok
	}

	return false
}

func (d *Deserializer) End() bool {
	t := trace.Begin("")

	if d.Error == nil {
		/* TODO(anton2920): add CRC verification after it's added to Serializer. */
		if d.Pos != len(d.Buffer) {
			d.Error = fmt.Errorf("expected to consume %d bytes, consumed only %d", len(d.Buffer), d.Pos)
		}
	}

	trace.End(t)
	return d.Error == nil
}
