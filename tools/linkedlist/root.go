package linkedlist

import "sync"

type LinkedList[T any] struct {
	root *LinkedListNode[T]
	lck  sync.RWMutex
}

func (l *LinkedList[T]) FindItem(predicate func(a T) bool) *LinkedListNode[T] {
	l.lck.RLock()
	defer l.lck.RUnlock()

	if l.root == nil {
		return nil
	}
	return l.root.findItemNext(predicate)
}

func (l *LinkedList[T]) Append(item T) *LinkedListNode[T] {
	l.lck.Lock()
	defer l.lck.Unlock()

	if l.root == nil {
		entry := &LinkedListNode[T]{
			list: l,
			Item: item,
		}
		l.root = entry
		return entry
	}

	//find tail
	tail := l.root
	for {
		if tail.next == nil {
			break
		}
		tail = tail.next
	}

	newEntry := &LinkedListNode[T]{
		list: l,
		Item: item,
		prev: tail,
	}
	tail.next = newEntry
	return newEntry
}
