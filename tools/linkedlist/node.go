package linkedlist

type LinkedListNode[T any] struct {
	list *LinkedList[T]

	prev *LinkedListNode[T]
	next *LinkedListNode[T]

	Item T
}

func (n *LinkedListNode[T]) findItemNext(predicate func(a T) bool) *LinkedListNode[T] {
	for node := n; node != nil; node = node.next {
		if predicate(node.Item) {
			return node
		}
	}
	return nil
}

// func (n *LinkedListNode[T]) findItemPrev(predicate func(a T) bool) *LinkedListNode[T] {
// 	for node := n; node != nil; node = node.prev {
// 		if predicate(node.Item) {
// 			return node
// 		}
// 	}
// 	return nil
// }

// func (n *LinkedListNode[T]) findItem(predicate func(a T) bool) *LinkedListNode[T] {
// 	found := n.findItemPrev(predicate)
// 	if found != nil {
// 		return found
// 	}
// 	return n.findItemNext(predicate)
// }

func (n *LinkedListNode[T]) Remove() {
	n.list.lck.Lock()
	defer n.list.lck.Unlock()

	if n.prev != nil {
		n.prev.next = n.next
	}
	if n.next != nil {
		n.next.prev = n.prev
	}
}
