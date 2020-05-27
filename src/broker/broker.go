package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"github.com/ivanovaleksey/broker/src/topics"
	"github.com/ivanovaleksey/broker/src/tree"
	"go.uber.org/zap"
)

type Broker struct {
	logger      *zap.Logger
	topicParser TopicParser
	tree        *tree.Tree
}

type TopicParser interface {
	IsStatic(topic types.Topic) (bool, error)
	ParseTopic(types.Topic) ([]string, error)
	ParsePattern(types.Topic) ([]string, error)
}

func NewBroker(l *zap.Logger) *Broker {
	return &Broker{
		logger:      l,
		topicParser: topics.NewParser(),
		tree:        tree.NewTree(),
	}
}

func (b *Broker) GetConsumers(topic string) ([]types.ConsumerID, error) {
	parts, err := b.topicParser.ParseTopic(topic)
	if err != nil {
		return nil, err
	}

	consumerIDs := b.tree.GetConsumers(parts)
	return consumerIDs, nil
}

// func (b *Broker) SendMessage(topic string, data []byte) {
//
// 	// fetch consumer by topic
// 	topic := req.Key
// 	if len(t.consumers) == 0 {
// 		logger.Debug("no consumers", zap.String("topic", topic))
// 		continue
// 	}
//
// 	consumer := t.consumers[0]
// 	// todo: не дублировать при нескольких пересекающихся подписках?
// 	resp := &pb.ConsumeResponse{
// 		Key:     topic,
// 		Payload: req.Payload,
// 	}
// 	if err := consumer.Send(resp); err != nil {
// 		// todo: return this err (what would be with `stream`)?
// 		logger.Error("can't send to consumer", zap.Error(err))
// 	}
// }
//
