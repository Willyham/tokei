// Package tokei provides a cron parser and scheduler.
//
// Tokei works by parsing the cron string and generating an Enumerator for each
// part of the expression. It these uses these Enumerators to enumerate possible valid
// combinations of times which match the expression.
package tokei

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

// CronExpression describes a parsed cron expression.
type CronExpression struct {
	Minutes    Enumerator
	Hours      Enumerator
	DayOfMonth Enumerator
	Month      Enumerator
	DayOfWeek  Enumerator
	Command    string
}

// Parse parses a cron expression from a string.
func Parse(input string) (*CronExpression, error) {
	parts := strings.Split(input, " ")
	if len(parts) < 5 {
		return nil, errors.New("invalid expression")
	}
	if len(parts) == 5 {
		parts = append(parts, "")
	}
	min, minErr := multiExpression.Parse(MinuteContext, parts[0])
	hour, hourErr := multiExpression.Parse(HourContext, parts[1])
	dom, domErr := multiExpression.Parse(DayOfMonthContext, parts[2])
	month, monthErr := multiExpression.Parse(MonthContext, parts[3])
	dow, dowErr := multiExpression.Parse(DayOfWeekContext, parts[4])
	command := strings.Join(parts[5:], " ")

	for _, err := range []error{minErr, hourErr, domErr, monthErr, dowErr} {
		if err != nil {
			return nil, err
		}
	}

	return &CronExpression{
		Minutes:    min,
		Hours:      hour,
		DayOfMonth: dom,
		Month:      month,
		DayOfWeek:  dow,
		Command:    command,
	}, nil
}

// Parser is anything that can parse an expression part.
type Parser interface {
	Parse(ExpressionContext, string) (Enumerator, error)
}

// ParseFunc allows us to adapt a func to Parser.
type ParseFunc func(ExpressionContext, string) (Enumerator, error)

// Parse adapts ParseFunc to parser.
func (f ParseFunc) Parse(ex ExpressionContext, input string) (Enumerator, error) {
	return f(ex, input)
}

// A MultiExpression detects the type of expression and hands off to another parser.
// It uses regex to do this for simplicity rather than a traditional tokenizer.
type MultiExpression struct {
	rangeRegex   *regexp.Regexp
	repeatRegex  *regexp.Regexp
	literalRegex *regexp.Regexp
}

// Parse parses any expression by deferring to other parsers.
func (m MultiExpression) Parse(ex ExpressionContext, input string) (Enumerator, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "*" {
		return KleeneExpression(ex, input)
	}
	if m.rangeRegex.MatchString(trimmed) {
		return RangeExpression(ex, trimmed)
	}
	if m.repeatRegex.MatchString(trimmed) {
		return RepeatExpression(ex, trimmed)
	}
	if m.literalRegex.MatchString(trimmed) {
		return LiteralExpression(ex, trimmed)
	}
	return nil, errors.New("unknown expression")
}

var multiExpression = MultiExpression{
	rangeRegex:   regexp.MustCompile(`\d-\d`),
	repeatRegex:  regexp.MustCompile(`./\d`),
	literalRegex: regexp.MustCompile(`(\d+)(,\s*\d+)*`),
}

// KleeneExpression parses the "*" expression only.
var KleeneExpression = ParseFunc(func(ex ExpressionContext, input string) (Enumerator, error) {
	if input != "*" {
		return nil, errors.New("input must be *")
	}
	return Sequence{
		start: ex.Min(),
		end:   ex.Max(),
		step:  1,
	}, nil
})

// RangeExpression parses expressions of the for x-y.
var RangeExpression = ParseFunc(func(ex ExpressionContext, input string) (Enumerator, error) {
	parts := strings.Split(input, "-")
	if len(parts) > 2 {
		return nil, errors.New("must be of form x-y")
	}
	start, err := parseStartValue(ex, parts[0])
	if err != nil {
		return nil, err
	}

	if len(parts) == 1 {
		return NewIrregularSequence([]int{start}), nil
	}

	end, err := parseEndValue(ex, parts[1])
	if err != nil {
		return nil, err
	}

	if start > end {
		return nil, errors.New("invalid range")
	}

	return Sequence{
		start: start,
		end:   end,
		step:  1,
	}, nil
})

// RepeatExpression parses expressions of the form x/y, including */y.
var RepeatExpression = ParseFunc(func(ex ExpressionContext, input string) (Enumerator, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return nil, errors.New("Invalid repeat expression, must be of form x/y")
	}
	end, err := parseEndValue(ex, parts[1])
	if err != nil {
		return nil, err
	}

	if parts[0] == "*" {
		return Sequence{
			start: ex.Min(),
			end:   ex.Max(),
			step:  end,
		}, nil
	}

	start, err := parseStartValue(ex, parts[0])
	if err != nil {
		return nil, err
	}

	return Sequence{
		start: start,
		end:   ex.Max(),
		step:  end,
	}, nil
})

// LiteralExpression parses expressions of the form "x,y[,z].."
var LiteralExpression = ParseFunc(func(ex ExpressionContext, input string) (Enumerator, error) {
	parts := strings.Split(input, ",")
	times := make([]int, len(parts))
	for i, part := range parts {
		time, err := strconv.Atoi(strings.TrimSpace(part))
		if err != nil {
			return nil, err
		}
		if time < ex.Min() || time > ex.Max() {
			return nil, errors.New("invalid time part for literal expression")
		}
		times[i] = time
	}

	return IrregularSequence{
		entries: times,
	}, nil
})

func parseStartValue(ex ExpressionContext, input string) (int, error) {
	start, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, err
	}
	if start < ex.Min() {
		return 0, errors.New("invalid start value")
	}
	return start, nil
}

func parseEndValue(ex ExpressionContext, input string) (int, error) {
	end, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, err
	}
	if end > ex.Max() {
		return 0, errors.New("invalid end value")
	}
	return end, nil
}
