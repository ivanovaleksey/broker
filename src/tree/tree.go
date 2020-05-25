package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
)

type Tree struct {
	root          *Node
	nodeConsumers *ConsumersLog
}

func NewTree() *Tree {
	root := NewNode()
	root.ID = -1

	log := NewConsumersLog()

	t := &Tree{
		root:          root,
		nodeConsumers: log,
	}
	return t
}

// AddSubscription receives already prepared parts
func (t *Tree) AddSubscription(consumerID types.ConsumerID, subscription []string) {
	if t.root == nil {
		return
	}

	var (
		lastPart    bool
		currentNode *Node
		// prevNode    *Node
	)
	currentNode = t.root

	for i := 0; i < len(subscription); i++ {
		part := subscription[i]
		lastPart = i == (len(subscription) - 1)

		childNode := currentNode.Child(part)
		if childNode == nil {
			newNode := NewNode()
			newNode.SetType(part)
			newNode.Part = part

			if lastPart {
				newNode.Stop = true
				t.nodeConsumers.AddConsumer(newNode.ID, consumerID)
				// todo: is it ok or should be done in smarter way?
				if newNode.IsHash() {
					currentNode.SetStop(true)
					// todo: check this is not root node, maybe better check
					if currentNode.ID > 0 {
						t.nodeConsumers.AddConsumer(currentNode.ID, consumerID)
					}
				}
			}

			// if currentNode.IsHash() && prevNode != nil {
			// 	prevNode.SetChild(newNode, part)
			// }
			currentNode.SetChild(newNode, part)
			childNode = newNode
		} else {
			if lastPart {
				// set stop in else-branch to avoid re-setting stop=true for new node
				childNode.SetStop(true)
				t.nodeConsumers.AddConsumer(childNode.ID, consumerID)
			}
		}

		// prevNode = currentNode
		currentNode = childNode
	}
}

// todo: if no consumers left node should not be stop anymore
func (t *Tree) RemoveSubscription(consumerID types.ConsumerID, subscription []string) {

}
