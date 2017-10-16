package tokei

import (
	"errors"
	"sort"
)

// Enumerator is anything which can enumerate a list of ints.
type Enumerator interface {
	Enumerate() []int
}

// Sequence describes a sequence of ints which can be enumerated.
type Sequence struct {
	start int
	end   int
	step  int
}

// NewSequence creates a new sequence.
func NewSequence(start, end, step int) (*Sequence, error) {
	if step < 0 || end < start {
		return nil, errors.New("invalid sequence")
	}
	return &Sequence{
		start: start,
		end:   end,
		step:  step,
	}, nil
}

// Enumerate returns a list of all ints in this sequence.
func (s Sequence) Enumerate() []int {
	output := make([]int, 0)
	for i := s.start; i <= s.end; i += s.step {
		output = append(output, i)
	}
	return output
}

// IrregularSequence is a sequence which doesn't follow a regular pattern.
type IrregularSequence struct {
	entries []int
}

// NewIrregularSequence creates a sequence from a list of ints.
func NewIrregularSequence(entries []int) IrregularSequence {
	return IrregularSequence{
		entries: entries,
	}
}

// Enumerate returns the list of ints, in order.
// It removes any duplicates.
func (s IrregularSequence) Enumerate() []int {
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
