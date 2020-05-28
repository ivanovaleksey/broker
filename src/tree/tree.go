package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"sync"
)

type Tree struct {
	root          *Node
	nodeConsumers *ConsumersLog

	staticConsumersMu sync.RWMutex
	staticConsumers   map[types.Topic]map[types.ConsumerID]struct{}
}

func NewTree() *Tree {
	const (
		rootSize = 1 << 18
		hashSize = 1 << 15
		starSize = 1 << 15
	)

	root := NewNode()
	root.ID = -1
	root.Next = make(map[string]*Node, rootSize)

	star := NewNode()
	star.Type = NodeTypeStar
	star.Part = NodeStar
	star.Next = make(map[string]*Node, starSize)
	root.SetChild(star, NodeStar)

	hash := NewNode()
	hash.Type = NodeTypeHash
	hash.Part = NodeHash
	hash.Next = make(map[string]*Node, hashSize)
	root.SetChild(hash, NodeHash)

	log := NewConsumersLog()

	const staticConsumersSize = 3000000
	staticConsumers := make(map[types.Topic]map[types.ConsumerID]struct{}, staticConsumersSize)

	t := &Tree{
		root:            root,
		nodeConsumers:   log,
		staticConsumers: staticConsumers,
	}

	// go func() {
	// 	fn := func() {
	// 		t.staticConsumersMu.RLock()
	// 		size := len(t.staticConsumers)
	//
	// 		var max, avg, sum int
	// 		for _, m := range t.staticConsumers {
	// 			if len(m) > max {
	// 				max = len(m)
	// 			}
	// 			sum += len(m)
	// 		}
	// 		if size > 0 {
	// 			avg = sum / size
	// 		}
	//
	// 		t.staticConsumersMu.RUnlock()
	// 		fmt.Printf("size=%d, max=%d,avg=%d,sum=%d\n", size, max, avg, sum)
	// 	}
	// 	fn()
	// 	tc := time.Tick(time.Second * 10)
	// 	for {
	// 		select {
	// 		case <-tc:
	// 			fn()
	// 		case <-ctx.Done():
	// 			return
	// 		}
	// 	}
	// }()

	return t
}

func (t *Tree) AddSubscriptionStatic(consumerID types.ConsumerID, topic types.Topic) {
	t.staticConsumersMu.Lock()
	defer t.staticConsumersMu.Unlock()

	inner, ok := t.staticConsumers[topic]
	if !ok {
		inner = make(map[types.ConsumerID]struct{})
		t.staticConsumers[topic] = inner
	}
	inner[consumerID] = struct{}{}
}

// AddSubscription receives already prepared parts
func (t *Tree) AddSubscription(consumerID types.ConsumerID, parts []string) {
	if t.root == nil {
		return
	}

	var (
		lastPart    bool
		currentNode *Node
		// prevNode    *Node
	)
	currentNode = t.root

	fn := func(childNode, parentNode *Node) {
		t.nodeConsumers.AddConsumer(childNode.ID, consumerID)
		// todo: is it ok or should be done in smarter way?
		if childNode.IsHash() {
			// todo: check this is not root node, maybe better check
			if parentNode.ID > 0 {
				// todo: it was outside if, move in case of strange problems
				parentNode.SetStop(true)
				t.nodeConsumers.AddConsumer(parentNode.ID, consumerID)
			}
		}
	}

	for i := 0; i < len(parts); i++ {
		part := parts[i]
		lastPart = i == (len(parts) - 1)

		childNode := currentNode.Child(part)
		if childNode == nil {
			newNode := NewNode()
			newNode.SetType(part)
			newNode.Part = part

			if lastPart {
				newNode.Stop = true
				fn(newNode, currentNode)
				// t.nodeConsumers.AddConsumer(newNode.ID, consumerID)
				// // todo: is it ok or should be done in smarter way?
				// if newNode.IsHash() {
				// 	// todo: check this is not root node, maybe better check
				// 	if currentNode.ID > 0 {
				// 		// todo: it was outside if, move in case of strange problems
				// 		currentNode.SetStop(true)
				// 		t.nodeConsumers.AddConsumer(currentNode.ID, consumerID)
				// 	}
				// }
			}

			// note: it was initial idea, some tests weren't pass, not used for now
			// if currentNode.IsHash() && prevNode != nil {
			// 	prevNode.SetChild(newNode, part)
			// }
			currentNode.SetChild(newNode, part)
			childNode = newNode
		} else {
			if lastPart {
				// set stop in else-branch (not outside of if-else) to avoid re-setting stop=true for new node
				childNode.SetStop(true)
				fn(childNode, currentNode)
				// t.nodeConsumers.AddConsumer(childNode.ID, consumerID)

				// this code is deduplicated in fn
				// if childNode.IsHash() {
				// 	if currentNode.ID > 0 {
				// 		currentNode.SetStop(true)
				// 		t.nodeConsumers.AddConsumer(currentNode.ID, consumerID)
				// 	}
				// }
			}
		}

		// prevNode = currentNode
		currentNode = childNode
	}
}

func (t *Tree) RemoveSubscriptionStatic(consumerID types.ConsumerID, topic types.Topic) {
	t.staticConsumersMu.Lock()
	defer t.staticConsumersMu.Unlock()

	inner, ok := t.staticConsumers[topic]
	if !ok {
		return
	}
	_, ok = inner[consumerID]
	if !ok {
		return
	}
	delete(inner, consumerID)
	if len(inner) == 0 {
		delete(t.staticConsumers, topic)
	}
}

func (t *Tree) RemoveSubscription(consumerID types.ConsumerID, parts []string) {
	if t.root == nil {
		return
	}

	var (
		lastPart    bool
		currentNode *Node
		prevNode    *Node
	)
	currentNode = t.root

	removeFromNode := func(node *Node) {
		left := t.nodeConsumers.RemoveConsumer(node.ID, consumerID)
		if left == 0 {
			node.SetStop(false)
		}
	}

	for i, part := range parts {
		childNode := currentNode.Child(part)
		if childNode == nil {
			// todo: break or continue?
			break
		}
		prevNode = currentNode
		currentNode = childNode

		lastPart = i == (len(parts) - 1)
		if lastPart {
			if !currentNode.IsStop() {
				continue
			}
			// todo: тут плохо то, что флаг stop-node еще не значит, что имеено этот консьюмер на нее подписан
			removeFromNode(currentNode)
			if currentNode.IsHash() {
				// todo: can be nil and panic?
				if prevNode.ID > 0 {
					removeFromNode(prevNode)
				}
			}
		}
	}
}
