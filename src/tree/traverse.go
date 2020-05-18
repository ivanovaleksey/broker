package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
)

func (t *Tree) GetConsumers(parts []string) []types.ConsumerID {
	nodeIDs := t.traverse(parts)

	// todo: think about duplicates!!!

	out := make([]types.ConsumerID, 0, len(nodeIDs))
	for _, nodeID := range nodeIDs {
		// todo: consider using bulk method to avoid multiple waits on lock
		ids := t.nodeConsumers.GetConsumers(nodeID)
		out = append(out, ids...)
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

	for _, child := range node.ChildrenForTraverse(part) {
		if child == nil {
			continue
		}
		t.traverseNode(child, parts[1:], stopNodes)
	}

	if node.IsHash() && node.IsStop() {
		stopNodes[node.ID] = struct{}{}
	}
}
