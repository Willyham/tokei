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
	minutes    enumerator
	hours      enumerator
	dayOfMonth enumerator
	month      enumerator
	dayOfWeek  enumerator
	command    string
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
	min, minErr := defaultMultiExpression.Parse(minuteContext, parts[0])
	hour, hourErr := defaultMultiExpression.Parse(hourContext, parts[1])
	dom, domErr := defaultMultiExpression.Parse(dayOfMonthContext, parts[2])
	month, monthErr := defaultMultiExpression.Parse(monthContext, parts[3])
	dow, dowErr := defaultMultiExpression.Parse(dayOfWeekContext, parts[4])
	command := strings.Join(parts[5:], " ")

	for _, err := range []error{minErr, hourErr, domErr, monthErr, dowErr} {
		if err != nil {
			return nil, err
		}
	}

	return &CronExpression{
		minutes:    min,
		hours:      hour,
		dayOfMonth: dom,
		month:      month,
		dayOfWeek:  dow,
		command:    command,
	}, nil
}

// parser is anything that can parse an expression part.
type parser interface {
	Parse(expressionContext, string) (enumerator, error)
}

// parseFunc allows us to adapt a func to Parser.
type parseFunc func(expressionContext, string) (enumerator, error)

// Parse adapts ParseFunc to parser.
func (f parseFunc) Parse(ex expressionContext, input string) (enumerator, error) {
	return f(ex, input)
}

// A multiExpression detects the type of expression and hands off to another parser.
// It uses regex to do this for simplicity rather than a traditional tokenizer.
type multiExpression struct {
	rangeRegex   *regexp.Regexp
	repeatRegex  *regexp.Regexp
	literalRegex *regexp.Regexp
}

// Parse parses any expression by deferring to other parsers.
func (m multiExpression) Parse(ex expressionContext, input string) (enumerator, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "*" {
		return kleeneExpression(ex, input)
	}
	if m.rangeRegex.MatchString(trimmed) {
		return rangeExpression(ex, trimmed)
	}
	if m.repeatRegex.MatchString(trimmed) {
		return repeatExpression(ex, trimmed)
	}
	if m.literalRegex.MatchString(trimmed) {
		return literalExpression(ex, trimmed)
	}
	return nil, errors.New("unknown expression")
}

var defaultMultiExpression = multiExpression{
	rangeRegex:   regexp.MustCompile(`\d-\d`),
	repeatRegex:  regexp.MustCompile(`./\d`),
	literalRegex: regexp.MustCompile(`(\d+)(,\s*\d+)*`),
}

// kleeneExpression parses the "*" expression only.
var kleeneExpression = parseFunc(func(ex expressionContext, input string) (enumerator, error) {
	if input != "*" {
		return nil, errors.New("input must be *")
	}
	return sequence{
		start: ex.Min(),
		end:   ex.Max(),
		step:  1,
	}, nil
})

// rangeExpression parses expressions of the for x-y.
var rangeExpression = parseFunc(func(ex expressionContext, input string) (enumerator, error) {
	parts := strings.Split(input, "-")
	if len(parts) > 2 {
		return nil, errors.New("must be of form x-y")
	}
	start, err := parseStartValue(ex, parts[0])
	if err != nil {
		return nil, err
	}

	if len(parts) == 1 {
		return newIrregularSequence([]int{start}), nil
	}

	end, err := parseEndValue(ex, parts[1])
	if err != nil {
		return nil, err
	}

	if start > end {
		return nil, errors.New("invalid range")
	}

	return sequence{
		start: start,
		end:   end,
		step:  1,
	}, nil
})

// repeatExpression parses expressions of the form x/y, including */y.
var repeatExpression = parseFunc(func(ex expressionContext, input string) (enumerator, error) {
	parts := strings.Split(input, "/")
	if len(parts) != 2 {
		return nil, errors.New("Invalid repeat expression, must be of form x/y")
	}
	end, err := parseEndValue(ex, parts[1])
	if err != nil {
		return nil, err
	}

	if parts[0] == "*" {
		return sequence{
			start: ex.Min(),
			end:   ex.Max(),
			step:  end,
		}, nil
	}

	start, err := parseStartValue(ex, parts[0])
	if err != nil {
		return nil, err
	}

	return sequence{
		start: start,
		end:   ex.Max(),
		step:  end,
	}, nil
})

// literalExpression parses expressions of the form "x,y[,z].."
var literalExpression = parseFunc(func(ex expressionContext, input string) (enumerator, error) {
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

	return irregularSequence{
		entries: times,
	}, nil
})

func parseStartValue(ex expressionContext, input string) (int, error) {
	start, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, err
	}
	if start < ex.Min() {
		return 0, errors.New("invalid start value")
	}
	return start, nil
}

func parseEndValue(ex expressionContext, input string) (int, error) {
	end, err := strconv.Atoi(strings.TrimSpace(input))
	if err != nil {
		return 0, err
	}
	if end > ex.Max() {
		return 0, errors.New("invalid end value")
	}
	return end, nil
}
