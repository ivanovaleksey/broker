package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"go.uber.org/zap"
)

func (b *Broker) Unsubscribe(id types.ConsumerID, topics []types.Topic) {
	for _, topic := range topics {
		if err := b.handleUnsubscribe(id, topic); err != nil {
			b.logger.Error("can't unsubscribe", zap.Int64("consumer_id", id), zap.String("topic", topic), zap.Error(err))
		}
	}
}

func (b *Broker) handleUnsubscribe(id types.ConsumerID, pattern types.Topic) error {
	parts, err := b.topicParser.ParsePattern(pattern)
	if err != nil {
		return err
	}

	// parts = prepareParts(parts)
	b.tree.RemoveSubscription(id, parts)
	return nil
}
