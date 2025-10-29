package bplus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/bools"
	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/trace"
)

type Leaf struct {
	PageHeader

	Next int64

	/* Data is structured as follows: | N*sizeof(uint16) bytes of keyOffsets | keys... | ...empty space... | ...values | N*sizeof(uint16) bytes of valueOffsets | */
	Data [PageSize - PageHeaderSize - 1*unsafe.Sizeof(int64(0))]byte
}

func init() {
	var page Page
	page.Init(PageTypeLeaf)

	leaf := page.Leaf()
	debug.Printf("[leaf]: len(Leaf.Data) == %d\n", len(leaf.Data))

	leaf.InsertKeyValueAt([]byte{1, 2, 3, 4}, []byte{5, 6, 7, 8, 9, 10}, 0)
	debug.Printf("[leaf]: %v\n", leaf)

	leaf.InsertKeyValueAt([]byte{192, 168, 0, 1}, []byte{253, 253, 253, 0}, 1)
	debug.Printf("[leaf]: %v\n", leaf)

	leaf.InsertKeyValueAt([]byte{254}, []byte{254}, 1)
	debug.Printf("[leaf]: %v\n", leaf)

	leaf.SetKeyValueAt([]byte{1, 2, 3}, []byte{4, 5, 6}, 0)
	debug.Printf("[leaf]: %v\n", leaf)

	leaf.SetKeyValueAt([]byte{255, 255, 255, 255, 255}, []byte{255, 255, 255, 255, 255}, 1)
	debug.Printf("[leaf]: %v\n", leaf)
}

func (l *Leaf) Find(key []byte) (int, bool) {
	defer trace.End(trace.Begin(""))

	if l.N == 0 {
		return -1, false
	} else if res := bytes.Compare(key, l.GetKeyAt(int(l.N)-1)); res >= 0 {
		return int(l.N) - 1 - bools.ToInt(res == 0), res == 0
	}

	for i := 0; i < int(l.N); i++ {
		if res := bytes.Compare(key, l.GetKeyAt(i)); res <= 0 {
			return i - 1, res == 0
		}
	}

	return int(l.N) - 1, false
}

func (l *Leaf) GetExtraOffset(count int) int {
	return GetExtraOffset(int(l.N), count)
}

func (l *Leaf) GetFirstKeyOffset() int {
	return GetExtraOffset(0, int(l.N))
}

func (l *Leaf) GetFirstValueOffset() int {
	return len(l.Data) - GetExtraOffset(0, int(l.N))
}

func (l *Leaf) GetKeyAt(index int) []byte {
	offset, length := l.GetKeyOffsetAndLength(index)
	return l.Data[offset : offset+length]
}

func (l *Leaf) GetKeyOffsetAndLength(index int) (offset int, length int) {
	switch {
	case index < int(l.N)-1:
		offset = int(binary.LittleEndian.Uint16(l.Data[l.GetKeyOffsetInData(index):]))
		length = int(binary.LittleEndian.Uint16(l.Data[l.GetKeyOffsetInData(index+1):])) - offset
	case index == int(l.N)-1:
		offset = int(binary.LittleEndian.Uint16(l.Data[l.GetKeyOffsetInData(index):]))
		length = int(l.Head) - offset
	case index > int(l.N)-1:
		offset = int(l.Head)
		length = 0
	}
	return
}

func (l *Leaf) GetKeyOffsetInData(index int) int {
	var keyOffset uint16
	return int(unsafe.Sizeof(keyOffset)) * index
}

func (l *Leaf) GetKeyOffsets() []uint16 {
	return *(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&l.Data[0])), Len: int(l.N), Cap: int(l.N)}))
}

func (l *Leaf) GetValueAt(index int) []byte {
	offset, length := l.GetValueOffsetAndLength(index)
	return l.Data[offset-length : offset]
}

/* value3 | value2 | value1 | value0 | offt3 | offt2 | offt1 | offt0 | */
/* value0 | offt0 | */
func (l *Leaf) GetValueOffsetAndLength(index int) (offset int, length int) {
	switch {
	case index < int(l.N)-1:
		offset = int(binary.LittleEndian.Uint16(l.Data[l.GetValueOffsetInData(index):]))
		length = offset - int(binary.LittleEndian.Uint16(l.Data[l.GetValueOffsetInData(index+1):]))
	case index == int(l.N)-1:
		offset = int(binary.LittleEndian.Uint16(l.Data[l.GetValueOffsetInData(index):]))
		length = offset - (len(l.Data) - int(l.Tail))
	case index > int(l.N)-1:
		offset = len(l.Data) - int(l.Tail)
		length = 0
	}
	return
}

