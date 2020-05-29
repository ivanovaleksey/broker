package tree

import (
	"github.com/ivanovaleksey/broker/src/topics"
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

var traversePool = sync.Pool{
	New: func() interface{} {
		return make([]*Node, 0, 4)
	},
}

var NodePoolCount int32
var NodePool = sync.Pool{
	New: func() interface{} {
		atomic.AddInt32(&NodePoolCount, 1)
		return NewNode()
	},
}

type Node struct {
	ID   int64
	Type NodeType

	// stopMu sync.RWMutex
	// Stop   bool
	stop int32

	childMu sync.RWMutex
	// todo: are int hashes more efficient?
	Next map[uint64]*Node
}

func NewNode() *Node {
	n := &Node{
		ID:   rand.Int63(),
		stop: 0,
		// todo: maps from pool?
		Next: make(map[uint64]*Node),
	}
	return n
}

func (n *Node) ChildrenForTraverse(word uint64, withSelfHash bool) []*Node {
	n.childMu.RLock()
	wordChild := n.Next[word]
	hashChild := n.Next[topics.HashHash]
	starChild := n.Next[topics.HashStar]
	n.childMu.RUnlock()
	// out := make([]*Node, 0, 4)
	out := traversePool.Get().([]*Node)
	out = append(out, wordChild, hashChild, starChild)
	// if n.IsHash() {
	// 	out[3] = n
	// }
	if n.IsHash() && withSelfHash {
		out = append(out, n)
	}
	return out
}

func (n *Node) Child(part uint64) *Node {
	n.childMu.RLock()
	defer n.childMu.RUnlock()
	return n.Next[part]
}

func (n *Node) SetChild(child *Node, part uint64) {
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

func (n *Node) SetType(value uint64) {
	switch value {
	case topics.HashStar:
		n.Type = NodeTypeStar
	case topics.HashHash:
		n.Type = NodeTypeHash
	default:
		n.Type = NodeTypeWord
	}
}

func (n *Node) IsHash() bool {
	return n.Type == NodeTypeHash
}
