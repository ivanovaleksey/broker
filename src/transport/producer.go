package transport

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"go.uber.org/zap"
	"google.golang.org/grpc/status"
	"io"
)

func (t *Transport) Produce(stream pb.MessageBroker_ProduceServer) error {
	ctx := stream.Context()
	logger := ctxzap.Extract(ctx)

	logger.Debug("start reading produce")

	for {
		select {
		case <-ctx.Done():
			logger.Debug("produce ctx done")
			return ctx.Err()
		default:
		}

		req, err := stream.Recv()
		switch {
		case err == io.EOF:
			logger.Debug("--EOF--")
			return stream.SendAndClose(&pb.ProduceResponse{})
		case err != nil:
			logger.Error("producer recv error", zap.String("code", status.Code(err).String()))
			return err
		}

		logger.Debug("produced message", zap.String("key", req.Key), zap.ByteString("payload", req.Payload))
		consumerIDs, err := t.broker.GetConsumers(req.Key)
		if err != nil {
			logger.Error("can't get consumers", zap.String("topic", req.Key))
			continue
		}
		logger.Debug("got consumers", zap.String("topic", req.Key), zap.Int64s("ids", consumerIDs))

		// todo: consider using queue with N workers for sending
		for _, consumerID := range consumerIDs {
			consumer, ok := t.GetConsumer(consumerID)
			if !ok {
				logger.Debug("no consumer", zap.Int64("id", consumerID))
				continue
			}

			logger.Debug("sending to consumer", zap.Int64("id", consumerID))
			go func(target Consumer) {
				// todo: how to handle err?
				consumer.Send(&pb.ConsumeResponse{
					Key:     req.Key,
					Payload: req.Payload,
				})
			}(consumer)
		}
	}
}
