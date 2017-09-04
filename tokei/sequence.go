package tokei

import (
	"errors"
	"sort"
)

type Enumerator interface {
	Enumerate() []int
}

type Sequence struct {
	start int
	stop  int
	step  int
}

func NewSequence(start, stop, step int) (*Sequence, error) {
	if step < 0 || stop < start {
		return nil, errors.New("invalid sequence")
	}
	return &Sequence{
		start: start,
		stop:  stop,
		step:  step,
	}, nil
}

func (s Sequence) Enumerate() []int {
	output := make([]int, 0)
	for i := s.start; i <= s.stop; i += s.step {
		output = append(output, i)
	}
	return output
}

type IrregularSequence struct {
	entries []int
}

func NewIrregularSequence(entries []int) IrregularSequence {
	return IrregularSequence{
		entries: entries,
	}
}

func (s IrregularSequence) Enumerate() []int {
	sort.Ints(s.entries)
	return s.entries
}
