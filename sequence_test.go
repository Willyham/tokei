package tokei

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSequenceEnumerate(t *testing.T) {
	cases := []struct {
		name             string
		start, end, step int
		expected         []int
		errExpected      bool
	}{
		{"monotonic", 0, 5, 1, []int{0, 1, 2, 3, 4, 5}, false},
		{"arithmetic", 0, 10, 2, []int{0, 2, 4, 6, 8, 10}, false},
		{"odd", 0, 9, 2, []int{0, 2, 4, 6, 8}, false},
		{"bad bounds", 0, -10, 1, nil, true},
		{"negative step", 0, 10, -1, nil, true},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			seq, err := newSequence(test.start, test.end, test.step)
			if test.errExpected {
				assert.Error(t, err)
				assert.Nil(t, seq)
				return
			}
			assert.Equal(t, test.expected, seq.Enumerate())
		})
	}
}

func TestIrregularSequence(t *testing.T) {
	assert.Equal(t, []int{1, 3, 4, 10}, newIrregularSequence([]int{1, 3, 4, 10}).Enumerate())
	assert.Equal(t, []int{1, 3, 4, 10}, newIrregularSequence([]int{1, 10, 3, 4}).Enumerate())
	assert.Equal(t, []int{1, 2, 3}, newIrregularSequence([]int{1, 1, 2, 2, 3, 3}).Enumerate())
}
