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

	// todo: remove, just for debug
	Part string

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

func (n *Node) ChildrenForTraverse(word string, withSelfHash bool) []*Node {
	n.childMu.RLock()
	wordChild := n.Next[word]
	hashChild := n.Next[NodeHash]
	starChild := n.Next[NodeStar]
	n.childMu.RUnlock()
	out := make([]*Node, 0, 4)
	out = append(out, wordChild, hashChild, starChild)
	// if n.IsHash() {
	// 	out[3] = n
	// }
	if n.IsHash() && withSelfHash {
		out = append(out, n)
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
