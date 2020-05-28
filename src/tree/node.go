package tree

import (
	"math/rand"
	"sync"
	"sync/atomic"
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

	// stopMu sync.RWMutex
	// Stop   bool
	stop int32

	childMu sync.RWMutex
	// todo: are int hashes more efficient?
	Next map[string]*Node
}

func NewNode() *Node {
	n := &Node{
		ID:   rand.Int63(),
		stop: 0,
		Next: make(map[string]*Node),
	}
	return n
}

var pool = sync.Pool{
	New: func() interface{} {
		return make([]*Node, 0, 4)
	},
}

func (n *Node) ChildrenForTraverse(word string, withSelfHash bool) []*Node {
	n.childMu.RLock()
	wordChild := n.Next[word]
	hashChild := n.Next[NodeHash]
	starChild := n.Next[NodeStar]
	n.childMu.RUnlock()
	// out := make([]*Node, 0, 4)
	out := pool.Get().([]*Node)
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
	// n.stopMu.RLock()
	// defer n.stopMu.RUnlock()
	// return n.Stop
	return atomic.LoadInt32(&n.stop) == 1
}

func (n *Node) SetStop(value bool) {
	// n.stopMu.Lock()
	// n.Stop = value
	// n.stopMu.Unlock()
	var v int32
	if value {
		v = 1
	}
	atomic.StoreInt32(&n.stop, v)
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
