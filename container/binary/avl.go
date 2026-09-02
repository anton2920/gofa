package binary

type AVLTree struct {
	N int

	Root *TreeNode

	Comparator TreeKeyComparator
}

func (avl *AVLTree) Add(node *TreeNode) {}

func (avl *AVLTree) Get(key TreeKeyPointer) *TreeNode {
	return nil
}

func (avl *AVLTree) Has(key TreeKeyPointer) bool {
	return avl.Get(key) != nil
}

func (avl *AVLTree) Del(key TreeKeyPointer) *TreeNode {
	return nil
}
