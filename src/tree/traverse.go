package tree

import (
	"container/list"
	"github.com/ivanovaleksey/broker/pkg/types"
	"sync/atomic"
)

func (t *Tree) GetConsumers(topicHash uint64, parts []uint64) []types.ConsumerID {
	// todo: there may be duplicates, because placeholders may be treated differently
	// e.g., in hash case both 2 and 2-># can stop nodes for the same pattern
	// uniq := make(map[types.ConsumerID]struct{}, len(nodeIDs))

	// topicHash := types.Topic(topic).Hash()
	idx := topicHash % bucketsCountConsumers
	t.staticConsumersLocks[idx].RLock()
	uniq := make(map[types.ConsumerID]struct{}, len(t.staticConsumers[idx][topicHash]))
	for _, consumerID := range t.staticConsumers[idx][topicHash] {
		uniq[consumerID] = struct{}{}
	}
	t.staticConsumersLocks[idx].RUnlock()

	if atomic.LoadInt32(&t.nodeConsumersActive) > 0 {
		nodeIDs := t.traverse(parts)
		for _, nodeID := range nodeIDs {
			// todo: consider using bulk method to avoid multiple waits on lock
			ids := t.nodeConsumers.GetConsumers(nodeID)
			for consumerID := range ids {
				_, ok := uniq[consumerID]
				if ok {
					continue
				}
				uniq[consumerID] = struct{}{}
			}
		}
	}

	out := make([]types.ConsumerID, 0, len(uniq))
	for id := range uniq {
		out = append(out, id)
	}
	return out
}

func (t *Tree) traverse(parts []uint64) []types.NodeID {
	// todo: need this check?
	if len(parts) == 0 {
		return nil
	}
	visited := make(map[types.NodeID]struct{})
	t.traverseQueue(t.root, parts, visited)
	if len(visited) == 0 {
		return nil
	}
	out := make([]types.NodeID, 0, len(visited))
	for nodeID := range visited {
		out = append(out, nodeID)
	}
	return out
}

// todo: for now it doesn't work with subsequent hashes
func (t *Tree) traverseQueue(node *Node, parts []uint64, stopNodes map[types.NodeID]struct{}) {
	type NodeWithParts struct {
		node  *Node
		parts []uint64
	}

	if node == nil {
		panic("node should not be nil")
	}

	queue := list.New()
	queue.PushBack(&NodeWithParts{
		node:  node,
		parts: parts,
	})

	fn := func(hashNode *Node, part uint64, parts []uint64) {
		children := hashNode.ChildrenForTraverse(part, false)
		if n := children.Word; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
		}
		if n := children.Star; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
		}
		if n := children.Hash; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
		}
		PutTraverseNode(children)
	}

	for queue.Len() > 0 {
		front := queue.Front()
		nodePart := front.Value.(*NodeWithParts)
		node := nodePart.node
		parts := nodePart.parts

		if len(parts) == 0 {
			if node.IsStop() {
				stopNodes[node.ID] = struct{}{}
			}
			queue.Remove(front)
			continue
		}

		part := parts[0]

		children := node.ChildrenForTraverse(part, true)
		if n := children.Word; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
		}
		if n := children.Star; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
		}
		if n := children.Hash; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
			fn(n, part, parts)
		}
		if n := children.SelfHash; n != nil {
			queue.PushBack(&NodeWithParts{
				node:  n,
				parts: parts[1:],
			})
			fn(n, part, parts)
		}
		PutTraverseNode(children)

		queue.Remove(front)
	}

	// todo: where to put it?
	// if node.IsHash() && node.IsStop() {
	// 	stopNodes[node.ID] = struct{}{}
	// }
}
