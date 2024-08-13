package http

import (
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/util"
)

type ContextPoolItem struct {
	Next *ContextPoolItem
	Item Context
}

type ContextPool struct {
	sync.Mutex

	Items []ContextPoolItem
	Head  *ContextPoolItem
}

func NewContextPool(n int) *ContextPool {
	var p ContextPool

	p.Items = make([]ContextPoolItem, n)
	p.Head = nil
	p.PutAll()

	return &p
}

func (p *ContextPool) Get() (*Context, error) {
	p.Lock()

	item := p.Head
	if item == nil {
		p.Unlock()
		return nil, NoSpaceLeft
	}
	p.Head = item.Next

	p.Unlock()
	return &item.Item, nil
}

func (p *ContextPool) Put(t *Context) {
	if t == nil {
		return
	}
	if uintptr(unsafe.Pointer(t)) < uintptr(unsafe.Pointer(&p.Items[0].Item)) || uintptr(unsafe.Pointer(t)) > uintptr(unsafe.Pointer(&p.Items[len(p.Items)-1].Item)) {
		/* NOTE(anton2920): do nothing as pointer does not come from the pool. */
		return
	}

	p.Lock()

	item := (*ContextPoolItem)(util.PtrAdd(unsafe.Pointer(t), -int(unsafe.Sizeof(p.Head))))
	item.Next = p.Head
	p.Head = item

	p.Unlock()
}

func (p *ContextPool) PutAll() {
	p.Lock()

	for i := len(p.Items) - 1; i >= 0; i-- {
		item := &p.Items[i]
		item.Next = p.Head
		p.Head = item
	}

	p.Unlock()
}
