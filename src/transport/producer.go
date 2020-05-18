package transport

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/logging/zap/ctxzap"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"go.uber.org/zap"
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
			logger.Error("produce receive error", zap.Error(err))
			return err
		}

		// todo: handle produce request
		logger.Debug("produced message", zap.String("key", req.Key), zap.ByteString("payload", req.Payload))

		consumerIDs, err := t.broker.GetConsumers(req.Key)
		if err != nil {
			logger.Error("can't get consumers", zap.String("topic", req.Key))
			continue
		}

		logger.Debug("got consumers", zap.String("topic", req.Key), zap.Int64s("ids", consumerIDs))

		// topic := req.Key
		// consumers := t.broker.Consumers(topic)
		// if len(consumers) == 0 {
		// 	logger.Debug("no consumers", zap.String("topic", topic))
		// 	continue
		// }
		//
		// if err != nil {
		// 	logger.Error("send message", zap.Error(err))
		// 	return err
		// }
	}
}

// func Produce() error {
//
// }
