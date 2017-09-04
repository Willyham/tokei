package tokei

import (
	"io"
	"strconv"
	"strings"
)

const tableColWidth = 14

type Reporter interface {
	Report(*CronExpression, io.Writer)
}

type reportFunc func(*CronExpression, io.Writer)

func (f reportFunc) Report(c *CronExpression, w io.Writer) {
	f(c, w)
}

var TableReporter = reportFunc(func(c *CronExpression, w io.Writer) {
	parts := []struct {
		label string
		enum  Enumerator
	}{
		{"minute", c.minutes},
		{"hours", c.hours},
		{"day of month", c.dayOfMonth},
		{"month", c.month},
		{"day of week", c.dayOfWeek},
	}

	for _, part := range parts {
		w.Write([]byte(padRight(part.label, " ", tableColWidth)))
		w.Write([]byte(separatedSlice(part.enum.Enumerate())))
		w.Write([]byte("\n"))
	}

	w.Write([]byte(padRight("command", " ", tableColWidth)))
	w.Write([]byte(c.command))
	w.Write([]byte("\n"))
})

func separatedSlice(nums []int) string {
	asStrings := make([]string, len(nums))
	for i, s := range nums {
		asStrings[i] = strconv.Itoa(s)
	}
	return strings.Join(asStrings, " ")
}

func padRight(str, pad string, length int) string {
	for {
		str += pad
		if len(str) > length {
			return str[0:length]
		}
	}
}
