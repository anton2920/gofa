package bplus

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/debug"
	"github.com/anton2920/gofa/trace"
)

type Node struct {
	PageHeader

	/* Data is structured as follows: | N*sizeof(uint16) bytes of keyOffsets | keys... | ...empty space... | N*sizeof(int64) bytes of children | */
	Data [PageSize - PageHeaderSize]byte
}

func init() {
	var page Page
	page.Init(PageTypeNode)

	node := page.Node()
	debug.Printf("[node]: len(Node.Data) == %d\n", len(node.Data))

	node.Init([]byte{1, 2, 3, 4}, 0, 1)
	debug.Printf("[node]: %v\n", node)

	node.InsertKeyChildAt([]byte{5, 6, 7, 8, 9, 10}, 2, 1)
	debug.Printf("[node]: %v\n", node)

	node.InsertKeyChildAt([]byte{192, 168, 0, 1}, 3, 1)
	debug.Printf("[node]: %v\n", node)

	node.SetKeyAt([]byte{254}, 0)
	node.SetChildAt(254, 0)
	debug.Printf("[node]: %v\n", node)

	node.SetKeyAt([]byte{255, 255, 255, 255, 255, 255, 255, 255}, 1)
	node.SetChildAt(255, 1)
	debug.Printf("[node]: %v\n", node)
}

func (n *Node) Init(key []byte, child0 int64, child1 int64) {
	extraOffset := GetExtraOffset(0, 1)

	n.SetChildAt(child0, -1)
	n.SetChildAt(child1, 0)

	n.Head = uint16(extraOffset)
	n.Tail = uint16(unsafe.Sizeof(child0)) * 2
	n.N = 1

	binary.LittleEndian.PutUint16(n.Data[n.GetKeyOffsetInData(0):], uint16(extraOffset))
	n.SetKeyAt(key, 0)
}

func (n *Node) Find(key []byte) int {
	defer trace.End(trace.Begin(""))

	if res := bytes.Compare(key, n.GetKeyAt(int(n.N)-1)); res >= 0 {
		return int(n.N) - 1
	}
	for i := 0; i < int(n.N); i++ {
		if bytes.Compare(key, n.GetKeyAt(i)) < 0 {
			return i - 1
		}
	}
	return int(n.N) - 1
}

func (n *Node) GetChildAt(index int) int64 {
	return int64(binary.LittleEndian.Uint64(n.Data[n.GetChildOffsetInData(index):]))
}

func (n *Node) GetChildOffsetInData(index int) int {
	var i int64
	return len(n.Data) - (index+2)*int(unsafe.Sizeof(i))
}

func (n *Node) GetExtraOffset(count int) int {
	return GetExtraOffset(int(n.N), count)
}

func (n *Node) GetFirstKeyOffset() int {
	return GetExtraOffset(0, int(n.N))
}

func (n *Node) GetKeyAt(index int) []byte {
	offset, length := n.GetKeyOffsetAndLength(index)
	return n.Data[offset : offset+length]
}

func (n *Node) GetKeyOffsetAndLength(index int) (offset int, length int) {
	switch {
	case index < int(n.N)-1:
		offset = int(binary.LittleEndian.Uint16(n.Data[n.GetKeyOffsetInData(index):]))
		length = int(binary.LittleEndian.Uint16(n.Data[n.GetKeyOffsetInData(index+1):])) - offset
	case index == int(n.N)-1:
		offset = int(binary.LittleEndian.Uint16(n.Data[n.GetKeyOffsetInData(index):]))
		length = int(n.Head) - offset
	case index > int(n.N)-1:
		offset = int(n.Head)
		length = 0
	}
	return
}

func (n *Node) GetKeyOffsetInData(index int) int {
	var keyOffset uint16
	return int(unsafe.Sizeof(keyOffset)) * index
}

func (n *Node) GetKeyOffsets() []uint16 {
	return *(*[]uint16)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&n.Data[0])), Len: int(n.N), Cap: int(n.N)}))
}

