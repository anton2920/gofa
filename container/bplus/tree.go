package bplus

import (
	"bytes"
	"fmt"
	"sync"
	"unsafe"

	"github.com/anton2920/gofa/errors"
	"github.com/anton2920/gofa/trace"
)

/* Tree is an implementation of a B+tree. */
type Tree struct {
	sync.RWMutex
	Pager
	Meta

	SearchPath []TreePathItem
}

type TreeForwardIterator struct {
	*Tree
	Leaf
	Current int
}

type TreePathItem struct {
	Page
	Index int64
	Pos   int
}

const (
	//TreeMaxOrder = 1 << 8
	TreeMaxOrder = 5

	TreeMagic   = 0xFAFEFAAF
	TreeVersion = 0x1
)

func duplicate(buffer []byte, x []byte) []byte {
	if len(buffer) < len(x) {
		panic("insufficient space in buffer")
	}
	return buffer[:copy(buffer, x)]
}

func int2Slice(x int) []byte {
	defer trace.End(trace.Begin(""))

	buf := make([]byte, unsafe.Sizeof(x))
	*(*uint64)(unsafe.Pointer(&buf[0])) = uint64(x)
	return buf
}

func slice2Int(buf []byte) int {
	defer trace.End(trace.Begin(""))

	return int(*(*int)(unsafe.Pointer(&buf[0])))
}

func GetTreeAt(pager Pager, index int64) (*Tree, error) {
	defer trace.End(trace.Begin(""))

	var t Tree

	t.Pager = pager

	base, err := t.ReadPageAt(t.Meta.Page(), index)
	if err != nil {
		const (
			Meta = iota
			Root
			End
			Count
		)

		var pages [Count]Page

		pages[Meta].Init(PageTypeMeta)
		meta := pages[Meta].Meta()
		meta.Magic = TreeMagic
		meta.Version = TreeVersion
		meta.Root = base + Root
		meta.EndSentinel = base + End

		pages[Root].Init(PageTypeLeaf)
		leaf := pages[Root].Leaf()
		leaf.Next = base + End

		pages[End].Init(PageTypeLeaf)

		if _, err := t.Pager.WritePagesAt(pages[:], base); err != nil {
			return nil, fmt.Errorf("failed to write initial pages: %v", err)
		}

		t.Meta.Magic = meta.Magic
		t.Meta.Version = meta.Version
		t.Meta.Root = meta.Root
		t.Meta.EndSentinel = meta.EndSentinel
	}

	if t.Meta.Magic != TreeMagic {
		return nil, fmt.Errorf("wrong tree magic: %d != %d", TreeMagic, t.Meta.Magic)
	}

	return &t, nil

}

func (it *TreeForwardIterator) Next() bool {
	it.Current++
	if it.Current >= int(it.Leaf.N) {
		if it.Leaf.Next == it.Meta.EndSentinel {
			return false
		}
		if _, err := it.ReadPageAt(it.Leaf.Page(), it.Leaf.Next); err != nil {
			return false
		}
		it.Current = 0
	}
	return true
}

func (it *TreeForwardIterator) Key() []byte {
	return it.Leaf.GetKeyAt(it.Current)
}

func (it *TreeForwardIterator) Value() []byte {
	return it.Leaf.GetValueAt(it.Current)
}

func (t *Tree) ReadPageAt(page *Page, index int64) (int64, error) {
	return t.Pager.ReadPagesAt(Page2Slice(page), index)
}

func (t *Tree) WritePageAt(page *Page, index int64) (int64, error) {
	return t.Pager.WritePagesAt(Page2Slice(page), index)
}

func (t *Tree) Begin() (*TreeForwardIterator, error) {
	var it TreeForwardIterator
	var page Page

	it.Tree = t

	index := t.Meta.Root
	for index != 0 {
		if _, err := t.ReadPageAt(&page, index); err != nil {
			return nil, fmt.Errorf("failed to read page: %v", err)
		}

		switch page.Type() {
		case PageTypeNode:
			node := page.Node()
			index = node.GetChildAt(-1)
		case PageTypeLeaf:
			it.Current = -1
			it.Leaf = *page.Leaf()
			return &it, nil
		}
	}

	panic("unreachable")
}

