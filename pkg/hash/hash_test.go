package hash

import (
	"bytes"
	"hash/fnv"
	"hash/maphash"
	"math/rand"
	"testing"
)

const alphabet = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func RandStringBytesRmndr(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = alphabet[rand.Int63()%int64(len(alphabet))]
	}
	return string(b)
}

func BenchmarkHashes(b *testing.B) {
	const size = 3000000
	strings := make([]string, 0, size)
	for i := 0; i < cap(strings); i++ {
		strings = append(strings, RandStringBytesRmndr(20))
	}

	b.ResetTimer()

	b.Run("fnv", func(b *testing.B) {
		b.ReportAllocs()
		h := fnv.New64()

		for i := 0; i < b.N; i++ {
			if i > size-1 {
				break
			}

			h.Reset()
			h.Write([]byte(strings[i]))
			h.Sum64()
		}
	})

	b.Run("maphash", func(b *testing.B) {
		b.ReportAllocs()
		h := maphash.Hash{}

		for i := 0; i < b.N; i++ {
			if i > size-1 {
				break
			}

			h.Reset()
			h.Write([]byte(strings[i]))
			h.Sum64()
		}
	})
}

func BenchmarkIndex(b *testing.B) {
	const size = 3000000
	strings := make([]string, 0, size)
	for i := 0; i < cap(strings); i++ {
		strings = append(strings, RandStringBytesRmndr(20))
	}

	b.ResetTimer()
	b.ReportAllocs()

	var idx int
	for i := 0; i < b.N; i++ {
		idx = i
		if i > size-1 {
			idx = i % size
		}

		bytes.IndexByte([]byte(strings[idx]), '1')
	}
}
