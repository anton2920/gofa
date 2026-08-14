package binary

import (
	"unsafe"

	"github.com/anton2920/gofa/mem"
)

type TreeKeyComparator func(*TreeNode, *TreeNode) int

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
}

type SearchTree struct {
	N int

	Root *TreeNode

	Arena             *mem.Arena
	ElementSize       uintptr
	ElementAlignment  uintptr
	ElementComparator TreeKeyComparator
}

func ArenaPushTreeNode(a *mem.Arena, elementSize uintptr, elementAlignment uintptr) *TreeNode {
	alignment := unsafe.Alignof(TreeNode{})
	if elementAlignment > alignment {
		alignment = elementAlignment
	}
	return (*TreeNode)(a.PushSizeWithAlignment(unsafe.Sizeof(TreeNode{})+elementSize, alignment))
}

func (bt *SearchTree) Init(a *mem.Arena, elementSize uintptr, elementAlignment uintptr, cmp TreeKeyComparator) {
	bt.Arena = a
	bt.ElementSize = elementSize
	bt.ElementAlignment = elementAlignment
	bt.ElementComparator = cmp
}

func (bt *SearchTree) NewNode() *TreeNode {
	return ArenaPushTreeNode(bt.Arena, bt.ElementSize, bt.ElementAlignment)
}

func (bt *SearchTree) insertAsAChildOf(n *TreeNode, parent **TreeNode) {
	if *parent == nil {
		*parent = n
		return
	}

	p := *parent
	cmp := bt.ElementComparator(p, n)
	if cmp <= 0 {
		bt.insertAsAChildOf(n, &p.Left)
	} else {
		bt.insertAsAChildOf(n, &p.Right)
	}
}

func (bt *SearchTree) Insert(n *TreeNode) {
	if bt.Root == nil {
		bt.Root = n
	} else {
		bt.insertAsAChildOf(n, &bt.Root)
	}
	bt.N++
}
