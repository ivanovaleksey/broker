package main

import (
	"context"
	"flag"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"io"
	"log"
	"time"
)

func main() {
	num := flag.Int("num", 1, "Client number")
	flag.Parse()

	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatal(err)
	}

	if err := run(logger, *num); err != nil {
		logger.Fatal("run error", zap.Error(err))
	}
}

func run(l *zap.Logger, num int) error {
	addr := "127.0.0.1:3000"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "can't dial")
	}

	client := pb.NewMessageBrokerClient(conn)

	// ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	// defer cancel()

	consumer, err := client.Consume(context.Background())
	if err != nil {
		return errors.Wrap(err, "can't consume")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	go func() {
		subscribe := func(topic string) error {
			l.Debug("subscribing", zap.String("topic", topic))
			req := &pb.ConsumeRequest{
				Action: pb.ConsumeRequest_SUBSCRIBE,
				Keys:   []string{topic},
			}
			return consumer.Send(req)
		}

		subscribe("topic_1")
		time.Sleep(time.Second * 5)
		subscribe("topic_2")

		consumer.CloseSend()
	}()

	l.Debug("ready to recv")
	for {
		resp, err := consumer.Recv()
		switch {
		case err == io.EOF:
			l.Debug("--EOF--")
			return nil
		case err != nil:
			return errors.Wrap(err, "can't recv")
		}

		l.Debug("got message", zap.String("topic", resp.Key), zap.ByteString("data", resp.Payload))
	}
}
