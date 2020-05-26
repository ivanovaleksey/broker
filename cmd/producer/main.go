package main

import (
	"context"
	"flag"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"github.com/pkg/errors"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"log"
	"strconv"
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
	addr := "127.0.0.1:8000"
	conn, err := grpc.Dial(addr, grpc.WithInsecure())
	if err != nil {
		return errors.Wrap(err, "can't dial")
	}

	client := pb.NewMessageBrokerClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
	defer cancel()

	producer, err := client.Produce(context.Background())
	if err != nil {
		return errors.Wrap(err, "can't produce")
	}

	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	var (
		key  = "topic" + strconv.Itoa(num)
		data []byte
	)

loop:
	for {
		select {
		case <-ctx.Done():
			l.Debug("ctx done")
			break loop
		case <-ticker.C:
			l.Debug("tick")
			data = []byte("data_" + strconv.Itoa(num) + "_" + time.Now().Format(time.RFC3339))
			req := &pb.ProduceRequest{
				Key:     key,
				Payload: data,
			}
			if err := producer.Send(req); err != nil {
				// l.Error("can't send", zap.Error(err))
				return errors.Wrap(err, "can't send")
			}
		}
	}

	l.Debug("closing")
	if _, err := producer.CloseAndRecv(); err != nil {
		return errors.Wrap(err, "can't close and recv")
	}
	l.Debug("closed")
	return nil
}
