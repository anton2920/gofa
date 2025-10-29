package bplus

import (
	"log"
	"reflect"
	"unsafe"

	"github.com/anton2920/gofa/bools"
)

type PageType uint8

type Page [PageSize]byte

type PageHeader struct {
	Type PageType
	N    uint8
	Head uint16
	Tail uint16
	_    [2]byte
}

const (
	PageSize       = 4096
	PageHeaderSize = unsafe.Sizeof(PageHeader{})
)

const (
	PageTypeNone = PageType(iota)
	PageTypeMeta
	PageTypeNode
	PageTypeLeaf
	PageTypeOverflow
)

/* TODO(anton2920): find the best constant for time-space tradeoff. */
const ExtraOffsetAfter = 16

func init() {
	var p Page
	var m Meta
	var n Node
	var l Leaf
	var o Overflow

	const (
		psize = unsafe.Sizeof(p)
		msize = unsafe.Sizeof(m)
		nsize = unsafe.Sizeof(n)
		lsize = unsafe.Sizeof(l)
		osize = unsafe.Sizeof(o)
	)

	if (psize != msize) || (psize != nsize) || (psize != lsize) || (psize != osize) {
		log.Panicf("[tree]: sizeof(Page) == %d, sizeof(Meta) == %d, sizeof(Node) == %d, sizeof(Leaf) == %d, sizeof(Oveflow) == %d", psize, msize, nsize, lsize, osize)
	}
}

func (p *Page) Init(typ PageType) {
	hdr := p.Header()
	hdr.Type = typ
	hdr.N = 0
	hdr.Head = 0
	hdr.Tail = 0
}

func (p *Page) Header() *PageHeader {
	return (*PageHeader)(unsafe.Pointer(p))
}

func (p *Page) Type() PageType {
	return p.Header().Type
}

func (p *Page) Meta() *Meta {
	hdr := p.Header()
	if hdr.Type != PageTypeMeta {
		log.Panicf("Page has type %d, but tried to use it as '*Meta'", hdr.Type)
	}
	return (*Meta)(unsafe.Pointer(p))
}

func (p *Page) Node() *Node {
	hdr := p.Header()
	if hdr.Type != PageTypeNode {
		log.Panicf("Page has type %d, but tried to use it as '*Node'", hdr.Type)
	}
	return (*Node)(unsafe.Pointer(p))
}

func (p *Page) Leaf() *Leaf {
	hdr := p.Header()
	if hdr.Type != PageTypeLeaf {
		log.Panicf("Page has type %d, but tried to use it as '*Leaf'", hdr.Type)
	}
	return (*Leaf)(unsafe.Pointer(p))
}

func (p *Page) Overflow() *Overflow {
	hdr := p.Header()
	if hdr.Type != PageTypeOverflow {
		log.Panicf("Page has type %d, but tried to use it as '*Overflow'", hdr.Type)
	}
	return (*Overflow)(unsafe.Pointer(p))
}

func GetExtraOffset(n int, count int) int {
	var offset uint16

	if count > 0 {
		return (((count + (n % ExtraOffsetAfter) - 1) / ExtraOffsetAfter) + bools.ToInt((n%ExtraOffsetAfter) == 0)) * int(unsafe.Sizeof(offset)) * ExtraOffsetAfter
	} else {
		return (((-count + (ExtraOffsetAfter - (n % ExtraOffsetAfter))) / ExtraOffsetAfter) - bools.ToInt((n%ExtraOffsetAfter) == 0)) * int(unsafe.Sizeof(offset)) * ExtraOffsetAfter
	}
}

func Page2Slice(p *Page) []Page {
	return *(*[]Page)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(p)), Len: 1, Cap: 1}))
}

func Pages2Bytes(ps []Page) []byte {
	return *(*[]byte)(unsafe.Pointer(&reflect.SliceHeader{Data: uintptr(unsafe.Pointer(&ps[0])), Len: len(ps) * PageSize, Cap: cap(ps) * PageSize}))
}
