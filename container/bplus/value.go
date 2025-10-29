package bplus

import (
	"encoding/binary"
	"unsafe"
)

type ValueType uint8

const (
	ValueTypeNone = ValueType(iota)
	ValueTypeFull
	ValueTypePartial
)

func FullValue(value []byte) []byte {
	buffer := make([]byte, len(value)+int(unsafe.Sizeof(ValueTypeFull)))
	buffer[0] = byte(ValueTypeFull)
	copy(buffer[unsafe.Sizeof(ValueTypeFull):], value)
	return buffer
}

func FullValueLen(value []byte) int {
	return len(value) + int(unsafe.Sizeof(ValueTypeFull))
}

func PartialValue(value []byte, next int64) []byte {
	buffer := make([]byte, len(value)+int(unsafe.Sizeof(ValueTypePartial))+int(unsafe.Sizeof(next)))
	buffer[0] = byte(ValueTypePartial)
	binary.LittleEndian.PutUint64(buffer[unsafe.Sizeof(ValueTypePartial):], uint64(next))
	copy(buffer[unsafe.Sizeof(ValueTypePartial)+unsafe.Sizeof(next):], value)
	return buffer
}

func PartialValueLen(value []byte) int {
	return len(value) + int(unsafe.Sizeof(ValueTypePartial)) + int(unsafe.Sizeof(int64(0)))
}

func ValueGetType(value []byte) ValueType {
	return ValueType(value[0])
}

func ValueGetFull(value []byte) []byte {
	return value[unsafe.Sizeof(ValueTypeFull):]
}

func ValueGetNext(value []byte) int64 {
	return int64(binary.LittleEndian.Uint64(value[unsafe.Sizeof(ValueTypePartial):]))
}

func ValueGetPartial(value []byte) ([]byte, int64) {
	return value[unsafe.Sizeof(ValueTypePartial)+unsafe.Sizeof(int64(0)):], ValueGetNext(value)
}

func ValueSetNext(value []byte, next int64) {
	binary.LittleEndian.PutUint64(value[unsafe.Sizeof(ValueTypePartial):], uint64(next))
}
