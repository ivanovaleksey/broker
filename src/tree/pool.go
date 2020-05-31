package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
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

func GetNodeFromPool() *node.Node {
	n := NodePool.Get().(*node.Node)
	return n
}

func PutNodeToPool(n *node.Node) {
	n.Reset()
	NodePool.Put(n)
}

var PoolCnt int32
var Pool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&PoolCnt, 1)
		return make([]types.ConsumerID, 3)
	},
}

func GetFromPool() []types.ConsumerID {
	return Pool.Get().([]types.ConsumerID)
}

func PutToPool(sl []types.ConsumerID) {
	sl = sl[:0]
	Pool.Put(sl)
}