func (l *Leaf) GetValueOffsetInData(index int) int {
	var valueOffset uint16
	return len(l.Data) - int(unsafe.Sizeof(valueOffset))*(index+1)
}

func (l *Leaf) GetValueOffsets() []uint16 {
	var valueOffset uint16

	return *(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&l.Data[0])) + uintptr(len(l.Data)) - unsafe.Sizeof(valueOffset)*uintptr(l.N), Len: int(l.N), Cap: int(l.N)}))
}

func (l *Leaf) InsertKeyValueAt(key []byte, value []byte, index int) {
	if index > int(l.N) {
		panic("index out of range for insert")
	}

	extraOffset := l.GetExtraOffset(1)
	keyOffset, _ := l.GetKeyOffsetAndLength(index)
	valueOffset, _ := l.GetValueOffsetAndLength(index)
	if int(l.Head)+int(l.Tail)+len(key)+len(value)+2*extraOffset > len(l.Data) {
		panic("insert key-value causes overflow")
	}

	keyOffsets := l.GetKeyOffsets()
	valueOffsets := l.GetValueOffsets()
	if extraOffset > 0 {
		for i := 0; i < index; i++ {
			keyOffsets[i] += uint16(extraOffset)
			valueOffsets[int(l.N)-1-i] -= uint16(extraOffset)
		}
	}
	for i := index; i < int(l.N); i++ {
		keyOffsets[i] += uint16(len(key) + extraOffset)
		valueOffsets[int(l.N)-1-i] -= uint16(len(value) + extraOffset)
	}

	copy(l.Data[keyOffset+len(key)+extraOffset:], l.Data[keyOffset:l.Head])
	copy(l.Data[l.GetFirstKeyOffset()+extraOffset:], l.Data[l.GetFirstKeyOffset():keyOffset])
	copy(l.Data[keyOffset+extraOffset:], key)
	copy(l.Data[l.GetKeyOffsetInData(index+1):], l.Data[l.GetKeyOffsetInData(index):l.GetKeyOffsetInData(int(l.N))])
	binary.LittleEndian.PutUint16(l.Data[l.GetKeyOffsetInData(index):], uint16(keyOffset+extraOffset))

	copy(l.Data[len(l.Data)-int(l.Tail)-len(value)-extraOffset:], l.Data[len(l.Data)-int(l.Tail):valueOffset])
	copy(l.Data[valueOffset-extraOffset:], l.Data[valueOffset:l.GetValueOffsetInData(int(l.N)-1)])
	copy(l.Data[valueOffset-len(value)-extraOffset:], value)
	copy(l.Data[l.GetValueOffsetInData(int(l.N)):], l.Data[l.GetValueOffsetInData(int(l.N)-1):l.GetValueOffsetInData(index-1)])
	binary.LittleEndian.PutUint16(l.Data[l.GetValueOffsetInData(index):], uint16(valueOffset-extraOffset))

	l.Head += uint16(len(key) + extraOffset)
	l.Tail += uint16(len(value) + extraOffset)
	l.N++
}

