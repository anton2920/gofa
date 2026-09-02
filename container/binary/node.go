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