func (t *Tree) Get(key []byte) ([]byte, error) {
	defer trace.End(trace.Begin(""))

	var buffer []byte
	var page Page
	var v []byte

	index := t.Meta.Root
	for index != 0 {
		if _, err := t.ReadPageAt(&page, index); err != nil {
			return nil, fmt.Errorf("failed to read page: %v", err)
		}

		switch page.Type() {
		case PageTypeNode:
			node := page.Node()
			pos := node.Find(key)
			index = node.GetChildAt(pos)
		case PageTypeLeaf:
			index = 0
			leaf := page.Leaf()
			pos, ok := leaf.Find(key)
			if ok {
				v = leaf.GetValueAt(pos + 1)
				switch ValueGetType(v) {
				default:
					panic("unknown value type")
				case ValueTypeFull:
					return ValueGetFull(v), nil
				case ValueTypePartial:
					v, next := ValueGetPartial(v)
					buffer = append(buffer, v...)

					for next != 0 {
						if _, err := t.ReadPageAt(&page, next); err != nil {
							return nil, fmt.Errorf("failed to read page: %v", err)
						}
						overflow := page.Overflow()
						buffer = append(buffer, overflow.GetValue()...)
						next = overflow.Next
					}

					return buffer, nil
				}
			}
		}
	}

	return nil, nil
}

func (t *Tree) Del(key []byte) error {
	return errors.New("not implemented")
}

func (t *Tree) Has(key []byte) (bool, error) {
	defer trace.End(trace.Begin(""))

	var page Page

	offset := t.Meta.Root
	for offset != 0 {
		if _, err := t.ReadPageAt(&page, offset); err != nil {
			return false, fmt.Errorf("failed to read page: %v", err)
		}

		switch page.Type() {
		case PageTypeNode:
			node := page.Node()
			index := node.Find(key)
			offset = node.GetChildAt(index)
		case PageTypeLeaf:
			leaf := page.Leaf()
			_, ok := leaf.Find(key)
			return ok, nil
		}
	}

	return false, nil
}

