package tokei

import (
	"sort"
	"time"
)

// Schedule represents the schedule on which the job will fire for a given timezone.
type Schedule struct {
	location   *time.Location
	expression *CronExpression
}

// NewSchedule creates a new schedule for an expression in the given timezone.
func NewSchedule(location *time.Location, ex *CronExpression) *Schedule {
	return &Schedule{
		location:   location,
		expression: ex,
	}
}

// NewScheduleUTC creates a new schedule for the expression in UTC.
func NewScheduleUTC(ex *CronExpression) *Schedule {
	return NewSchedule(time.UTC, ex)
}

// Timer returns a ScheduleTimer which fires on this schedule.
func (s *Schedule) Timer() *ScheduleTimer {
	return NewScheduleTimer(s)
}

// Next returns the next time that matches the schedule.
func (s *Schedule) Next() time.Time {
	return s.Project(1)[0]
}

// NextFrom returns the next time >= t which matches the schedule.
func (s *Schedule) NextFrom(t time.Time) time.Time {
	return s.ProjectFrom(t, 1)[0]
}

// Project returns the next N times that the expression is matched.
func (s *Schedule) Project(n int) []time.Time {
	return s.ProjectFrom(time.Now(), n)
}

// ProjectFrom returns the next N matching times after t. If t matches the expression,
// it is counted in the results.
func (s *Schedule) ProjectFrom(t time.Time, n int) []time.Time {
	last := t.In(s.location)
	results := make([]time.Time, n)
	for i := 0; i < n; i++ {
		next := s.calculateNextFromTime(last, i == 0)
		results[i] = next
		last = next
	}
	return results
}

func (s *Schedule) calculateNextFromTime(t time.Time, matchSame bool) time.Time {
	getNext := func(current time.Time) time.Time {
		// Starting at t, find the next month in which t matches the schedule
		for !s.matchesMonth(current.Month()) {
			current = current.AddDate(0, 1, 0)
			// NOTE: Add similar logic to below if we want to match specific years.
		}

		for !s.matchesDayOfMonth(current.Day()) || !s.matchesDayOfWeek(current.Weekday()) {
			current = current.AddDate(0, 0, 1)

			// Wrapped around to the 1st which increments the month, so now we need to start from the top
			// of that month again.
			if current.Day() == 1 {
				reset := time.Date(current.Year(), current.Month(), current.Day(), 0, 0, 0, 0, t.Location())
				return s.calculateNextFromTime(reset, true)
			}
		}

		for !s.matchesHour(current.Hour()) {
			current = current.Add(time.Hour)
		}

		for !s.matchesMinute(current.Minute()) {
			current = current.Add(time.Minute)
		}
		return current
	}

	// If we don't want to match the current time (maybe becasue we want to generate the next N times from now)
	// add a minute to move us along.
	if matchSame {
		return getNext(t)
	}
	return getNext(t.Add(time.Minute))
}

func (s *Schedule) matchesMonth(month time.Month) bool {
	return contains(s.expression.month.Enumerate(), int(month))
}

func (s *Schedule) matchesDayOfMonth(dayOfMonth int) bool {
	return contains(s.expression.dayOfMonth.Enumerate(), dayOfMonth)
}

func (s *Schedule) matchesDayOfWeek(dayOfWeek time.Weekday) bool {
	weekday := int(dayOfWeek)
	if weekday == 0 {
		weekday = 7
	}
	return contains(s.expression.dayOfWeek.Enumerate(), weekday)
}

func (s *Schedule) matchesHour(hour int) bool {
	return contains(s.expression.hours.Enumerate(), hour)
}

func (s *Schedule) matchesMinute(minute int) bool {
	return contains(s.expression.minutes.Enumerate(), minute)
}

// ScheduleTimer is a timer which runs on the cron schedule.
type ScheduleTimer struct {
	schedule  *Schedule
	timeChan  chan time.Time
	closeChan chan struct{}
}

// NewScheduleTimer creates a new timer.
func NewScheduleTimer(schedule *Schedule) *ScheduleTimer {
	return &ScheduleTimer{
		schedule:  schedule,
		timeChan:  make(chan time.Time),
		closeChan: make(chan struct{}),
	}
}

// Next returns a channel which upon which times will be sent when
// the schedule matches the cron expression.
// Timers must be started with Start().
func (st *ScheduleTimer) Next() <-chan time.Time {
	return st.timeChan
}

// Start starts the timer.
func (st *ScheduleTimer) Start() {
	for {
		next := st.schedule.Next()
		diff := next.Sub(time.Now().In(st.schedule.location))
		time.Sleep(diff)
		st.timeChan <- next
	}
}

// contains makes use of the fact that all expression enumerations are inherently sorted
// and uses a binary search to determine if there is a match.
func contains(haystack []int, needle int) bool {
	index := sort.Search(len(haystack), func(i int) bool { return haystack[i] >= needle })
	if index < len(haystack) && haystack[index] == needle {
		return true
	}
	return false
}
