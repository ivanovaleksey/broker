package alloc

import (
	"github.com/ivanovaleksey/broker/pkg/hash"
	"github.com/ivanovaleksey/broker/pkg/list"
	"github.com/ivanovaleksey/broker/src/node"
	"github.com/ivanovaleksey/broker/src/tree"
)

func Alloc() {
	// allocHashes()
	// allocNodes()
	// allocListElements()
	// allocTraverseNodes()
	// go func() {
	// 	fn := func() {
	// 		fmt.Printf("list=%d,nodes=%d,hashes=%d\n",
	// 			// atomic.LoadInt32(&node.TraversePoolCount),
	// 			atomic.LoadInt32(&list.ElementPoolCount),
	// 			atomic.LoadInt32(&tree.NodePoolCount),
	// 			atomic.LoadInt32(&hash.PoolCount),
	// 		)
	// 	}
	// 	fn()
	// 	tc := time.Tick(time.Second * 10)
	// 	for {
	// 		select {
	// 		case <-tc:
	// 			fn()
	// 		}
	// 	}
	// }()
}

func allocHashes() {
	{
		const size = 1000
		var items [size]interface{}
		for i := 0; i < size; i++ {
			items[i] = hash.Pool.Get()
		}
		for i := 0; i < size; i++ {
			hash.Pool.Put(items[i])
		}
	}
}

func allocNodes() {
	{
		const size = 1000000
		var items [size]*node.Node
		for i := 0; i < size; i++ {
			items[i] = tree.NodePool.Get().(*node.Node)
		}
		for i := 0; i < size; i++ {
			tree.NodePool.Put(items[i])
		}
	}
}

func allocTraverseNodes() {
	{
		const size = 100000
		var items [size]interface{}
		for i := 0; i < size; i++ {
			items[i] = tree.TraversePool.Get()
		}
		for i := 0; i < size; i++ {
			tree.TraversePool.Put(items[i])
		}
	}
}

func allocListElements() {
	{
		const size = 1000000
		var items [size]*list.Element
		for i := 0; i < size; i++ {
			items[i] = list.ElementPool.Get().(*list.Element)
		}
		for i := 0; i < size; i++ {
			list.ElementPool.Put(items[i])
		}
	}
}
