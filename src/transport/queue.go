package transport

import (
	"context"
	pb "github.com/ivanovaleksey/broker/pkg/pb/broker_fast"
	"go.uber.org/zap"
	"sync"
)

const (
	defaultQueueWorkers = 1024
	defaultQueueSize    = 1024
)

type Queue struct {
	logger *zap.Logger

	queue      chan Message
	numWorkers int

	wg     sync.WaitGroup
	ctx    context.Context
	cancel func()
}

type Message struct {
	To   Consumer
	Key  string
	Data []byte
}

func NewQueue(ctx context.Context, logger *zap.Logger) *Queue {
	ctx, cancel := context.WithCancel(ctx)

	q := &Queue{
		ctx:        ctx,
		cancel:     cancel,
		logger:     logger,
		queue:      make(chan Message),
		numWorkers: defaultQueueWorkers,
	}
	return q
}

func (q *Queue) Push(msg Message) {
	// todo: push with select and timeout?
	q.queue <- msg
}

func (q *Queue) RunBackground() {
	for i := 0; i < q.numWorkers; i++ {
		q.wg.Add(1)
		go func() {
			defer q.wg.Done()

			for {
				select {
				case <-q.ctx.Done():
					return
				case env := <-q.queue:
					// q.logger.Debug("got envelope", zap.String("topic", env.Message.Key))
					q.handleMessage(env)
				}
			}
		}()
	}
}

func (q *Queue) handleMessage(env Message) {
	// todo: how to handle err?
	err := env.To.Send(&pb.ConsumeResponse{
		Key:     env.Key,
		Payload: env.Data,
	})
	if err != nil {
		q.logger.Error("can't send", zap.Error(err))
	}
}

// Close waits for workers to stop
func (q *Queue) Close() error {
	select {
	case <-q.queue:
		return nil
	default:
	}

	q.cancel()
	q.wg.Wait()
	close(q.queue)
	return nil
}
