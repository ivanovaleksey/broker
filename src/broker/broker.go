package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"github.com/ivanovaleksey/broker/src/topics"
	"github.com/ivanovaleksey/broker/src/tree"
	"go.uber.org/zap"
)

type Broker struct {
	logger *zap.Logger

	bytesTopic  *topics.BytesParser
	topicParser TopicParser
	tree        *tree.Tree
}

type TopicParser interface {
	IsStatic(topic string) (bool, error)
	ParseTopic(string) ([]string, error)
	ParsePattern(string) ([]string, error)
}

func NewBroker(l *zap.Logger) *Broker {
	return &Broker{
		logger:      l,
		bytesTopic:  topics.NewBytesParser(),
		topicParser: topics.NewParser(),
		tree:        tree.NewTree(),
	}
}

func (b *Broker) GetConsumers(topic string) ([]types.ConsumerID, error) {
	// parts, err := b.topicParser.ParseTopic(topic)
	// if err != nil {
	// 	return nil, err
	// }
	if len(topic) == 0 {
		return nil, nil
	}

	topicHash := b.bytesTopic.Hash(topic)
	parts := b.bytesTopic.Parts(topic)
	consumerIDs := b.tree.GetConsumers(topicHash, parts)
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
