package hash

import (
	"hash/maphash"
	"sync"
)

var seed maphash.Seed

func init() {
	const numHashes = 10000000

	seed = maphash.MakeSeed()
	for i := 0; i < numHashes; i++ {
		hashPool.Put(new(maphash.Hash))
	}
}

var hashPool = sync.Pool{
	New: func() interface{} {
		return new(maphash.Hash)
	},
}

func GetHash() *maphash.Hash {
	h := hashPool.Get().(*maphash.Hash)
	h.SetSeed(seed)
	return h
}

func PutHash(h *maphash.Hash) {
	h.Reset()
	hashPool.Put(h)
}