func (n *Node) InsertKeyChildAt(key []byte, child int64, index int) {
	if (index < 0) || (index > int(n.N)) {
		panic("node index out of range")
	}

	extraOffset := n.GetExtraOffset(1)
	offset, _ := n.GetKeyOffsetAndLength(index)
	if int(n.Head)+int(n.Tail)+len(key)+int(unsafe.Sizeof(child))+extraOffset > len(n.Data) {
		panic("insert key-child causes overflow")
	}

	keyOffsets := n.GetKeyOffsets()
	if extraOffset > 0 {
		for i := 0; i < index; i++ {
			keyOffsets[i] += uint16(extraOffset)
		}
	}
	for i := index; i < int(n.N); i++ {
		keyOffsets[i] += uint16(len(key) + int(extraOffset))
	}

	copy(n.Data[offset+len(key)+extraOffset:], n.Data[offset:n.Head])
	copy(n.Data[n.GetFirstKeyOffset()+extraOffset:], n.Data[n.GetFirstKeyOffset():offset])
	copy(n.Data[offset+int(extraOffset):], key)
	copy(n.Data[n.GetKeyOffsetInData(index+1):], n.Data[n.GetKeyOffsetInData(index):n.GetKeyOffsetInData(int(n.N))])
	binary.LittleEndian.PutUint16(n.Data[n.GetKeyOffsetInData(index):], uint16(offset+int(extraOffset)))

	copy(n.Data[n.GetChildOffsetInData(int(n.N)):], n.Data[n.GetChildOffsetInData(int(n.N)-1):n.GetChildOffsetInData(index-1)])
	n.SetChildAt(child, index)

	n.Head += uint16(len(key) + extraOffset)
	n.Tail += uint16(unsafe.Sizeof(child))
	n.N++
}

func (src *Node) MoveData(dst *Node, where int, from int, to int) {
	var keyLengths int
	var child int64

	if where > int(dst.N) {
		panic("move destination index forces sparseness")
	}
	if to == -1 {
		to = int(src.N)
	}
	count := to - from

	/* Bulk insert keys to 'dst[where+1:]' from 'src[from+1:to]' and children to 'dst[where:]' from 'src[from:to]'. */
	extraOffset := dst.GetExtraOffset(count - 1)
	whereKeyOffset, _ := dst.GetKeyOffsetAndLength(where + 1)

	fromKeyOffset, fromKeyLength := src.GetKeyOffsetAndLength(from + 1)
	keyLengths += fromKeyLength

	for i := from + 2; i < to; i++ {
		_, keyLength := src.GetKeyOffsetAndLength(i)
		keyLengths += int(keyLength)
	}

	childrenLengths := int(unsafe.Sizeof(child)) * count

	if int(dst.Head)+int(dst.Tail)+keyLengths+childrenLengths+extraOffset > len(dst.Data) {
		panic("move data causes overflow")
	}

	keyOffsets := dst.GetKeyOffsets()
	if extraOffset > 0 {
		for i := 0; i < where+1; i++ {
			keyOffsets[i] += uint16(extraOffset)
		}
	}
	for i := where + 1; i < int(dst.N); i++ {
		keyOffsets[i] += uint16(keyLengths + extraOffset)
	}

	copy(dst.Data[whereKeyOffset+keyLengths+extraOffset:], dst.Data[whereKeyOffset:dst.Head])
	copy(dst.Data[dst.GetFirstKeyOffset()+extraOffset:], dst.Data[dst.GetFirstKeyOffset():whereKeyOffset])
	copy(dst.Data[whereKeyOffset+extraOffset:], src.Data[fromKeyOffset:fromKeyOffset+keyLengths])
	copy(dst.Data[dst.GetKeyOffsetInData(where+count):], dst.Data[dst.GetKeyOffsetInData(where+1):dst.GetKeyOffsetInData(int(dst.N))])

	offset := uint16(whereKeyOffset + extraOffset)
	binary.LittleEndian.PutUint16(dst.Data[dst.GetKeyOffsetInData(where+1):], offset)
	w := where + 2
	for i := from + 1; i < to-1; i++ {
		_, keyLength := src.GetKeyOffsetAndLength(i)
		offset += uint16(keyLength)
		binary.LittleEndian.PutUint16(dst.Data[dst.GetKeyOffsetInData(w):], offset)
		w++
	}

	copy(dst.Data[dst.GetChildOffsetInData(int(dst.N)+where-1):], dst.Data[dst.GetChildOffsetInData(int(dst.N)-1):dst.GetChildOffsetInData(where-1)])

	w = where
	for i := from; i < to; i++ {
		dst.SetChildAt(src.GetChildAt(i), w)
		w++
	}

	dst.Head += uint16(keyLengths + extraOffset)
	dst.Tail += uint16(childrenLengths)
	dst.N += uint8(count - 1)

	/* Bulk remove of 'src[from:to]'.*/
	extraOffset = src.GetExtraOffset(-count)

	keyOffsets = src.GetKeyOffsets()
	if extraOffset > 0 {
		for i := 0; i < from; i++ {
			keyOffsets[i] -= uint16(extraOffset)
		}
	}
	for i := to; i < int(src.N); i++ {
		keyOffsets[i] -= uint16(keyLengths + extraOffset)
	}

	keyLengths = 0
	fromKeyOffset, fromKeyLength = src.GetKeyOffsetAndLength(from)
	keyLengths += fromKeyLength

	for i := from + 1; i < to; i++ {
		_, keyLength := src.GetKeyOffsetAndLength(i)
		keyLengths += int(keyLength)
	}

	copy(src.Data[src.GetKeyOffsetInData(from):], src.Data[src.GetKeyOffsetInData(to):src.GetKeyOffsetInData(int(src.N))])
	copy(src.Data[src.GetFirstKeyOffset()-extraOffset:], src.Data[src.GetFirstKeyOffset():fromKeyOffset])
	copy(src.Data[fromKeyOffset-extraOffset:], src.Data[fromKeyOffset+keyLengths:src.Head])

	copy(src.Data[src.GetChildOffsetInData(from):], src.Data[src.GetChildOffsetInData(int(src.N)-1):src.GetChildOffsetInData(to-1)])

	src.Head -= uint16(keyLengths + extraOffset)
	src.Tail -= uint16(childrenLengths)
	src.N -= uint8(count)
}