func (src *Leaf) MoveData(dst *Leaf, where int, from int, to int) {
	var keyLengths, valueLengths int

	if where > int(dst.N) {
		panic("move destination index forces sparseness")
	}
	if to == -1 {
		to = int(src.N)
	}
	count := to - from

	/* Bulk insert to 'dst[where:]' from 'src[from:to]'. */
	extraOffset := dst.GetExtraOffset(count)
	whereKeyOffset, _ := dst.GetKeyOffsetAndLength(where)
	whereValueOffset, _ := dst.GetValueOffsetAndLength(where)

	fromKeyOffset, fromKeyLength := src.GetKeyOffsetAndLength(from)
	keyLengths += fromKeyLength

	fromValueOffset, fromValueLength := src.GetValueOffsetAndLength(from)
	valueLengths += fromValueLength

	for i := from + 1; i < to; i++ {
		_, keyLength := src.GetKeyOffsetAndLength(i)
		keyLengths += int(keyLength)

		_, valueLength := src.GetValueOffsetAndLength(i)
		valueLengths += int(valueLength)
	}

	if int(dst.Head)+int(dst.Tail)+keyLengths+valueLengths+extraOffset > len(dst.Data) {
		panic("move data causes overflow")
	}

	keyOffsets := dst.GetKeyOffsets()
	valueOffsets := dst.GetValueOffsets()
	if extraOffset > 0 {
		for i := 0; i < where; i++ {
			keyOffsets[i] += uint16(extraOffset)
			valueOffsets[int(dst.N)-1-i] -= uint16(extraOffset)
		}
	}
	for i := where; i < int(dst.N); i++ {
		keyOffsets[i] += uint16(keyLengths + extraOffset)
		valueOffsets[int(dst.N)-1-i] -= uint16(valueLengths + extraOffset)
	}

	copy(dst.Data[whereKeyOffset+keyLengths+extraOffset:], dst.Data[whereKeyOffset:dst.Head])
	copy(dst.Data[dst.GetFirstKeyOffset()+extraOffset:], dst.Data[dst.GetFirstKeyOffset():whereKeyOffset])
	copy(dst.Data[whereKeyOffset+extraOffset:], src.Data[fromKeyOffset:fromKeyOffset+keyLengths])
	copy(dst.Data[dst.GetKeyOffsetInData(where+count):], dst.Data[dst.GetKeyOffsetInData(where):dst.GetKeyOffsetInData(int(dst.N))])

	offset := uint16(whereKeyOffset + extraOffset)
	binary.LittleEndian.PutUint16(dst.Data[dst.GetKeyOffsetInData(where):], offset)
	w := where + 1
	for i := from; i < to-1; i++ {
		_, keyLength := src.GetKeyOffsetAndLength(i)
		offset += uint16(keyLength)
		binary.LittleEndian.PutUint16(dst.Data[dst.GetKeyOffsetInData(w):], offset)
		w++
	}

	/* value3 | value2 | value1 | value0 | offt3 | offt2 | offt1 | offt0 | */
	copy(dst.Data[len(dst.Data)-int(dst.Tail)-valueLengths-extraOffset:], dst.Data[len(dst.Data)-int(dst.Tail):whereValueOffset])
	copy(dst.Data[whereValueOffset-extraOffset:], dst.Data[whereValueOffset:dst.GetFirstValueOffset()])
	copy(dst.Data[whereValueOffset-valueLengths-extraOffset:], src.Data[fromValueOffset-valueLengths:fromValueOffset])
	copy(dst.Data[dst.GetValueOffsetInData(int(dst.N)+count-1):], dst.Data[dst.GetValueOffsetInData(int(dst.N)-1):dst.GetValueOffsetInData(where-1)])

	offset = uint16(whereValueOffset - extraOffset)
	binary.LittleEndian.PutUint16(dst.Data[dst.GetValueOffsetInData(where):], offset)
	w = where + 1
	for i := from; i < to-1; i++ {
		_, valueLength := src.GetValueOffsetAndLength(i)
		offset -= uint16(valueLength)
		binary.LittleEndian.PutUint16(dst.Data[dst.GetValueOffsetInData(w):], offset)
		w++
	}

	dst.Head += uint16(keyLengths + extraOffset)
	dst.Tail += uint16(valueLengths + extraOffset)
	dst.N += uint8(count)

	/* Bulk remove of 'src[from:to]'.*/
	extraOffset = src.GetExtraOffset(-count)

	keyOffsets = src.GetKeyOffsets()
	valueOffsets = src.GetValueOffsets()
	if extraOffset > 0 {
		for i := 0; i < from; i++ {
			keyOffsets[i] -= uint16(extraOffset)
			valueOffsets[int(src.N)-1-i] += uint16(extraOffset)
		}
	}
	for i := to; i < int(src.N); i++ {
		keyOffsets[i] -= uint16(keyLengths + extraOffset)
		valueOffsets[int(src.N)-1-i] += uint16(valueLengths + extraOffset)
	}

	copy(src.Data[src.GetKeyOffsetInData(from):], src.Data[src.GetKeyOffsetInData(to):src.GetKeyOffsetInData(int(src.N))])
	copy(src.Data[src.GetFirstKeyOffset()-extraOffset:], src.Data[src.GetFirstKeyOffset():fromKeyOffset])
	copy(src.Data[fromKeyOffset-extraOffset:], src.Data[fromKeyOffset+keyLengths:src.Head])

	copy(src.Data[src.GetValueOffsetInData(from):], src.Data[src.GetValueOffsetInData(int(src.N)-1):src.GetValueOffsetInData(to-1)])
	copy(src.Data[fromValueOffset+extraOffset:], src.Data[fromValueOffset:src.GetFirstValueOffset()])
	copy(src.Data[len(src.Data)-int(src.Tail)+valueLengths+extraOffset:], src.Data[len(src.Data)-int(src.Tail):fromValueOffset-valueLengths])

	src.Head -= uint16(keyLengths + extraOffset)
	src.Tail -= uint16(valueLengths + extraOffset)
	src.N -= uint8(count)
}

