package broker

import (
	"github.com/ivanovaleksey/broker/pkg/types"
	"github.com/ivanovaleksey/broker/src/node"
	"go.uber.org/zap"
)

func (b *Broker) Subscribe(id types.ConsumerID, topics []string) {
	for _, topic := range topics {
		if err := b.handleSubscribe(id, topic); err != nil {
			b.logger.Error("can't subscribe", zap.Int64("consumer_id", id), zap.Error(err))
		}
	}
}

func (b *Broker) handleSubscribe(id types.ConsumerID, pattern string) error {
	if len(pattern) == 0 {
		return nil
	}
	info := b.bytesTopic.Parse(pattern)
	if info.IsStatic {
		b.tree.AddSubscriptionStatic(id, info.Part)
		return nil
	}

	b.tree.AddSubscription(id, info.Parts)
	return nil

	// isStatic, err := b.topicParser.IsStatic(pattern)
	// if err != nil {
	// 	return err
	// }
	// if isStatic {
	// 	b.tree.AddSubscriptionStatic(id, types.Topic(pattern))
	// 	return nil
	// }
	//
	// parts, err := b.topicParser.ParsePattern(pattern)
	// if err != nil {
	// 	return err
	// }
	//
	// // parts = prepareParts(parts)
	// b.tree.AddSubscription(id, parts)
	// return nil
}

// todo: really need to combine sequent hashes?
// it mutates input parts
func prepareParts(parts []string) []string {
	out := parts[:0]
	for i := 0; i < len(parts); i++ {
		part := parts[i]
		if part == node.NodeHash {
			nextIdx := i + 1
			if nextIdx < len(parts) && parts[nextIdx] == node.NodeHash {
				// skip sequent #
				continue
			}
		}
		out = append(out, part)
	}
	return out
}
