package topics

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestParser_ParseTopic(t *testing.T) {
	t.Run("table test", func(t *testing.T) {
		tests := []struct {
			topic string
			parts []string
			err   error
		}{
			{
				topic: "",
				err:   ErrTopicEmpty,
			},
			{
				topic: strings.Repeat("a", 65),
				err:   ErrTopicTooLong,
			},
			{
				topic: "one@two",
				err:   ErrTopicInvalidChar,
			},
			{
				topic: "@one",
				err:   ErrTopicInvalidChar,
			},
			{
				topic: "one@",
				err:   ErrTopicInvalidChar,
			},
			{
				topic: ".one",
				err:   ErrTopicEmptyPart,
			},
			{
				topic: "one..two",
				err:   ErrTopicEmptyPart,
			},
			{
				topic: "one.",
				err:   ErrTopicEmptyPart,
			},
			{
				topic: ".",
				err:   ErrTopicEmptyPart,
			},
			{
				topic: "one.two",
				err:   nil,
				parts: []string{"one", "two"},
			},
		}

		parser := NewParser()

		for _, test := range tests {
			parts, err := parser.ParseTopic(test.topic)

			require.Equal(t, test.err, err, "topic=%s", test.topic)
			assert.Equal(t, test.parts, parts, "topic=%s", test.topic)
		}
	})
}