func (l *Leaf) OverflowAfterInsertKeyValue(keyLength int, valueLength int) bool {
	return (int8(l.N) == ^0) || (int(l.Head)+int(l.Tail)+keyLength+valueLength+2*l.GetExtraOffset(1) > len(l.Data))
}

func (l *Leaf) OverflowAfterInsertKeyValueInEmpty(keyLength int, valueLength int) bool {
	return keyLength+valueLength+2*l.GetExtraOffset(1) > len(l.Data)
}

func (l *Leaf) OverflowAfterInsertValue(valueLength int) bool {
	return (int8(l.N) == ^0) || (int(l.Head)+int(l.Tail)+valueLength+2*l.GetExtraOffset(1) > len(l.Data))
}

func (l *Leaf) SetKeyValueAt(key []byte, value []byte, index int) {
	if (index < 0) || (index >= int(l.N)) {
		panic("leaf index out of range")
	}

	keyOffset, keyLength := l.GetKeyOffsetAndLength(index)
	valueOffset, valueLength := l.GetValueOffsetAndLength(index)
	if int(l.Head)+int(l.Tail)+len(key)+len(value)-keyLength-valueLength > len(l.Data) {
		panic("set key-value causes overflow")
	}

	copy(l.Data[keyOffset+len(key):], l.Data[keyOffset+keyLength:l.Head])
	copy(l.Data[keyOffset:], key)

	/* value3 | value2 | value1 | value0 | offt3 | offt2 | offt1 | offt0 | */
	copy(l.Data[len(l.Data)-int(l.Tail)-len(value)+valueLength:], l.Data[len(l.Data)-int(l.Tail):valueOffset-valueLength])
	copy(l.Data[valueOffset-len(value):], value)

	keyOffsets := l.GetKeyOffsets()
	valueOffsets := l.GetValueOffsets()
	for i := index + 1; i < int(l.N); i++ {
		keyOffsets[i] += uint16(len(key) - keyLength)
		valueOffsets[int(l.N)-1-i] -= uint16(len(value) - valueLength)
	}

	l.Head += uint16(len(key) - keyLength)
	l.Tail += uint16(len(value) - valueLength)
}

func (l *Leaf) SetValueAt(value []byte, index int) {
	if (index < 0) || (index >= int(l.N)) {
		panic("leaf index out of range")
	}

	valueOffset, valueLength := l.GetValueOffsetAndLength(index)
	if int(l.Head)+int(l.Tail)+len(value)-valueLength > len(l.Data) {
		panic("set value causes overflow")
	}

	/* value3 | value2 | value1 | value0 | offt3 | offt2 | offt1 | offt0 | */
	copy(l.Data[len(l.Data)-int(l.Tail)-len(value)+valueLength:], l.Data[len(l.Data)-int(l.Tail):valueOffset-valueLength])
	copy(l.Data[valueOffset-len(value):], value)

	valueOffsets := l.GetValueOffsets()
	for i := index + 1; i < int(l.N); i++ {
		valueOffsets[int(l.N)-1-i] -= uint16(len(value) - valueLength)
	}

	l.Tail += uint16(len(value) - valueLength)
}

func (l *Leaf) String() string {
	var buf bytes.Buffer

	buf.WriteString("{ Keys: [")
	for i := 0; i < int(l.N); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%v", l.GetKeyAt(i))
	}

	buf.WriteString("], Values: [")

	for i := 0; i < int(l.N); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%v", l.GetValueAt(i))
	}

	buf.WriteString("] }")
	return buf.String()
}

func (l *Leaf) Reset() {
	l.N = 0
	l.Head = 0
	l.Tail = 0
}

func (l *Leaf) Page() *Page {
	return (*Page)(unsafe.Pointer(l))
}