func (t *Tree) Set(key []byte, value []byte) error {
	defer trace.End(trace.Begin(""))

	var page Page

	var err error
	var ok bool
	var pos int

	t.SearchPath = t.SearchPath[:0]

	index := t.Meta.Root
forIndex:
	for index != 0 {
		if _, err := t.ReadPageAt(&page, index); err != nil {
			return fmt.Errorf("failed to read page: %v", err)
		}

		switch page.Type() {
		case PageTypeNode:
			node := page.Node()
			pos = node.Find(key)
			t.SearchPath = append(t.SearchPath, TreePathItem{page, index, pos})
			index = node.GetChildAt(pos)
		case PageTypeLeaf:
			leaf := page.Leaf()
			pos, ok = leaf.Find(key)
			break forIndex
		}
	}

	var overflow bool
	leaf := page.Leaf()

	if leaf.OverflowAfterInsertKeyValueInEmpty(len(key), FullValueLen(value)) {
		var page Page
		page.Init(PageTypeOverflow)
		overflow := page.Overflow()

		value = overflow.SetValue(value)
		index, err := t.WritePageAt(&page, -1)
		if err != nil {
			return fmt.Errorf("failed to write new overflow: %v", err)
		}

		for (len(value) != 0) && (leaf.OverflowAfterInsertKeyValueInEmpty(len(key), PartialValueLen(value))) {
			overflow.Next = index
			value = overflow.SetValue(value)
			index, err = t.WritePageAt(&page, -1)
			if err != nil {
				return fmt.Errorf("failed to write new overflow: %v", err)
			}
		}

		/* TODO(anton2920): remove extra memory allocation. */
		value = PartialValue(value, index)
	} else {
		/* TODO(anton2920): remove extra memory allocation. */
		value = FullValue(value)
	}

	if ok {
		/* Found key, check for overflow before updating value. */
		overflow = leaf.OverflowAfterInsertValue(len(value))
	} else {
		/* Check for overflow before inserting new key. */
		overflow = leaf.OverflowAfterInsertKeyValue(len(key), len(value)) || (leaf.N >= TreeMaxOrder-1)
	}

	if !overflow {
		if ok {
			/* Updating value for existing key. */
			leaf.SetValueAt(value, pos+1)
		} else {
			/* Insering new key-value. */
			leaf.InsertKeyValueAt(key, value, pos+1)
		}
		if _, err = t.WritePageAt(&page, index); err != nil {
			return fmt.Errorf("failed to write updated leaf: %v", err)
		}
		return nil
	}

	/* Split leaf into two. */
	var newLeaf Leaf
	newLeaf.Page().Init(PageTypeLeaf)
	newBuffer := make([]byte, PageSize)

	half := int(leaf.N) / 2
	if pos < half-1 {
		leaf.MoveData(&newLeaf, 0, half-1, -1)
		if ok {
			leaf.SetValueAt(value, pos)
		} else {
			leaf.InsertKeyValueAt(key, value, pos+1)
		}
	} else {
		leaf.MoveData(&newLeaf, 0, half, -1)
		if ok {
			newLeaf.SetValueAt(value, pos-half)
		} else {
			newLeaf.InsertKeyValueAt(key, value, pos+1-half)
		}
	}

	newLeaf.Next = leaf.Next
	newKey := duplicate(newBuffer, newLeaf.GetKeyAt(0))
	newPage, err := t.WritePageAt(newLeaf.Page(), -1)
	if err != nil {
		return fmt.Errorf("failed to write new leaf: %v", err)
	}

	leaf.Next = newPage
	if _, err = t.WritePageAt(&page, index); err != nil {
		return fmt.Errorf("failed to write updated leaf: %v", err)
	}

	/* Update posing structure. */
	for p := len(t.SearchPath) - 1; p >= 0; p-- {
		page := t.SearchPath[p].Page
		pos := t.SearchPath[p].Pos
		node := page.Node()

		node.SetChildAt(index, pos)

		overflow = node.OverflowAfterInsertKeyChild(len(key)) || (node.N >= TreeMaxOrder-1)
		if !overflow {
			node.InsertKeyChildAt(newKey, newPage, pos+1)
			if _, err = t.WritePageAt(&page, t.SearchPath[p].Index); err != nil {
				return fmt.Errorf("failed to write updated node: %v", err)
			}
			return nil
		}

		var insertKey []byte
		var newNode Page
		newNode.Init(PageTypeNode)

		insertBuffer := make([]byte, PageSize)

		half = int(node.N) / 2
		if pos < half-1 {
			insertKey = duplicate(insertBuffer, newKey)
			newKey = duplicate(newBuffer, node.GetKeyAt(half-1))

			node.MoveData(newNode.Node(), -1, half-1, -1)
			node.InsertKeyChildAt(insertKey, newPage, pos+1)
		} else if pos == half-1 {
			insertKey = duplicate(insertBuffer, node.GetKeyAt(half))
			insertPage := node.GetChildAt(half)

			node.MoveData(newNode.Node(), -1, half, -1)
			newNode.Node().SetChildAt(newPage, -1)
			newNode.Node().InsertKeyChildAt(insertKey, insertPage, pos+1-half)
		} else {
			insertKey = duplicate(insertBuffer, newKey)
			newKey = duplicate(newBuffer, node.GetKeyAt(half))

			node.MoveData(newNode.Node(), -1, half, -1)
			newNode.Node().InsertKeyChildAt(insertKey, newPage, pos-half)
		}

		newPage, err = t.WritePageAt(&newNode, -1)
		if err != nil {
			return fmt.Errorf("failed to write new node: %v", err)
		}

		index, err = t.WritePageAt(&page, t.SearchPath[p].Index)
		if err != nil {
			return fmt.Errorf("failed to write updated node: %v", err)
		}
	}

	var root Page
	root.Init(PageTypeNode)
	node := root.Node()
	node.Init(newKey, t.Meta.Root, newPage)

	t.Meta.Root, err = t.WritePageAt(&root, -1)
	if err != nil {
		return fmt.Errorf("failed to write new root: %v", err)
	}

	return nil
}

func (t *Tree) stringImpl(buf *bytes.Buffer, index int64, level int) error {
	var page Page

	if index != 0 {
		if _, err := t.ReadPageAt(&page, index); err != nil {
			return fmt.Errorf("failed to read page: %v", err)
		}

		for i := 0; i < level; i++ {
			buf.WriteRune('\t')
		}

		switch page.Type() {
		case PageTypeNode:
			node := page.Node()
			for i := 0; i < int(node.N); i++ {
				fmt.Fprintf(buf, "%4d", slice2Int(node.GetKeyAt(i)))
			}
			buf.WriteRune('\n')

			for i := -1; i < int(node.N); i++ {
				t.stringImpl(buf, node.GetChildAt(i), level+1)
			}
		case PageTypeLeaf:
			leaf := page.Leaf()
			for i := 0; i < int(leaf.N); i++ {
				fmt.Fprintf(buf, "%4d", slice2Int(leaf.GetKeyAt(i)))
			}
			buf.WriteRune('\n')
		}
	}

	return nil
}

func (t *Tree) String() string {
	var buf bytes.Buffer

	if err := t.stringImpl(&buf, t.Meta.Root, 0); err != nil {
		return err.Error()
	}

	return buf.String()
}
