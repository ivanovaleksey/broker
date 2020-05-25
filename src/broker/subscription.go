package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"github.com/ivanovaleksey/broker/src/tree"
	"go.uber.org/zap"
)

func (b *Broker) Subscribe(id types.ConsumerID, topics []types.Topic) {
	for _, topic := range topics {
		if err := b.handleSubscribe(id, topic); err != nil {
			b.logger.Error("can't subscribe", zap.Int64("consumer_id", id), zap.String("topic", topic), zap.Error(err))
		}
	}
}

func (b *Broker) handleSubscribe(id types.ConsumerID, pattern types.Topic) error {
	parts, err := b.topicParser.ParsePattern(pattern)
	if err != nil {
		return err
	}

	// parts = prepareParts(parts)
	b.tree.AddSubscription(id, parts)
	return nil
}

// todo: really need to combine sequent hashes?
// it mutates input parts
func prepareParts(parts []string) []string {
	out := parts[:0]
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == tree.NodeHash {
			nextIdx := i+1
			if nextIdx < len(parts) && parts[nextIdx] == tree.NodeHash {
				// skip sequent #
				continue
			}
		}
		out = append(out, part)
	}
	return out
}
