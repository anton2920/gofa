package http

import (
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/errors"
)

type ConnPoolItem struct {
	Conn
	Next *ConnPoolItem
}

type ConnPool struct {
	sync.Mutex

	Conns []ConnPoolItem
	Head  *ConnPoolItem
}

func NewConnPool(n int) *ConnPool {
	var p ConnPool

	p.Conns = make([]ConnPoolItem, n)
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
	return &item.Conn, nil
}

func (p *ConnPool) Put(conn *Conn) {
	if (conn == nil) || (uintptr(unsafe.Pointer(conn)) < uintptr(unsafe.Pointer(&p.Conns[0].Conn))) || (uintptr(unsafe.Pointer(conn)) > uintptr(unsafe.Pointer(&p.Conns[len(p.Conns)-1].Conn))) {
		/* Do nothing as pointer does not come from the pool. */
		return
	}
	*conn = Conn{}

	p.Lock()

	item := (*ConnPoolItem)(unsafe.Pointer(conn))
	item.Next = p.Head
	p.Head = item

	p.Unlock()
}

func (p *ConnPool) PutAll() {
	p.Lock()

	for i := len(p.Conns) - 1; i >= 0; i-- {
		item := &p.Conns[i]
		item.Next = p.Head
		p.Head = item
	}

	p.Unlock()
}
