package alloc

import (
	"github.com/ivanovaleksey/broker/src/node"
	"testing"
)

func BenchmarkAlloc(b *testing.B) {
	b.Run("map alloc", func(b *testing.B) {
		b.ReportAllocs()
		n := node.Node{}
		for i := 0; i < b.N; i++ {
			m := make(map[int64]*node.Node, 10)
			m[0] = &n
		}
	})

	b.Run("slice alloc", func(b *testing.B) {
		b.ReportAllocs()
		n := node.Node{}
		for i := 0; i < b.N; i++ {
			m := make([]*node.Node, 10)
			m[0] = &n
		}
	})
}
