package list

import (
	"github.com/ivanovaleksey/broker/src/node"
)

type List struct {
	head *Element
}

type Element struct {
	next *Element

	Node  *node.Node
	Parts []uint64
}

func NewList() List {
	return List{}
}

func (l *List) Push(node *node.Node, parts []uint64) {
	el := GetElementFromPool()
	el.Node = node
	el.Parts = parts

	el.next = l.head
	l.head = el
	// return el
}

func (l *List) Pop() *Element {
	el := l.head
	if el == nil {
		return nil
	}
	l.head = el.next

	return el
}
