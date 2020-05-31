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
	Part uint64

	Stop int32

	childMu   sync.RWMutex
	childHash *Node
	childStar *Node
	Next      map[uint64]*Node
	Children  []*Node
}

func NewNode() *Node {
	n := &Node{
		ID: rand.Int63(),
	}
	// go func() {
	// 	tc := time.Tick(10 * time.Second)
	// 	for {
	// 		select {
	// 		case <-tc:
	// 			n.childMu.RLock()
	// 			nextLen := len(n.Next)
	// 			childLen := len(n.Children)
	// 			if nextLen > 10 || childLen > 10 {
	// 				var (
	// 					starChild, hashChild int64
	// 				)
	// 				if n.childHash != nil {
	// 					hashChild = n.childHash.ID
	// 				}
	// 				if n.childStar != nil {
	// 					hashChild = n.childStar.ID
	// 				}
	//
	// 				fmt.Printf("node_id=%d,type=%d,hash_child=%d,star_child=%d,len(next)=%d,len(child)=%d\n", n.ID, n.Type, hashChild, starChild, nextLen, childLen)
	// 			}
	// 			n.childMu.RUnlock()
	// 		}
	// 	}
	// }()

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
	if n.IsSpecial() {
		out.Word = n.Next[word]
	} else {
		for i := 0; i < len(n.Children); i++ {
			if n.Children[i].Part == word {
				out.Word = n.Children[i]
				break
			}
		}
	}
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
		if n.IsSpecial() {
			return n.Next[part]
		} else {
			for i := 0; i < len(n.Children); i++ {
				if n.Children[i].Part == part {
					return n.Children[i]
				}
			}
			return nil
		}
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
		// if n.Next == nil {
		if n.IsSpecial() {
			// n.Next = make(map[uint64]*Node)
			n.Next[part] = child
		} else {
			child.Part = part
			n.Children = append(n.Children, child)
		}
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
		if n.IsSpecial() {
			delete(n.Next, part)
		} else {
			for i := 0; i < len(n.Children); i++ {
				if n.Children[i].Part == part {
					n.Children[i] = n.Children[len(n.Children)-1]
					n.Children[len(n.Children)-1] = nil
					n.Children = n.Children[:len(n.Children)-1]
					break
				}
			}
		}
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
	n.Part = 0
	n.childMu.Lock()
	n.childStar = nil
	n.childHash = nil
	for k := range n.Next {
		delete(n.Next, k)
	}
	n.Children = nil
	n.childMu.Unlock()
}

func (n *Node) IsRoot() bool {
	return n.ID == -1
}

func (n *Node) IsSpecial() bool {
	return n.ID < 0
}
