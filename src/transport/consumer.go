package transport

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"github.com/ivanovaleksey/broker/pkg/types"
	"go.uber.org/zap"
	"io"
	"math/rand"
)

func (t *Transport) Consume(stream pb.MessageBroker_ConsumeServer) error {
	ctx := stream.Context()
	logger := ctxzap.Extract(ctx)

	consumerID := rand.Int63()
	logger.Debug("start consume", zap.Int64("consumer_id", consumerID))

	for {
		select {
		case <-ctx.Done():
			logger.Debug("consume ctx done")
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		switch {
		case err == io.EOF:
			logger.Debug("--EOF--")
			return nil
		case err != nil:
			logger.Error("recv error", zap.Error(err))
			return err
		}

		logger.Debug("consumed message", zap.String("action", req.Action.String()), zap.Strings("keys", req.Keys), zap.Int64("consumer", consumerID))

		switch req.Action {
		case pb.ConsumeRequest_SUBSCRIBE:
			t.AddConsumer(consumerID, stream)
			t.broker.Subscribe(consumerID, req.Keys)
		case pb.ConsumeRequest_UNSUBSCRIBE:
			// todo: maybe remove from consumers if no subscriptions?
			t.broker.Unsubscribe(consumerID, req.Keys)
		}
	}
}

func (t *Transport) AddConsumer(id types.ConsumerID, consumer Consumer) {
	index := id % bucketsCount
	t.consumersLocks[index].Lock()
	t.consumers[index][id] = consumer
	t.consumersLocks[index].Unlock()
}

func (t *Transport) GetConsumer(id types.ConsumerID) (Consumer, bool) {
	index := id % bucketsCount
	t.consumersLocks[index].RLock()
	c, ok := t.consumers[index][id]
	t.consumersLocks[index].RUnlock()
	return c, ok
}
