package tokei

import (
	"errors"
	"regexp"
	"strconv"
	"strings"
)

type CronExpression struct {
	minutes    Enumerator
	hours      Enumerator
	dayOfMonth Enumerator
	month      Enumerator
	dayOfWeek  Enumerator
	command    string
}

func Parse(input string) (*CronExpression, error) {
	parts := strings.Split(input, " ")
	if len(parts) < 6 {
		return nil, errors.New("invalid expression")
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
		minutes:    min,
		hours:      hour,
		dayOfMonth: dom,
		month:      month,
		dayOfWeek:  dow,
		command:    command,
	}, nil
}

type Parser interface {
	Parse(ExpressionContext, string) (Enumerator, error)
}

type ParseFunc func(ExpressionContext, string) (Enumerator, error)

func (f ParseFunc) Parse(ex ExpressionContext, input string) (Enumerator, error) {
	return f(ex, input)
}

type MultiExpression struct {
	rangeRegex   *regexp.Regexp
	repeatRegex  *regexp.Regexp
	literalRegex *regexp.Regexp
}

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

var KleeneExpression = ParseFunc(func(ex ExpressionContext, input string) (Enumerator, error) {
	if input != "*" {
		return nil, errors.New("input must be *")
	}
	return Sequence{
		start: ex.Min(),
		stop:  ex.Max(),
		step:  1,
	}, nil
})

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
		stop:  end,
		step:  1,
	}, nil
})

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
			stop:  ex.Max(),
			step:  end,
		}, nil
	}

	start, err := parseStartValue(ex, parts[0])
	if err != nil {
		return nil, err
	}

	return Sequence{
		start: start,
		stop:  ex.Max(),
		step:  end,
	}, nil
})

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
