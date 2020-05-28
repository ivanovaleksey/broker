package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"strings"
)

func (t *Tree) GetConsumers(parts []string) []types.ConsumerID {
	nodeIDs := t.traverse(parts)

	// todo: there may be duplicates, because placeholders may be treated differently
	// e.g., in hash case both 2 and 2-># can stop nodes for the same pattern
	// uniq := make(map[types.ConsumerID]struct{}, len(nodeIDs))

	topic := strings.Join(parts, ".")
	topicHash := types.Topic(topic).Hash()
	idx := topicHash % bucketsCountConsumers
	t.staticConsumersLocks[idx].RLock()
	uniq := make(map[types.ConsumerID]struct{}, len(t.staticConsumers[idx][topicHash]))
	for _, consumerID := range t.staticConsumers[idx][topicHash] {
		uniq[consumerID] = struct{}{}
	}
	t.staticConsumersLocks[idx].RUnlock()

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

	out := make([]types.ConsumerID, 0, len(uniq))
	for id := range uniq {
		out = append(out, id)
	}
	return out
}

func (t *Tree) traverse(parts []string) []types.NodeID {
	// todo: need this check?
	if len(parts) == 0 {
		return nil
	}
	visited := make(map[types.NodeID]struct{})
	t.traverseNode(t.root, parts, visited)
	if len(visited) == 0 {
		return nil
	}
	out := make([]types.NodeID, 0, len(visited))
	for nodeID := range visited {
		out = append(out, nodeID)
	}
	return out
}

// todo: think about non-recursive variant
// traverseNode returns true if last part is happen to be on a stop-node
func (t *Tree) traverseNode(node *Node, parts []string, stopNodes map[types.NodeID]struct{}) {
	if node == nil {
		panic("node should not be nil")
	}

	if len(parts) == 0 {
		if !node.IsStop() {
			return
		}
		stopNodes[node.ID] = struct{}{}
		return
	}

	part := parts[0]
	// wordChild, hashChild, starChild := node.ChildrenForTraverse(part)
	// children := []*Node{wordChild, hashChild, starChild}
	// if node.IsHash() {
	// 	children = append(children, node)
	// }

	children := node.ChildrenForTraverse(part, true)
	for {
		if len(children) == 0 {
			break
		}
		child := children[0]
		children = children[1:]
		if child == nil {
			continue
		}
		if child.IsHash() {
			hashChildren := child.ChildrenForTraverse(part, false)
			children = append(children, hashChildren...)
		}
		t.traverseNode(child, parts[1:], stopNodes)
	}

	// for _, child := range children {
	// 	if child == nil {
	// 		continue
	// 	}
	// 	if child.IsHash() {
	// 		hashChildren := child.ChildrenForTraverse(part)
	// 		children = append(children, hashChildren...)
	// 	}
	// 	t.traverseNode(child, parts[1:], stopNodes)
	// }

	if node.IsHash() && node.IsStop() {
		stopNodes[node.ID] = struct{}{}
	}
}
