package tree

import (
	"math/rand"
	"sync"
)

const (
	NodeStar = "*"
	NodeHash = "#"
)

type NodeType int

const (
	NodeTypeInvalid NodeType = iota
	NodeTypeStar
	NodeTypeHash
	NodeTypeWord
)

type Node struct {
	ID   int64
	Type NodeType

	stopMu sync.RWMutex
	Stop   bool

	childMu sync.RWMutex
	// todo: are int hashes more efficient?
	Next map[string]*Node
}

func NewNode() *Node {
	n := &Node{
		ID:   rand.Int63(),
		Stop: false,
		Next: make(map[string]*Node),
	}
	return n
}

func (n *Node) ChildrenForTraverse(word string) []*Node {
	n.childMu.RLock()
	out := []*Node{
		n.Next[word],
		n.Next[NodeHash],
		n.Next[NodeStar],
		nil,
	}
	n.childMu.RUnlock()
	if n.IsHash() {
		out[3] = n
	}
	return out
}

func (n *Node) Child(part string) *Node {
	n.childMu.RLock()
	defer n.childMu.RUnlock()
	return n.Next[part]
}

func (n *Node) SetChild(child *Node, part string) {
	n.childMu.Lock()
	n.Next[part] = child
	n.childMu.Unlock()
}

func (n *Node) IsStop() bool {
	n.stopMu.RLock()
	defer n.stopMu.RUnlock()
	return n.Stop
}

func (n *Node) SetStop(value bool) {
	n.stopMu.Lock()
	n.Stop = value
	n.stopMu.Unlock()
}

func (n *Node) SetType(value string) {
	switch value {
	case NodeStar:
		n.Type = NodeTypeStar
	case NodeHash:
		n.Type = NodeTypeHash
	default:
		n.Type = NodeTypeWord
	}
}

func (n *Node) IsHash() bool {
	return n.Type == NodeTypeHash
}
