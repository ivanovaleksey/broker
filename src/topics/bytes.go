package topics

import (
	"bytes"
	"github.com/ivanovaleksey/broker/pkg/hash"
)

const (
	byteStar byte = '*'
	byteHash byte = '#'
	byteDot  byte = '.'
)

var (
	HashStar uint64
	HashHash uint64

	bytesStar = [1]byte{byteStar}
	bytesHash = [1]byte{byteHash}
	bytesDot  = [1]byte{byteDot}
)

func init() {
	h := hash.GetHash()

	h.WriteByte(byteStar)
	HashStar = h.Sum64()
	h.Reset()

	h.WriteByte(byteHash)
	HashHash = h.Sum64()
	h.Reset()

	hash.PutHash(h)
}

type BytesParser struct {
}

func NewBytesParser() *BytesParser {
	return &BytesParser{}
}

func (t *BytesParser) IsStatic(str string) bool {
	for i := 0; i < len(str); i++ {
		if str[i] == byteHash || str[i] == byteStar {
			return false
		}
	}
	return true
}

type Info struct {
	IsStatic bool
	Part     uint64
	Parts    []uint64
}

func (t *BytesParser) Parse(str string) Info {
	if t.IsStatic(str) {
		return Info{
			IsStatic: true,
			Part:     t.Hash(str),
		}
	}

	return Info{
		IsStatic: false,
		Parts:    t.Parts(str),
	}
}

func (t *BytesParser) Hash(str string) uint64 {
	h := hash.GetHash()
	h.WriteString(str)
	sum := h.Sum64()
	hash.PutHash(h)
	return sum
}

func (t *BytesParser) Parts(str string) []uint64 {
	h := hash.GetHash()
	parts := bytes.Split([]byte(str), bytesDot[:])
	out := make([]uint64, 0, len(parts))
	for i := 0; i < len(parts); i++ {
		if bytes.Equal(parts[i], bytesHash[:]) {
			nextIdx := i + 1
			if nextIdx < len(parts) && bytes.Equal(parts[nextIdx], bytesHash[:]) {
				// skip sequent #
				continue
			}
		}

		h.Write(parts[i])
		out = append(out, h.Sum64())
		h.Reset()
	}
	hash.PutHash(h)
	return out
}
