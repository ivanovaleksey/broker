package topics

import (
	"github.com/ivanovaleksey/broker/pkg/hash"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBytesTopic_Parse(t *testing.T) {
	tests := []struct {
		in  string
		exp []string
	}{
		{
			in:  "a.bb.cd",
			exp: []string{"a", "bb", "cd"},
		},
		{
			in:  "a.#.*.d",
			exp: []string{"a", "#", "*", "d"},
		},
		{
			in:  "a.#.#.d",
			exp: []string{"a", "#", "d"},
		},
		{
			in:  "#.#.#.d",
			exp: []string{"#", "d"},
		},
		{
			in:  "#.#.a.d",
			exp: []string{"#", "a", "d"},
		},
		{
			in:  "#.a.#.d",
			exp: []string{"#", "a", "#", "d"},
		},
	}

	h := hash.GetHash()
	p := NewBytesParser()

	for _, test := range tests {
		t.Run(test.in, func(t *testing.T) {
			exp := make([]uint64, len(test.exp))
			for i, s := range test.exp {
				h.WriteString(s)
				exp[i] = h.Sum64()
				h.Reset()
			}

			got := p.Parse(test.in)

			assert.Equal(t, exp, got)
		})
	}
}
