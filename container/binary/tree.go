package binary

import (
	"unsafe"

	"github.com/anton2920/gofa/pointers"
)

type (
	TreeKeyPointer    unsafe.Pointer
	TreeKeyComparator func(TreeKeyPointer, TreeKeyPointer) int
)

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
}

func (n *TreeNode) KeyPointer() TreeKeyPointer {
	return (TreeKeyPointer)(pointers.Add(unsafe.Pointer(n), unsafe.Sizeof(*n)))
}

type SearchTree struct {
	N int

	Root *TreeNode

	Comparator TreeKeyComparator
}

func (bt *SearchTree) addAsAChildOf(node *TreeNode, pparent **TreeNode) {
	for {
		if *pparent == nil {
			*pparent = node
			bt.N++
			return
		}
		parent := *pparent

		res := bt.Comparator(node.KeyPointer(), parent.KeyPointer())
		switch {
		case res < 0:
			pparent = &parent.Left
		case res > 0:
			pparent = &parent.Right
		case res == 0:
			return
		}
	}
}

func (bt *SearchTree) Add(node *TreeNode) {
	bt.addAsAChildOf(node, &bt.Root)
}

func (bt *SearchTree) getFromChildOf(key TreeKeyPointer, node *TreeNode) *TreeNode {
	for {
		if node == nil {
			return nil
		}

		res := bt.Comparator(key, node.KeyPointer())
		switch {
		case res < 0:
			node = node.Left
		case res > 0:
			node = node.Right
		case res == 0:
			return node
		}
	}
}

func (bt *SearchTree) Get(key TreeKeyPointer) *TreeNode {
	return bt.getFromChildOf(key, bt.Root)
}

func (bt *SearchTree) Has(key TreeKeyPointer) bool {
	return bt.Get(key) != nil
}

func (bt *SearchTree) delFromChildOf(key TreeKeyPointer, pnode **TreeNode) *TreeNode {
	for {
		if *pnode == nil {
			return nil
		}
		node := *pnode

		res := bt.Comparator(key, node.KeyPointer())
		switch {
		case res < 0:
			return bt.delFromChildOf(key, &node.Left)
		case res > 0:
			return bt.delFromChildOf(key, &node.Right)
		case res == 0:
			if node.Left == nil {
				*pnode = node.Right
			} else if node.Right == nil {
				*pnode = node.Left
			} else {
				parent := node
				it := parent.Left
				for it.Right != nil {
					parent = it
					it = it.Right
				}

				if parent == node {
					parent.Left = it.Left
				} else {
					parent.Right = it.Left
				}

				it.Left = node.Left
				it.Right = node.Right
				*pnode = it
			}

			bt.N--
			node.Left, node.Right = nil, nil
			return node
		}
	}
}

func (bt *SearchTree) Del(key TreeKeyPointer) *TreeNode {
	return bt.delFromChildOf(key, &bt.Root)
}
