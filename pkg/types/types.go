package types

import (
	"hash/maphash"
	"sync"
)

type ConsumerID = int64

type NodeID = int64

type Topic string

var seed maphash.Seed

func init() {
	seed = maphash.MakeSeed()
	for i := 0; i < 10000000; i++ {
		hashPool.Put(new(maphash.Hash))
	}
}

var hashPool = sync.Pool{
	New: func() interface{} {
		return new(maphash.Hash)
	},
}

func (t Topic) Hash() uint64 {
	h := hashPool.Get().(*maphash.Hash)
	defer func() {
		h.Reset()
		hashPool.Put(h)
	}()
	h.SetSeed(seed)
	h.WriteString(string(t))
	return h.Sum64()
}
