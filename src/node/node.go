package node

import (
	"github.com/ivanovaleksey/broker/pkg/types"
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

type Node struct {
	ID   types.NodeID
	Type NodeType

	Stop int32

	childMu   sync.RWMutex
	childHash *Node
	childStar *Node
	Next      map[uint64]*Node
}

func NewNode() *Node {
	n := &Node{
		ID: rand.Int63(),
	}
	return n
}

type TraverseNode struct {
	Hash     *Node
	Star     *Node
	Word     *Node
	SelfHash *Node
}

func (n *Node) ChildrenForTraverse(word uint64, withSelfHash bool) TraverseNode {
	// out := GetTraverseNode()
	out := TraverseNode{}
	n.childMu.RLock()
	out.Hash = n.childHash
	out.Star = n.childStar
	out.Word = n.Next[word]
	n.childMu.RUnlock()
	if n.IsHash() && withSelfHash {
		out.SelfHash = n
	}
	return out
}

func (n *Node) Child(part uint64) *Node {
	n.childMu.RLock()
	defer n.childMu.RUnlock()
	switch part {
	case topics.HashStar:
		return n.childStar
	case topics.HashHash:
		return n.childHash
	default:
		return n.Next[part]
	}
}

func (n *Node) SetChild(child *Node, part uint64) {
	n.childMu.Lock()
	switch part {
	case topics.HashStar:
		n.childStar = child
	case topics.HashHash:
		n.childHash = child
	default:
		if n.Next == nil {
			n.Next = make(map[uint64]*Node)
		}
		n.Next[part] = child
	}
	n.childMu.Unlock()
}

func (n *Node) RemoveChild(part uint64) {
	n.childMu.Lock()
	switch part {
	case topics.HashStar:
		n.childStar = nil
	case topics.HashHash:
		n.childHash = nil
	default:
		delete(n.Next, part)
	}
	n.childMu.Unlock()
}

func (n *Node) IsStop() bool {
	// n.stopMu.RLock()
	// defer n.stopMu.RUnlock()
	// return n.Stop
	return atomic.LoadInt32(&n.Stop) == 1
}

func (n *Node) SetStop(value bool) {
	// n.stopMu.Lock()
	// n.Stop = value
	// n.stopMu.Unlock()
	var v int32
	if value {
		v = 1
	}
	atomic.StoreInt32(&n.Stop, v)
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

func (n *Node) Reset() {
	n.childMu.Lock()
	n.childStar = nil
	n.childHash = nil
	n.childMu.Unlock()
}
