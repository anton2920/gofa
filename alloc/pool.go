package alloc

import (
	"sync"
	"unsafe"
)

type PoolElement[T any] struct {
	Next *PoolElement[T]
	Elem T
}

type Pool[T any] struct {
	Elems []PoolElement[T]
	Head  *PoolElement[T]
}

type SyncPool[T any] struct {
	sync.Mutex
	Pool[T]
}

func NewPool[T any](n int) Pool[T] {
	var p Pool[T]

	p.Elems = make([]PoolElement[T], n)
	p.Head = nil
	p.PutAll()

	return p
}

func (p *Pool[T]) Get() (*T, error) {
	elem := p.Head
	if elem == nil {
		return nil, NoSpaceLeft
	}
	p.Head = elem.Next

	return &elem.Elem, nil
}

func (p *Pool[T]) Put(t *T) {
	if t == nil {
		return
	}
	if (uintptr(unsafe.Pointer(t)) < uintptr(unsafe.Pointer(&p.Elems[0].Elem))) || (uintptr(unsafe.Pointer(t)) > uintptr(unsafe.Pointer(&p.Elems[len(p.Elems)-1].Elem))) {
		panic("elem does not come from pool")
	}

	elem := (*PoolElement[T])(unsafe.Add(unsafe.Pointer(t), -int(unsafe.Sizeof(t))))
	elem.Next = p.Head
	p.Head = elem
}

func (p *Pool[T]) PutAll() {
	for i := len(p.Elems) - 1; i >= 0; i-- {
		elem := &p.Elems[i]
		elem.Next = p.Head
		p.Head = elem
	}
}

func NewSyncPool[T any](n int) SyncPool[T] {
	var p SyncPool[T]

	p.Elems = make([]PoolElement[T], n)
	p.Head = nil
	p.PutAll()

	return p
}

func (p *SyncPool[T]) Get() (*T, error) {
	p.Lock()
	defer p.Unlock()
	return p.Pool.Get()
}

func (p *SyncPool[T]) Put(t *T) {
	p.Lock()
	defer p.Unlock()
	p.Pool.Put(t)
}

func (p *SyncPool[T]) PutAll() {
	p.Lock()
	defer p.Unlock()
	p.Pool.PutAll()
}
