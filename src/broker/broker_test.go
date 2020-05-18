package broker_test

import (
	"github.com/ivanovaleksey/broker/src/broker"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap"
	"math/rand"
	"testing"
)

func TestBroker_GetConsumers(t *testing.T) {
	t.Run("spec table test", func(t *testing.T) {
		tests := []struct {
			subscription  string
			expectedTrue  []string
			expectedFalse []string
		}{
			{
				// сообщения с любым ключом
				subscription:  "#",
				expectedTrue:  []string{"logs", "user", "messages"},
				expectedFalse: nil,
			},
			{
				// сообщения с ключом содержащим одно слово, например, logs, user, messages
				subscription:  "*",
				expectedTrue:  []string{"logs", "user", "messages"},
				expectedFalse: []string{"logs.user"},
			},
			{
				subscription:  "one.*",
				expectedTrue:  []string{"one.two", "one.three", "one.four"},
				expectedFalse: []string{"one.two.three"},
			},
			{
				subscription:  "*.one",
				expectedTrue:  []string{"two.one", "three.one", "four.one"},
				expectedFalse: []string{"three.two.one"},
			},
			{
				subscription:  "#.one",
				expectedTrue:  []string{"one", "two.one", "three.two.one"},
				expectedFalse: []string{"one.two"},
			},
			{
				subscription:  "one.*.two",
				expectedTrue:  []string{"one.three.two", "one.four.two", "one.five.two"},
				expectedFalse: []string{"one.two"},
			},
			{
				subscription:  "one.#.two",
				expectedTrue:  []string{"one.two", "one.three.two", "one.three.four.two"},
				expectedFalse: []string{"one.one"},
			},
			{
				subscription:  "one.*.*.two",
				expectedTrue:  []string{"one.three.four.two", "one.five.six.two"},
				expectedFalse: []string{"one.two", "one.five.six.seven.two"},
			},
			{
				subscription:  "one.*.two.*",
				expectedTrue:  []string{"one.three.two.four", "one.four.two.five"},
				expectedFalse: []string{"one.two", "one.three.two", "one.two.four", "one.three.more.two.four"},
			},
			{
				subscription:  "one.#.two.*",
				expectedTrue:  []string{"one.two.three", "one.three.two.four", "one.three.four.two.five"},
				expectedFalse: []string{"one.two", "one.three.two"},
			},
			{
				subscription:  "one.*.two.#",
				expectedTrue:  []string{"one.three.two", "one.three.two.four", "one.three.two.four.five"},
				expectedFalse: []string{"one.two"},
			},
		}

		consumerID := rand.Int63()
		for _, test := range tests {
			t.Run(test.subscription, func(t *testing.T) {
				brk := broker.NewBroker(zap.NewNop())
				brk.Subscribe(consumerID, []string{test.subscription})

				for _, topic := range test.expectedTrue {
					consumers, err := brk.GetConsumers(topic)
					require.NoError(t, err, "topic=%s", topic)
					assert.Equal(t, []int64{consumerID}, consumers, "topic=%s", topic)
				}

				for _, topic := range test.expectedFalse {
					consumers, err := brk.GetConsumers(topic)
					require.NoError(t, err, "topic=%s", topic)
					assert.Empty(t, consumers, "topic=%s", topic)
				}
			})
		}
	})
}
