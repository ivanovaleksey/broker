package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"go.uber.org/zap"
)

func (b *Broker) Unsubscribe(id types.ConsumerID, topics []string) {
	for _, topic := range topics {
		if err := b.handleUnsubscribe(id, topic); err != nil {
			b.logger.Error("can't unsubscribe", zap.Int64("consumer_id", id), zap.Error(err))
		}
	}
}

func (b *Broker) handleUnsubscribe(id types.ConsumerID, pattern string) error {
	isStatic, err := b.topicParser.IsStatic(pattern)
	if err != nil {
		return err
	}
	if isStatic {
		b.tree.RemoveSubscriptionStatic(id, types.Topic(pattern))
		return nil
	}

	parts, err := b.topicParser.ParsePattern(pattern)
	if err != nil {
		return err
	}

	// parts = prepareParts(parts)
	b.tree.RemoveSubscription(id, parts)
	return nil
}
