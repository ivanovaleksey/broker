package tree

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTree_AddSubscription(t *testing.T) {
	const (
		consumerA int64 = (iota + 1) * 100
		consumerB
		consumerC
		consumerD
		consumerE
	)

	t.Run("initial tree idea", func(t *testing.T) {
		t.Skip("should not pass with plain tree logic")

		t.Run("hash case", func(t *testing.T) {
			tree := NewTree()
			tree.AddSubscription(consumerA, []string{"#", "1", "#", "2", "#"})
			tree.AddSubscription(consumerB, []string{"#"})
			tree.AddSubscription(consumerC, []string{"1", "#", "2"})

			require.Equal(t, 2, len(tree.root.Next))

			firstHashNode := tree.root.Child("#")
			require.NotNil(t, firstHashNode)
			require.Equal(t, 1, len(firstHashNode.Next))
			assert.Equal(t, NodeTypeHash, firstHashNode.Type)
			assert.Equal(t, true, firstHashNode.Stop)
			firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(firstHashNode.ID)
			assert.Equal(t, []int64{consumerB}, firstHashNodeConsumers)

			oneNode := tree.root.Child("1")
			require.NotNil(t, oneNode)
			require.Equal(t, 2, len(oneNode.Next))
			assert.Equal(t, NodeTypeWord, oneNode.Type)
			assert.Equal(t, false, oneNode.Stop)
			oneNodeConsumers := tree.nodeConsumers.GetConsumers(oneNode.ID)
			assert.Empty(t, oneNodeConsumers)

			secondHashNode := oneNode.Child("#")
			require.NotNil(t, secondHashNode)
			require.Equal(t, 1, len(secondHashNode.Next))
			assert.Equal(t, NodeTypeHash, secondHashNode.Type)
			assert.Equal(t, false, secondHashNode.Stop)
			secondHashNodeConsumers := tree.nodeConsumers.GetConsumers(secondHashNode.ID)
			assert.Empty(t, secondHashNodeConsumers)

			twoNode := oneNode.Child("2")
			require.NotNil(t, twoNode)
			require.Equal(t, 1, len(twoNode.Next))
			assert.Equal(t, NodeTypeWord, twoNode.Type)
			assert.Equal(t, true, twoNode.Stop)
			twoNodeConsumers := tree.nodeConsumers.GetConsumers(twoNode.ID)
			assert.Equal(t, []int64{consumerA, consumerC}, twoNodeConsumers)

			twoNodeFromSecondHash := secondHashNode.Child("2")
			require.NotNil(t, twoNodeFromSecondHash)
			require.Equal(t, twoNode, twoNodeFromSecondHash)
			// require.Equal(t, 1, len(twoNodeFromSecondHash.Next))
			// assert.Equal(t, NodeTypeWord, twoNode.Type)
			// assert.Equal(t, true, twoNodeFromSecondHash.Stop)
			// twoNodeFromSecondHashConsumers := tree.nodeConsumers.GetConsumers(twoNodeFromSecondHash.ID)
			// assert.Equal(t, []int64{consumerA, consumerC}, twoNodeFromSecondHashConsumers)

			thirdHashNode := twoNode.Child("#")
			require.NotNil(t, thirdHashNode)
			require.Equal(t, 0, len(thirdHashNode.Next))
			assert.Equal(t, NodeTypeHash, thirdHashNode.Type)
			assert.Equal(t, true, thirdHashNode.Stop)
			thirdHashNodeConsumers := tree.nodeConsumers.GetConsumers(thirdHashNode.ID)
			assert.Equal(t, []int64{consumerA}, thirdHashNodeConsumers)
		})

		t.Run("logs and events case", func(t *testing.T) {
			tree := NewTree()
			tree.AddSubscription(consumerA, []string{"#", "debug"})
			tree.AddSubscription(consumerB, []string{"#"})
			tree.AddSubscription(consumerC, []string{"events", "orders", "*"})
			tree.AddSubscription(consumerD, []string{"events"})
			tree.AddSubscription(consumerE, []string{"events", "orders", "#"})

			require.Equal(t, 3, len(tree.root.Next))

			firstHashNode := tree.root.Child("#")
			require.NotNil(t, firstHashNode)
			require.Equal(t, 1, len(firstHashNode.Next))
			assert.Equal(t, true, firstHashNode.IsHash())
			assert.Equal(t, true, firstHashNode.Stop)
			firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(firstHashNode.ID)
			assert.Equal(t, []int64{consumerB}, firstHashNodeConsumers)

			debugNode := tree.root.Child("debug")
			require.NotNil(t, debugNode)
			require.Equal(t, 0, len(debugNode.Next))
			assert.Equal(t, NodeTypeWord, debugNode.Type)
			assert.Equal(t, true, debugNode.Stop)
			debugNodeConsumers := tree.nodeConsumers.GetConsumers(debugNode.ID)
			assert.Equal(t, []int64{consumerA}, debugNodeConsumers)

			eventsNode := tree.root.Child("events")
			require.NotNil(t, eventsNode)
			require.Equal(t, 1, len(eventsNode.Next))
			assert.Equal(t, NodeTypeWord, eventsNode.Type)
			assert.Equal(t, true, eventsNode.Stop)
			eventsNodeConsumers := tree.nodeConsumers.GetConsumers(eventsNode.ID)
			assert.Equal(t, []int64{consumerD}, eventsNodeConsumers)

			ordersNode := eventsNode.Child("orders")
			require.NotNil(t, ordersNode)
			require.Equal(t, 2, len(ordersNode.Next))
			assert.Equal(t, NodeTypeWord, ordersNode.Type)
			assert.Equal(t, true, ordersNode.Stop)
			ordersNodeConsumers := tree.nodeConsumers.GetConsumers(ordersNode.ID)
			assert.Equal(t, []int64{consumerE}, ordersNodeConsumers)

			ordersStarNode := ordersNode.Child("*")
			require.NotNil(t, ordersStarNode)
			require.Equal(t, 0, len(ordersStarNode.Next))
			assert.Equal(t, NodeTypeStar, ordersStarNode.Type)
			assert.Equal(t, true, ordersStarNode.Stop)
			ordersStarNodeConsumers := tree.nodeConsumers.GetConsumers(ordersStarNode.ID)
			assert.Equal(t, []int64{consumerC}, ordersStarNodeConsumers)

			ordersHashNode := ordersNode.Child("#")
			require.NotNil(t, ordersHashNode)
			require.Equal(t, 0, len(ordersHashNode.Next))
			assert.Equal(t, NodeTypeHash, ordersHashNode.Type)
			assert.Equal(t, true, ordersHashNode.Stop)
			ordersHashNodeConsumers := tree.nodeConsumers.GetConsumers(ordersHashNode.ID)
			assert.Equal(t, []int64{consumerE}, ordersHashNodeConsumers)
		})
	})

	t.Run("plain tree", func(t *testing.T) {
		t.Run("hash case", func(t *testing.T) {
			tree := NewTree()
			tree.AddSubscription(consumerA, []string{"#", "1", "#", "2", "#"})
			tree.AddSubscription(consumerB, []string{"#"})
			tree.AddSubscription(consumerC, []string{"1", "#", "2"})

			require.Equal(t, 2, len(tree.root.Next))

			{
				hashNode := tree.root.Child("#")
				require.NotNil(t, hashNode)
				require.Equal(t, 1, len(hashNode.Next))
				assert.Equal(t, NodeTypeHash, hashNode.Type)
				assert.Equal(t, true, hashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(hashNode.ID)
				assert.Equal(t, []int64{consumerB}, firstHashNodeConsumers)

				oneNode := hashNode.Child("1")
				require.NotNil(t, oneNode)
				require.Equal(t, 1, len(oneNode.Next))
				assert.Equal(t, NodeTypeWord, oneNode.Type)
				assert.Equal(t, false, oneNode.Stop)
				oneNodeConsumers := tree.nodeConsumers.GetConsumers(oneNode.ID)
				assert.Empty(t, oneNodeConsumers)

				oneHashNode := oneNode.Child("#")
				require.NotNil(t, oneHashNode)
				require.Equal(t, 1, len(oneHashNode.Next))
				assert.Equal(t, NodeTypeHash, oneHashNode.Type)
				assert.Equal(t, false, oneHashNode.Stop)
				oneHashNodeConsumers := tree.nodeConsumers.GetConsumers(oneHashNode.ID)
				assert.Empty(t, oneHashNodeConsumers)

				twoNode := oneHashNode.Child("2")
				require.NotNil(t, twoNode)
				require.Equal(t, 1, len(twoNode.Next))
				assert.Equal(t, NodeTypeWord, twoNode.Type)
				assert.Equal(t, true, twoNode.Stop)
				twoNodeConsumers := tree.nodeConsumers.GetConsumers(twoNode.ID)
				assert.Equal(t, []int64{consumerA}, twoNodeConsumers)

				twoHashNode := twoNode.Child("#")
				require.NotNil(t, twoHashNode)
				require.Equal(t, 0, len(twoHashNode.Next))
				assert.Equal(t, NodeTypeHash, twoHashNode.Type)
				assert.Equal(t, true, twoHashNode.Stop)
				twoHashNodeConsumers := tree.nodeConsumers.GetConsumers(twoHashNode.ID)
				assert.Equal(t, []int64{consumerA}, twoHashNodeConsumers)
			}

			{
				oneNode := tree.root.Child("1")
				require.NotNil(t, oneNode)
				require.Equal(t, 1, len(oneNode.Next))
				assert.Equal(t, NodeTypeWord, oneNode.Type)
				assert.Equal(t, false, oneNode.Stop)
				oneNodeConsumers := tree.nodeConsumers.GetConsumers(oneNode.ID)
				assert.Empty(t, oneNodeConsumers)

				hashNode := oneNode.Child("#")
				require.NotNil(t, hashNode)
				require.Equal(t, 1, len(hashNode.Next))
				assert.Equal(t, NodeTypeHash, hashNode.Type)
				assert.Equal(t, false, hashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(hashNode.ID)
				assert.Empty(t, firstHashNodeConsumers)

				twoNode := hashNode.Child("2")
				require.NotNil(t, twoNode)
				require.Equal(t, 0, len(twoNode.Next))
				assert.Equal(t, NodeTypeWord, twoNode.Type)
				assert.Equal(t, true, twoNode.Stop)
				twoNodeConsumers := tree.nodeConsumers.GetConsumers(twoNode.ID)
				assert.Equal(t, []int64{consumerC}, twoNodeConsumers)
			}
		})

		t.Run("logs and events case", func(t *testing.T) {
			tree := NewTree()
			tree.AddSubscription(consumerA, []string{"#", "debug"})
			tree.AddSubscription(consumerB, []string{"#"})
			tree.AddSubscription(consumerC, []string{"events", "orders", "*"})
			tree.AddSubscription(consumerD, []string{"events"})
			tree.AddSubscription(consumerE, []string{"events", "orders", "#"})

			require.Equal(t, 2, len(tree.root.Next))

			{
				hashNode := tree.root.Child("#")
				require.NotNil(t, hashNode)
				require.Equal(t, 1, len(hashNode.Next))
				assert.Equal(t, NodeTypeHash, hashNode.Type)
				assert.Equal(t, true, hashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(hashNode.ID)
				assert.Equal(t, []int64{consumerB}, firstHashNodeConsumers)

				debugNode := hashNode.Child("debug")
				require.NotNil(t, debugNode)
				require.Equal(t, 0, len(debugNode.Next))
				assert.Equal(t, NodeTypeWord, debugNode.Type)
				assert.Equal(t, true, debugNode.Stop)
				debugNodeConsumers := tree.nodeConsumers.GetConsumers(debugNode.ID)
				assert.Equal(t, []int64{consumerA}, debugNodeConsumers)
			}

			{
				eventsNode := tree.root.Child("events")
				require.NotNil(t, eventsNode)
				require.Equal(t, 1, len(eventsNode.Next))
				assert.Equal(t, NodeTypeWord, eventsNode.Type)
				assert.Equal(t, true, eventsNode.Stop)
				eventsNodeConsumers := tree.nodeConsumers.GetConsumers(eventsNode.ID)
				assert.Equal(t, []int64{consumerD}, eventsNodeConsumers)

				ordersNode := eventsNode.Child("orders")
				require.NotNil(t, ordersNode)
				require.Equal(t, 2, len(ordersNode.Next))
				assert.Equal(t, NodeTypeWord, ordersNode.Type)
				assert.Equal(t, true, ordersNode.Stop)
				ordersNodeConsumers := tree.nodeConsumers.GetConsumers(ordersNode.ID)
				assert.Equal(t, []int64{consumerE}, ordersNodeConsumers)

				ordersStarNode := ordersNode.Child("*")
				require.NotNil(t, ordersStarNode)
				require.Equal(t, 0, len(ordersStarNode.Next))
				assert.Equal(t, NodeTypeStar, ordersStarNode.Type)
				assert.Equal(t, true, ordersStarNode.Stop)
				ordersStarNodeConsumers := tree.nodeConsumers.GetConsumers(ordersStarNode.ID)
				assert.Equal(t, []int64{consumerC}, ordersStarNodeConsumers)

				ordersHashNode := ordersNode.Child("#")
				require.NotNil(t, ordersHashNode)
				require.Equal(t, 0, len(ordersHashNode.Next))
				assert.Equal(t, NodeTypeHash, ordersHashNode.Type)
				assert.Equal(t, true, ordersHashNode.Stop)
				ordersHashNodeConsumers := tree.nodeConsumers.GetConsumers(ordersHashNode.ID)
				assert.Equal(t, []int64{consumerE}, ordersHashNodeConsumers)
			}
		})

		t.Run("hash.star combination", func(t *testing.T) {
			t.Run("single hash", func(t *testing.T) {
				tree := NewTree()
				tree.AddSubscription(consumerA, []string{"#"})

				require.Equal(t, 1, len(tree.root.Next))

				hashNode := tree.root.Child("#")
				require.NotNil(t, hashNode)
				require.Equal(t, 0, len(hashNode.Next))
				assert.Equal(t, NodeTypeHash, hashNode.Type)
				assert.Equal(t, true, hashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(hashNode.ID)
				assert.Equal(t, []int64{consumerA}, firstHashNodeConsumers)
			})

			t.Run("multiple hashes", func(t *testing.T) {
				tree := NewTree()
				tree.AddSubscription(consumerA, []string{"#"})
				tree.AddSubscription(consumerB, []string{"#", "#"})

				require.Equal(t, 1, len(tree.root.Next))

				firstHashNode := tree.root.Child("#")
				require.NotNil(t, firstHashNode)
				require.Equal(t, 1, len(firstHashNode.Next))
				assert.Equal(t, NodeTypeHash, firstHashNode.Type)
				assert.Equal(t, true, firstHashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(firstHashNode.ID)
				assert.Equal(t, []int64{consumerA, consumerB}, firstHashNodeConsumers)

				secondHashNode := firstHashNode.Child("#")
				require.NotNil(t, secondHashNode)
				require.Equal(t, 0, len(secondHashNode.Next))
				assert.Equal(t, NodeTypeHash, secondHashNode.Type)
				assert.Equal(t, true, secondHashNode.Stop)
				secondHashNodeConsumers := tree.nodeConsumers.GetConsumers(secondHashNode.ID)
				assert.Equal(t, []int64{consumerB}, secondHashNodeConsumers)
			})

			t.Run("hash and dot", func(t *testing.T) {
				tree := NewTree()
				tree.AddSubscription(consumerA, []string{"#"})
				tree.AddSubscription(consumerB, []string{"#", "*"})

				require.Equal(t, 1, len(tree.root.Next))

				firstHashNode := tree.root.Child("#")
				require.NotNil(t, firstHashNode)
				require.Equal(t, 1, len(firstHashNode.Next))
				assert.Equal(t, NodeTypeHash, firstHashNode.Type)
				assert.Equal(t, true, firstHashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(firstHashNode.ID)
				assert.Equal(t, []int64{consumerA}, firstHashNodeConsumers)

				secondHashNode := firstHashNode.Child("*")
				require.NotNil(t, secondHashNode)
				require.Equal(t, 0, len(secondHashNode.Next))
				assert.Equal(t, NodeTypeStar, secondHashNode.Type)
				assert.Equal(t, true, secondHashNode.Stop)
				secondHashNodeConsumers := tree.nodeConsumers.GetConsumers(secondHashNode.ID)
				assert.Equal(t, []int64{consumerB}, secondHashNodeConsumers)
			})

			t.Run("dot and hash", func(t *testing.T) {
				tree := NewTree()
				tree.AddSubscription(consumerA, []string{"#"})
				tree.AddSubscription(consumerB, []string{"#", "*"})

				require.Equal(t, 1, len(tree.root.Next))

				firstHashNode := tree.root.Child("#")
				require.NotNil(t, firstHashNode)
				require.Equal(t, 1, len(firstHashNode.Next))
				assert.Equal(t, NodeTypeHash, firstHashNode.Type)
				assert.Equal(t, true, firstHashNode.Stop)
				firstHashNodeConsumers := tree.nodeConsumers.GetConsumers(firstHashNode.ID)
				assert.Equal(t, []int64{consumerA}, firstHashNodeConsumers)

				secondHashNode := firstHashNode.Child("*")
				require.NotNil(t, secondHashNode)
				require.Equal(t, 0, len(secondHashNode.Next))
				assert.Equal(t, NodeTypeStar, secondHashNode.Type)
				assert.Equal(t, true, secondHashNode.Stop)
				secondHashNodeConsumers := tree.nodeConsumers.GetConsumers(secondHashNode.ID)
				assert.Equal(t, []int64{consumerB}, secondHashNodeConsumers)
			})
		})
	})
}
