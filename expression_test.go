package tokei

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKleeneExpression(t *testing.T) {
	ex, err := kleeneExpression(dayOfWeekContext, "*")
	assert.NoError(t, err)
	assert.Equal(t, sequence{start: 1, end: 7, step: 1}, ex)
}

func TestKleeneExpressionError(t *testing.T) {
	_, err := kleeneExpression(minuteContext, "5")
	assert.Error(t, err)
}

func TestRangeExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected enumerator
	}{
		{"normal", "1-10", sequence{start: 1, end: 10, step: 1}},
		{"single", "25", irregularSequence{entries: []int{25}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := rangeExpression(minuteContext, test.input)
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
		{"low start", "0-10"},
		{"high end", "10-13"},
		{"high both", "13-14"},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := rangeExpression(monthContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestRepeatExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected enumerator
	}{
		{"star", "*/10", sequence{start: 0, end: 59, step: 10}},
		{"normal", "5/10", sequence{start: 5, end: 59, step: 10}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := repeatExpression(minuteContext, test.input)
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
			_, err := repeatExpression(minuteContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestLiteralExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected enumerator
	}{
		{"single", "1,2,3", irregularSequence{entries: []int{1, 2, 3}}},
		{"normal", "1,2,3", irregularSequence{entries: []int{1, 2, 3}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := literalExpression(minuteContext, test.input)
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
			_, err := literalExpression(minuteContext, test.input)
			assert.Error(t, err)
		})
	}
}

func TestMultiExpression(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected enumerator
	}{
		{"star", "*", sequence{start: 0, end: 59, step: 1}},
		{"range", "1-10", sequence{start: 1, end: 10, step: 1}},
		{"range single", "25", irregularSequence{entries: []int{25}}},
		{"repeat", "5/10", sequence{start: 5, end: 59, step: 10}},
		{"repeat star", "*/10", sequence{start: 0, end: 59, step: 10}},
		{"literal", "1,2,3", irregularSequence{entries: []int{1, 2, 3}}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			re, err := defaultMultiExpression.Parse(minuteContext, test.input)
			assert.NoError(t, err)
			assert.Equal(t, test.expected, re)
		})
	}
}

func TestMultiError(t *testing.T) {
	_, err := defaultMultiExpression.Parse(minuteContext, "blah")
	assert.Error(t, err)
}

func TestContextInvalid(t *testing.T) {
	invalid := expressionContext(1000)
	assert.Panics(t, func() {
		invalid.Min()
	})
	assert.Panics(t, func() {
		invalid.Max()
	})
}

func TestParseInvalid(t *testing.T) {
	cases := []struct {
		name  string
		input string
	}{
		{"empty", ""},
		{"not enough parts", "1 1 1 1"},
		{"bad part", "a 1 1 1 1"},
	}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			_, err := Parse(test.input)
			assert.Error(t, err)
		})
	}
}
