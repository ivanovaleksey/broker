package tree

import (
	"github.com/ivanovaleksey/broker/pkg/sort"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestTree_GetConsumers(t *testing.T) {
	t.Run("test stop nodes logic 1", func(t *testing.T) {
		var (
			consumerA int64 = 100
			consumerB int64 = 200
			consumerC int64 = 300
			consumerD int64 = 400
			consumerE int64 = 500
		)

		tree := NewTree()
		tree.AddSubscription(consumerA, []string{"#", "debug"})
		tree.AddSubscription(consumerB, []string{"#"})
		tree.AddSubscription(consumerC, []string{"events", "orders", "*"})
		tree.AddSubscription(consumerD, []string{"events"})
		tree.AddSubscription(consumerE, []string{"events", "orders", "#"})

		tests := []struct {
			parts    []string
			expected []int64
		}{
			{
				parts:    []string{"any", "debug"},
				expected: []int64{consumerA, consumerB},
			},
			{
				parts:    []string{"debug"},
				expected: []int64{consumerA, consumerB},
			},
			{
				parts:    []string{"any"},
				expected: []int64{consumerB},
			},
			{
				parts:    []string{"events", "orders", "paid"},
				expected: []int64{consumerB, consumerC, consumerE},
			},
			{
				parts:    []string{"events", "orders", "paid", "some"},
				expected: []int64{consumerB, consumerE},
			},
			{
				parts:    []string{"events", "orders"},
				expected: []int64{consumerB, consumerE},
			},
			{
				parts:    []string{"events", "payments"},
				expected: []int64{consumerB},
			},
			{
				parts:    []string{"events"},
				expected: []int64{consumerB, consumerD},
			},
		}

		for _, test := range tests {
			t.Run(strings.Join(test.parts, "."), func(t *testing.T) {
				consumers := tree.GetConsumers(test.parts)
				exp := sort.Int64SLice(test.expected)
				exp.Sort()
				act := sort.Int64SLice(consumers)
				act.Sort()
				assert.Equal(t, exp, act)
			})
		}
	})

	t.Run("test termination", func(t *testing.T) {
		const (
			consumerA int64 = (iota + 1) * 100
		)

		tree := NewTree()
		tree.AddSubscription(consumerA, []string{"1", "#", "2"})

		tests := []struct {
			parts    []string
			expected []int64
		}{
			{
				parts:    []string{"1", "3", "2", "4", "2"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"1", "3", "2"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"1", "3", "2", "4"},
				expected: []int64{},
			},
		}

		for _, test := range tests {
			t.Run(strings.Join(test.parts, "."), func(t *testing.T) {
				consumers := tree.GetConsumers(test.parts)
				exp := sort.Int64SLice(test.expected)
				exp.Sort()
				act := sort.Int64SLice(consumers)
				act.Sort()
				assert.Equal(t, exp, act)
			})
		}
	})

	t.Run("hash case", func(t *testing.T) {
		var (
			consumerA int64 = 100
			// consumerB int64 = 200
			consumerC int64 = 300
		)

		tree := NewTree()
		// tree.AddSubscription(consumerA, []string{"#", "#", "1", "#", "2", "#"})
		tree.AddSubscription(consumerA, []string{"#", "1", "#", "2", "#"})
		// tree.AddSubscription(consumerB, []string{"#"})
		tree.AddSubscription(consumerC, []string{"1", "#", "2"})

		tests := []struct {
			parts    []string
			expected []int64
		}{
			{
				parts:    []string{"1", "2"},
				expected: []int64{consumerA, consumerC},
			},
			{
				parts:    []string{"1", "x", "2"},
				expected: []int64{consumerA, consumerC},
			},
			{
				parts:    []string{"x", "1", "2"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"1", "x", "y"},
				expected: []int64{},
			},
			{
				parts:    []string{"1", "2", "2"},
				expected: []int64{consumerA, consumerC},
			},
			{
				parts:    []string{"1", "2", "x"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"x", "y", "2"},
				expected: []int64{},
			},
			{
				parts:    []string{"x", "y", "z", "1"},
				expected: []int64{},
			},
			{
				parts:    []string{"x", "y", "z", "1", "2"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"x", "y", "z", "1", "2", "1"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"x", "y", "z", "1", "2", "1", "x"},
				expected: []int64{consumerA},
			},
			{
				parts:    []string{"1", "x", "2", "y", "2"},
				expected: []int64{consumerA, consumerC},
			},
		}

		for _, test := range tests {
			t.Run(strings.Join(test.parts, "."), func(t *testing.T) {
				consumers := tree.GetConsumers(test.parts)
				exp := sort.Int64SLice(test.expected)
				exp.Sort()
				act := sort.Int64SLice(consumers)
				act.Sort()
				assert.Equal(t, exp, act)
			})
		}
	})
}
