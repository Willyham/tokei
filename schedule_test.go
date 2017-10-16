package tokei

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var epoch = time.Unix(0, 0).In(time.UTC)

func TestNext(t *testing.T) {
	cases := []struct {
		name      string
		input     string
		startTime time.Time
		expected  time.Time
	}{
		{"every minute", "* * * * *", epoch, epoch},
		{"minute 5", "5 * * * *", epoch, epoch.Add(time.Minute * 5)},
		{"hour 5", "* 5 * * *", epoch, epoch.Add(time.Hour * 5)},
		{"day of month 5", "* * 5 * *", epoch, epoch.AddDate(0, 0, 4)},
		{"day of week 5", "* * * * 5", epoch, epoch.AddDate(0, 0, 1)}, // Epoch was a Thursday
		{"month 5", "* * * 5 *", epoch, epoch.AddDate(0, 4, 0)},

		// Tuesday in Jan which is the 2nd of the month. This is before the epoch (end of 1969) and doesn't occur again
		// until 1973
		{"wrap year", "* * 2 1 2", epoch, epoch.AddDate(3, 0, 1)},

		// Start at day 10 and ask for any 7th of the month. Should get Feb 7th.
		{"wrap months", "* * 7 * *", epoch.AddDate(0, 0, 10), epoch.AddDate(0, 1, 6)},

		// Start at hour 10 on 1st and ask for any hour 7. Should get 07:00 Jan 2nd.
		{"wrap hours", "* 7 * * *", epoch.Add(time.Hour * 10), epoch.AddDate(0, 0, 1).Add(time.Hour * 7)},

		// Start at minute 10 in hour 3 and ask for any mminute 7. Should get 04:07 Jan 1st.
		{"wrap minute", "7 * * * *", epoch.Add(time.Hour*3 + time.Minute*10), epoch.Add(time.Hour*4 + time.Minute*7)},

		// At every 5th minute from 10 through 59 past every hour from 3 through 5 on day-of-month 1 and 2 and on Tuesday in July.
		// First occurrence is Tuesday 2nd July 1974 03:10:00
		{"complex", "10/5 3-5 1,2 7 2", epoch, epoch.AddDate(4, 6, 1).Add(time.Hour*3 + time.Minute*10)},
	}

	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {
			ex, err := Parse(test.input)
			require.NoError(t, err)

			sched := NewScheduleUTC(ex)
			output := sched.NextFrom(test.startTime)
			assert.Equal(t, test.expected, output)
		})
	}
}

func TestNextMultiple(t *testing.T) {
	ex, err := Parse("* * * * *")
	require.NoError(t, err)

	sched := NewScheduleUTC(ex)
	next := sched.ProjectFrom(epoch, 3)

	expected := []time.Time{epoch, epoch.Add(time.Minute * 1), epoch.Add(time.Minute * 2)}
	assert.Equal(t, expected, next)
}

func TestTimer(t *testing.T) {
	ex, err := Parse("* * * * *")
	require.NoError(t, err)
	schedule := NewScheduleUTC(ex)

	startTime := time.Now()
	timer := schedule.Timer()
	go timer.Start()

	// Should immediately receive as current time is matched
	out := <-timer.Next()
	assert.WithinDuration(t, startTime, out, time.Minute)
}

func TestTimerLong(t *testing.T) {
	if testing.Short() {
		t.SkipNow()
		return
	}
	ex, err := Parse("*/2 * * * *")
	require.NoError(t, err)
	schedule := NewScheduleUTC(ex)

	startTime := time.Now()
	timer := schedule.Timer()
	go timer.Start()

	out := <-timer.Next()
	// Should fire within the next minute (+ a second buffer.)
	assert.WithinDuration(t, startTime, out, time.Minute+time.Second)
}

var benchCases = []struct {
	name  string
	input string
}{
	{"all", "* * * * *"},
	{"every 10 minutes", "*/10 * * * *"},
	{"waking hours", "00 09-18 * * 1-5"},
	{"years in future", "10/5 3-5 1,2 7 2"},
}

func BenchmarkNext(b *testing.B) {
	for _, bench := range benchCases {
		ex, err := Parse(bench.input)
		require.NoError(b, err)
		sched := NewScheduleUTC(ex)

		epoch := time.Unix(0, 0).In(time.UTC)
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sched.NextFrom(epoch)
			}
		})
	}
}

func BenchmarkProject(b *testing.B) {
	for _, bench := range benchCases {
		ex, err := Parse(bench.input)
		require.NoError(b, err)
		sched := NewScheduleUTC(ex)

		epoch := time.Unix(0, 0).In(time.UTC)
		b.Run(bench.name, func(b *testing.B) {
			for i := 0; i < b.N; i++ {
				sched.ProjectFrom(epoch, 5)
			}
		})
	}
}
