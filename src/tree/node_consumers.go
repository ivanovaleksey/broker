package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"sync"
)

const bucketsCount = 16

type NodeConsumers = map[types.NodeID]Consumers
type Consumers = map[types.ConsumerID]int

// ConsumersLog links nodes and registered consumers
type ConsumersLog struct {
	locks         [bucketsCount]sync.RWMutex
	nodeConsumers [bucketsCount]NodeConsumers
}

func NewConsumersLog() *ConsumersLog {
	log := &ConsumersLog{}
	for i := 0; i < bucketsCount; i++ {
		log.nodeConsumers[i] = make(NodeConsumers, 1)
	}
	return log
}

func (l *ConsumersLog) AddConsumer(nodeID types.NodeID, consumerID types.ConsumerID) {
	bucketIndex := nodeID % bucketsCount
	l.locks[bucketIndex].Lock()
	hash, ok := l.nodeConsumers[bucketIndex][nodeID]
	if !ok {
		hash = make(Consumers, 1)
		l.nodeConsumers[bucketIndex][nodeID] = hash
	}
	hash[consumerID]++
	l.locks[bucketIndex].Unlock()
}

// RemoveConsumer returns how much consumers are still subscribed to this node
func (l *ConsumersLog) RemoveConsumer(nodeID types.NodeID, consumerID types.ConsumerID) int {
	bucketIndex := nodeID % bucketsCount
	l.locks[bucketIndex].Lock()
	defer l.locks[bucketIndex].Unlock()
	hash, _ := l.nodeConsumers[bucketIndex][nodeID]
	// if !ok {
	// 	todo: need it? unsubscribe before subscribe?
	// 	return 0
	// }
	if _, ok := hash[consumerID]; !ok {
		// this consumer is not subscribed now
		return len(hash)
	}
	hash[consumerID]--
	if hash[consumerID] == 0 {
		// consumer has no more subscriptions on this node
		delete(hash, consumerID)
	}
	l.nodeConsumers[bucketIndex][nodeID] = hash
	return len(hash)
}

func (l *ConsumersLog) GetConsumers(nodeID types.NodeID) Consumers {
	bucketIndex := nodeID % bucketsCount
	l.locks[bucketIndex].RLock()
	defer l.locks[bucketIndex].RUnlock()
	return l.nodeConsumers[bucketIndex][nodeID]
}
