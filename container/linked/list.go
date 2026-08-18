package linked

type ListNode struct {
	Prev *ListNode
	Next *ListNode
}

type List struct {
	N int

	First *ListNode
	Last  *ListNode
}

func (l *List) At(index int) *ListNode {
	if index >= l.N {
		return nil
	}

	var it *ListNode
	for it = l.First; (it != nil) && (index > 0); it = it.Next {
		index--
	}

	return it
}

func (l *List) InsertBefore(n *ListNode, index int) {
	n.Next = l.At(index)

	if n.Next == nil {
		n.Prev = l.Last
		l.Last = n
	} else {
		n.Prev = n.Next.Prev
		n.Next.Prev = n
	}

	if n.Prev == nil {
		l.First = n
	} else {
		n.Prev.Next = n
	}

	l.N++
}

func (l *List) Append(n *ListNode) {
	l.InsertBefore(n, l.N)
}

func (l *List) Prepend(n *ListNode) {
	l.InsertBefore(n, 0)
}

func (l *List) Remove(n *ListNode) {
	if n != nil {
		if n.Next != nil {
			n.Next.Prev = n.Prev
		}
		if n.Prev != nil {
			n.Prev.Next = n.Next
		}

		if l.First == n {
			l.First = n.Next
		}
		if l.Last == n {
			l.Last = n.Prev
		}

		l.N--
	}
}
