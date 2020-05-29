package alloc

import (
	"github.com/ivanovaleksey/broker/pkg/hash"
	"github.com/ivanovaleksey/broker/src/tree"
)

func Alloc() {
	allocHashes()
	allocNodes()
	allocTraverseNodes()
	// go func() {
	// 	fn := func() {
	// 		fmt.Printf("traverse=%d,nodes=%d,hashes=%d\n",
	// 			atomic.LoadInt32(&node.TraversePoolCount),
	// 			atomic.LoadInt32(&node.NodePoolCount),
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
		const size = 100000
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
		var items [size]interface{}
		for i := 0; i < size; i++ {
			items[i] = tree.NodePool.Get()
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
