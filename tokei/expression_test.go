package tokei

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKleeneExpression(t *testing.T) {
	ex, err := KleeneExpression(DayOfWeekContext, "*")
	assert.NoError(t, err)
	assert.Equal(t, Sequence{start: 1, stop: 7, step: 1}, ex)
}

func TestKleeneExpressionError(t *testing.T) {
	_, err := KleeneExpression(MinuteContext, "5")
	assert.Error(t, err)
}

func TestRangeExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected Enumerator
	}{
		{"normal", "1-10", Sequence{start: 1, stop: 10, step: 1}},
		{"single", "25", IrregularSequence{entries: []int{25}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := RangeExpression(MinuteContext, test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, re)
		})
	}
}

func TestRangeExpressionErrors(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"too many values", "1-5-10"},
		{"star range", "*-10"},
		{"bad range", "20-10"},
		{"bad start", "a-10"},
		{"bad end", "10-a"},
		{"junk", "asdfghjkl"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := RangeExpression(MinuteContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestRepeatExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected Enumerator
	}{
		{"star", "*/10", Sequence{start: 0, stop: 59, step: 10}},
		{"normal", "5/10", Sequence{start: 5, stop: 59, step: 10}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := RepeatExpression(MinuteContext, test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, re)
		})
	}
}

func TestRepeatExpressionError(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"too many parts", "10/10/10"},
		{"bad star", "10/*"},
		{"bad start value", "a/10"},
		{"bad end value", "10/a"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := RepeatExpression(MinuteContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestLiteralExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected Enumerator
	}{
		{"single", "1,2,3", IrregularSequence{entries: []int{1, 2, 3}}},
		{"normal", "1,2,3", IrregularSequence{entries: []int{1, 2, 3}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := LiteralExpression(MinuteContext, test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, re)
		})
	}
}

func TestLiteralExpressionError(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"bad value", "a"},
		{"bad value multi", "1,2,foo"},
		{"bad format", "1|2"},
		{"too large", "1,2,1000"},
		{"too small", "-1,1"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := LiteralExpression(MinuteContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestMultiExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected Enumerator
	}{
		{"star", "*", Sequence{start: 0, stop: 59, step: 1}},
		{"range", "1-10", Sequence{start: 1, stop: 10, step: 1}},
		{"range single", "25", IrregularSequence{entries: []int{25}}},
		{"repeat", "5/10", Sequence{start: 5, stop: 59, step: 10}},
		{"repeat star", "*/10", Sequence{start: 0, stop: 59, step: 10}},
		{"literal", "1,2,3", IrregularSequence{entries: []int{1, 2, 3}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := multiExpression.Parse(MinuteContext, test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, re)
		})
	}
}

func TestMultiError(t *testing.T) {
	_, err := multiExpression.Parse(MinuteContext, "blah")
	assert.Error(t, err)
}