func (n *Node) OverflowAfterInsertKeyChild(keyLength int) bool {
	var child int64
	return int(n.Head)+int(n.Tail)+keyLength+int(unsafe.Sizeof(child))+n.GetExtraOffset(1) > len(n.Data)
}

func (n *Node) SetChildAt(offset int64, index int) {
	binary.LittleEndian.PutUint64(n.Data[n.GetChildOffsetInData(index):], uint64(offset))
}

func (n *Node) SetKeyAt(key []byte, index int) {
	if (index < 0) || (index >= int(n.N)) {
		panic("node index out of range")
	}

	offset, length := n.GetKeyOffsetAndLength(index)
	if int(n.Head)+int(n.Tail)+len(key)-length > len(n.Data) {
		panic("set key causes overflow")
	}

	keyOffsets := n.GetKeyOffsets()
	for i := index + 1; i < int(n.N); i++ {
		keyOffsets[i] += uint16(len(key) - length)
	}

	/* TODO(anton2920): find the minimum number of bytes so that this key is still distinct from other keys. */
	copy(n.Data[offset+len(key):], n.Data[offset+length:n.Head])
	copy(n.Data[offset:], key)

	n.Head += uint16(len(key) - length)
}

func (n *Node) String() string {
	var buf bytes.Buffer

	buf.WriteString("{ Children: [")
	for i := -1; i < int(n.N); i++ {
		if i > -1 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%d", n.GetChildAt(i))
	}

	buf.WriteString("], Keys: [")
	for i := 0; i < int(n.N); i++ {
		if i > 0 {
			buf.WriteString(", ")
		}
		fmt.Fprintf(&buf, "%v", n.GetKeyAt(i))
	}

	buf.WriteString("] }")
	return buf.String()
}
