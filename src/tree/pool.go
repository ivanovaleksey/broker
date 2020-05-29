package tree

import (
	"sync"
)

var TraversePoolCount int32
var TraversePool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&TraversePoolCount, 1)
		return new(TraverseNode)
	},
}

func GetTraverseNode() *TraverseNode {
	node := TraversePool.Get().(*TraverseNode)
	return node
}

func PutTraverseNode(n *TraverseNode) {
	TraversePool.Put(n)
}

var NodePoolCount int32
var NodePool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&NodePoolCount, 1)
		return NewNode()
	},
}
