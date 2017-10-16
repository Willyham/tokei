package tokei

import (
	"errors"
	"sort"
)

// enumerator is anything which can enumerate a list of ints.
type enumerator interface {
	Enumerate() []int
}

// sequence describes a sequence of ints which can be enumerated.
type sequence struct {
	start int
	end   int
	step  int
}

// newSequence creates a new sequence.
func newSequence(start, end, step int) (*sequence, error) {
	if step < 0 || end < start {
		return nil, errors.New("invalid sequence")
	}
	return &sequence{
		start: start,
		end:   end,
		step:  step,
	}, nil
}

// Enumerate returns a list of all ints in this sequence.
func (s sequence) Enumerate() []int {
	output := make([]int, 0)
	for i := s.start; i <= s.end; i += s.step {
		output = append(output, i)
	}
	return output
}

// irregularSequence is a sequence which doesn't follow a regular pattern.
type irregularSequence struct {
	entries []int
}

// newIrregularSequence creates a sequence from a list of ints.
func newIrregularSequence(entries []int) irregularSequence {
	return irregularSequence{
		entries: entries,
	}
}

// Enumerate returns the list of ints, in order.
// It removes any duplicates.
func (s irregularSequence) Enumerate() []int {
	output := make([]int, 0, len(s.entries))
	seen := map[int]struct{}{}
	for _, num := range s.entries {
		_, ok := seen[num]
		if ok {
			continue
		}
		output = append(output, num)
		seen[num] = struct{}{}
	}

	sort.Ints(output)
	return output
}
