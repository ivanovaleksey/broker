package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"github.com/ivanovaleksey/broker/src/node"
	"github.com/ivanovaleksey/broker/src/topics"
	"sync"
	"sync/atomic"
)

const bucketsCountConsumers = 16

type Tree struct {
	root *node.Node

	nodeConsumersActive int32
	nodeConsumers       *ConsumersLog

	staticConsumersLocks [bucketsCountConsumers]sync.RWMutex
	staticConsumers      [bucketsCountConsumers]map[uint64][]types.ConsumerID
}

func NewTree() *Tree {
	const (
		rootSize = 1 << 18
		hashSize = 1 << 15
		starSize = 1 << 15
	)

	root := node.NewNode()
	root.ID = -1
	root.Next = make(map[uint64]*node.Node, rootSize)

	star := node.NewNode()
	star.ID = -2
	star.Type = node.NodeTypeStar
	star.Next = make(map[uint64]*node.Node, starSize)
	root.SetChild(star, topics.HashStar)

	hash := node.NewNode()
	hash.ID = -3
	hash.Type = node.NodeTypeHash
	hash.Next = make(map[uint64]*node.Node, hashSize)
	root.SetChild(hash, topics.HashHash)

	starStar := node.NewNode()
	starStar.ID = -4
	starStar.Type = node.NodeTypeStar
	starStar.Next = make(map[uint64]*node.Node, starSize)
	star.SetChild(starStar, topics.HashStar)

	log := NewConsumersLog()

	t := &Tree{
		root:          root,
		nodeConsumers: log,
	}

	const staticConsumersSize = 200000
	// staticConsumers := make(map[uint64]map[types.ConsumerID]struct{}, staticConsumersSize)
	for i := 0; i < bucketsCountConsumers; i++ {
		m := make(map[uint64][]types.ConsumerID, staticConsumersSize)
		t.staticConsumers[i] = m
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

func (t *Tree) AddSubscriptionStatic(consumerID types.ConsumerID, topicHash uint64) {
	// topicHash := topic.Hash()
	idx := topicHash % bucketsCountConsumers

	t.staticConsumersLocks[idx].Lock()
	defer t.staticConsumersLocks[idx].Unlock()

	var alreadySubscribed bool
	inner, ok := t.staticConsumers[idx][topicHash]
	if !ok {
		// todo: can get slice from pool?
		// inner = make([]types.ConsumerID, 0, 12)
		inner = GetFromPool()
	} else {
		for i := 0; i < len(inner); i++ {
			if inner[i] == consumerID {
				alreadySubscribed = true
				break
			}
		}
	}
	if !alreadySubscribed {
		inner = append(inner, consumerID)
		t.staticConsumers[idx][topicHash] = inner
	}
}

// AddSubscription receives already prepared parts
func (t *Tree) AddSubscription(consumerID types.ConsumerID, parts []uint64) {
	if t.root == nil {
		return
	}
	atomic.CompareAndSwapInt32(&t.nodeConsumersActive, 0, 1)

	var (
		lastPart    bool
		currentNode *node.Node
		// prevNode    *Node
	)
	currentNode = t.root

	fn := func(childNode, parentNode *node.Node) {
		t.nodeConsumers.AddConsumer(childNode.ID, consumerID)
		// todo: is it ok or should be done in smarter way?
		if childNode.IsHash() {
			// todo: check this is not root node, maybe better check
			if !parentNode.IsRoot() {
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
			// newNode := NewNode()
			newNode := GetNodeFromPool()
			newNode.SetType(part)

			if lastPart {
				newNode.Stop = 1
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

func (t *Tree) RemoveSubscriptionStatic(consumerID types.ConsumerID, topicHash uint64) {
	// topicHash := topic.Hash()
	idx := topicHash % bucketsCountConsumers

	t.staticConsumersLocks[idx].Lock()
	defer t.staticConsumersLocks[idx].Unlock()

	inner, ok := t.staticConsumers[idx][topicHash]
	if !ok {
		return
	}

	var (
		found bool
	)
	for i := 0; i < len(inner); i++ {
		if inner[i] == consumerID {
			found = true
			inner[i] = inner[len(inner)-1]
			inner = inner[:len(inner)-1]
			t.staticConsumers[idx][topicHash] = inner
			break
		}
	}
	if !found {
		return
	}
	if len(inner) == 0 {
		PutToPool(inner)
		delete(t.staticConsumers[idx], topicHash)
	}
}

func (t *Tree) RemoveSubscription(consumerID types.ConsumerID, parts []uint64) {
	if t.root == nil {
		return
	}

	var (
		lastPart    bool
		currentNode *node.Node
		prevNode    *node.Node
	)
	currentNode = t.root

	removeConsumerFromNode := func(node *node.Node) int {
		left := t.nodeConsumers.RemoveConsumer(node.ID, consumerID)
		if left == 0 {
			// todo: remove node? put to pool?
			node.SetStop(false)
		}
		return left
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
			left := removeConsumerFromNode(currentNode)
			if currentNode.IsHash() {
				// todo: can be nil and panic?
				if !prevNode.IsRoot() {
					// todo: how to delete node here?
					removeConsumerFromNode(prevNode)
				}
			}
			if left == 0 && !prevNode.IsRoot() {
				prevNode.RemoveChild(part)
				PutNodeToPool(currentNode)
			}
		}
	}
}
