package tree

import (
	"github.com/ivanovaleksey/broker/src/node"
	"sync"
)

var TraversePoolCount int32
var TraversePool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&TraversePoolCount, 1)
		return new(node.TraverseNode)
	},
}

func GetTraverseNode() *node.TraverseNode {
	node := TraversePool.Get().(*node.TraverseNode)
	return node
}

func PutTraverseNode(n *node.TraverseNode) {
	TraversePool.Put(n)
}

var NodePoolCount int32
var NodePool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&NodePoolCount, 1)
		return node.NewNode()
	},
}
