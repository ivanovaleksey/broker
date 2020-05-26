package transport

import (
	"context"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"github.com/ivanovaleksey/broker/pkg/types"
	"go.uber.org/zap"
	"sync"
)

const bucketsCount = 256

// type Consumer = pb.MessageBroker_ConsumeServer
type Consumer interface {
	Send(*pb.ConsumeResponse) error
}

type Transport struct {
	broker Broker
	queue  *Queue

	// mu        sync.RWMutex
	// consumers map[types.ConsumerID]Consumer
	consumersLocks [bucketsCount]sync.RWMutex
	consumers      [bucketsCount]map[types.ConsumerID]Consumer
}

func NewTransport(ctx context.Context, logger *zap.Logger, b Broker) *Transport {
	t := &Transport{
		broker: b,
		queue:  NewQueue(ctx, logger),
		// consumers: make(map[types.ConsumerID]Consumer),
	}
	for i := 0; i < bucketsCount; i++ {
		t.consumers[i] = make(map[types.ConsumerID]Consumer, 1)
	}
	return t
}

func (t *Transport) Start() {
	t.queue.RunBackground()
}

type Broker interface {
	Subscribe(id types.ConsumerID, topics []types.Topic)
	// todo: consider returning list of remaining subscriptions
	Unsubscribe(id types.ConsumerID, topics []types.Topic)
	GetConsumers(topic string) ([]types.ConsumerID, error)
}

func (t *Transport) Close() error {
	return t.queue.Close()
}
