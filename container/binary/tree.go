package binary

type TreeNodeComparator func(*TreeNode, *TreeNode) int

type TreeNode struct {
	Left  *TreeNode
	Right *TreeNode
}

type SearchTree struct {
	N int

	Root *TreeNode
}

func (bt *SearchTree) addAsAChildOf(key *TreeNode, pnode **TreeNode, cmp TreeNodeComparator) {
	if *pnode == nil {
		*pnode = key
		return
	}
	node := *pnode

	res := cmp(key, node)
	switch {
	case res < 0:
		bt.addAsAChildOf(key, &node.Left, cmp)
	case res > 0:
		bt.addAsAChildOf(key, &node.Right, cmp)
	}
}

func (bt *SearchTree) Add(key *TreeNode, cmp TreeNodeComparator) {
	bt.addAsAChildOf(key, &bt.Root, cmp)
	bt.N++
}

func (bt *SearchTree) getFromChildOf(key *TreeNode, node *TreeNode, cmp TreeNodeComparator) *TreeNode {
	if node == nil {
		return nil
	}

	res := cmp(key, node)
	switch {
	case res < 0:
		return bt.getFromChildOf(key, node.Left, cmp)
	case res > 0:
		return bt.getFromChildOf(key, node.Right, cmp)
	}

	return node
}

func (bt *SearchTree) Get(key *TreeNode, cmp TreeNodeComparator) *TreeNode {
	return bt.getFromChildOf(key, bt.Root, cmp)
}

func (bt *SearchTree) Has(key *TreeNode, cmp TreeNodeComparator) bool {
	return bt.Get(key, cmp) != nil
}

func (bt *SearchTree) delFromParent(key *TreeNode, node *TreeNode, parent *TreeNode, cmp TreeNodeComparator) *TreeNode {
	if node == nil {
		return nil
	}

	res := cmp(key, node)
	switch {
	case res < 0:
		return bt.delFromParent(key, node.Left, node, cmp)
	case res > 0:
		return bt.delFromParent(key, node.Right, node, cmp)
	}

	if parent.Left == node {
		parent.Left = node.Left
		if node.Right != nil {
			bt.addAsAChildOf(node.Right, &parent.Right, cmp)
		}
	} else {
		parent.Right = node.Right
		if node.Left != nil {
			bt.addAsAChildOf(node.Left, &parent.Left, cmp)
		}
	}

	node.Left, node.Right = nil, nil
	return node
}

func (bt *SearchTree) Del(key *TreeNode, cmp TreeNodeComparator) *TreeNode {
	n := bt.delFromParent(key, bt.Root, nil, cmp)
	if n != nil {
		bt.N--
	}
	return n
}
