package hash

import (
	"hash/maphash"
	"sync"
)

var Seed maphash.Seed

func init() {
	Seed = maphash.MakeSeed()
}

var PoolCount int32
var Pool = sync.Pool{
	New: func() interface{} {
		// atomic.AddInt32(&PoolCount, 1)
		return new(maphash.Hash)
	},
}

func GetHash() *maphash.Hash {
	h := Pool.Get().(*maphash.Hash)
	h.SetSeed(Seed)
	return h
}

func PutHash(h *maphash.Hash) {
	h.Reset()
	Pool.Put(h)
}
