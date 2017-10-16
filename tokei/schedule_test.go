package tokei

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNext(t *testing.T) {
	epoch := time.Unix(0, 0).In(time.UTC)
	cases := []struct {
		name     string
		input    string
		n        int
		expected []time.Time
	}{
		{"every minute", "* * * * * /bin/foo", 2, []time.Time{epoch, epoch.Add(time.Minute)}},
		{"minute 5", "5 * * * * /bin/foo", 1, []time.Time{epoch.Add(time.Minute * 5)}},
		{"hour 5", "* 5 * * * /bin/foo", 1, []time.Time{epoch.Add(time.Hour * 5)}},
		{"day of month 5", "* * 5 * * /bin/foo", 1, []time.Time{epoch.AddDate(0, 0, 4)}},
		{"day of week 5", "* * * * 5 /bin/foo", 1, []time.Time{epoch.AddDate(0, 0, 1)}}, // Epoch was a Thursday
		{"month 5", "* * * 5 * /bin/foo", 1, []time.Time{epoch.AddDate(0, 4, 0)}},

		// Tuesday in Jan which is the 2nd of the month. This is before the epoch (end of 1969) and doesn't occur again
		// until 1973
		{"wrap", "* * 2 1 2 /bin/foo", 1, []time.Time{epoch.AddDate(3, 0, 1)}},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

			ex, err := Parse(test.input)
			require.NoError(t, err)

			sched := NewScheduleUTC(ex)
			output := sched.NextFrom(epoch, test.n)
			assert.Equal(t, test.expected, output)
		})
	}
}
