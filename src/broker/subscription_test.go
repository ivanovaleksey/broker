package broker

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func Test_prepareParts(t *testing.T) {
	tests := []struct {
		in       []string
		expected []string
	}{
		{
			in:       []string{"1", "2"},
			expected: []string{"1", "2"},
		},
		{
			in:       []string{"#", "1", "2"},
			expected: []string{"#", "1", "2"},
		},
		{
			in:       []string{"#", "#", "1", "*", "2"},
			expected: []string{"#", "1", "*", "2"},
		},
		{
			in:       []string{"#", "#", "1", "*", "2", "#", "#"},
			expected: []string{"#", "1", "*", "2", "#"},
		},
	}

	for _, test := range tests {
		out := prepareParts(test.in)

		assert.Equal(t, test.expected, out, test.in)
	}
}
