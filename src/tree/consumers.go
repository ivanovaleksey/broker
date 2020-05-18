package tree

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"sync"
)

const bucketsCount = 16

type ConsumersHash = map[types.NodeID][]types.ConsumerID

// ConsumersLog links nodes and registered consumers
type ConsumersLog struct {
	locks         [bucketsCount]sync.RWMutex
	nodeConsumers [bucketsCount]ConsumersHash
}

func NewConsumersLog() *ConsumersLog {
	log := &ConsumersLog{}
	for i := 0; i < bucketsCount; i++ {
		log.nodeConsumers[i] = make(ConsumersHash)
	}
	return log
}

func (l *ConsumersLog) AddConsumer(nodeID types.NodeID, consumerID types.ConsumerID) {
	bucketIndex := nodeID % bucketsCount
	l.locks[bucketIndex].Lock()
	l.nodeConsumers[bucketIndex][nodeID] = append(l.nodeConsumers[bucketIndex][nodeID], consumerID)
	l.locks[bucketIndex].Unlock()
}

func (l *ConsumersLog) GetConsumers(nodeID types.NodeID) []types.ConsumerID {
	bucketIndex := nodeID % bucketsCount
	l.locks[bucketIndex].RLock()
	defer l.locks[bucketIndex].RUnlock()
	return l.nodeConsumers[bucketIndex][nodeID]
}
