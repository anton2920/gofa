package http

import (
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/pointers"
)

type ConnPoolItem struct {
	Item Conn
	Next *ConnPoolItem
}

type ConnPool struct {
	sync.Mutex

	Items []ConnPoolItem
	Head  *ConnPoolItem
}

func NewConnPool(n int) *ConnPool {
	var p ConnPool

	p.Items = make([]ConnPoolItem, n)
	p.Head = nil
	p.PutAll()

	return &p
}

func (p *ConnPool) Get() (*Conn, error) {
	p.Lock()

	item := p.Head
	if item == nil {
		p.Unlock()
		return nil, errors.New("no space left")
	}
	p.Head = item.Next

	p.Unlock()
	return &item.Item, nil
}

func (p *ConnPool) Put(t *Conn) {
	if t == nil {
		return
	}
	if uintptr(unsafe.Pointer(t)) < uintptr(unsafe.Pointer(&p.Items[0].Item)) || uintptr(unsafe.Pointer(t)) > uintptr(unsafe.Pointer(&p.Items[len(p.Items)-1].Item)) {
		/* NOTE(anton2920): do nothing as pointer does not come from the pool. */
		return
	}

	p.Lock()

	item := (*ConnPoolItem)(pointers.Add(unsafe.Pointer(t), -int(unsafe.Sizeof(p.Head))))
	item.Next = p.Head
	p.Head = item

	p.Unlock()
}

func (p *ConnPool) PutAll() {
	p.Lock()

	for i := len(p.Items) - 1; i >= 0; i-- {
		item := &p.Items[i]
		item.Next = p.Head
		p.Head = item
	}

	p.Unlock()
}
