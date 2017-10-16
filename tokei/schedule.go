package tokei

import (
	"sort"
	"time"
)

type Schedule struct {
	location   *time.Location
	expression *CronExpression
}

func NewSchedule(location *time.Location, ex *CronExpression) Schedule {
	return Schedule{
		location:   location,
		expression: ex,
	}
}

func NewScheduleUTC(ex *CronExpression) Schedule {
	return NewSchedule(time.UTC, ex)
}

func (s Schedule) Next(n int) []time.Time {
	return s.NextFrom(time.Now(), n)
}

func (s Schedule) NextFrom(t time.Time, n int) []time.Time {
	last := t.In(s.location)
	results := make([]time.Time, n)
	for i := 0; i < n; i++ {
		next := s.calculateNextFromTime(last, i == 0)
		results[i] = next
		last = next
	}
	return results
}

func (s Schedule) Matches(t time.Time) bool {
	return (s.matchesMonth(t.Month()) &&
		s.matchesDayOfMonth(t.Day()) &&
		s.matchesDayOfWeek(t.Weekday()) &&
		s.matchesHour(t.Hour()) &&
		s.matchesMinute(t.Minute()))
}

func (s Schedule) calculateNextFromTime(t time.Time, matchSame bool) time.Time {

	getNext := func(current time.Time) time.Time {
		// Starting at t, find the next month in which t matches the schedule
		for !s.matchesMonth(current.Month()) {
			current = current.AddDate(0, 1, 0)

			// If we wrap around, we need to reverify all other values
			// if current.Month() == time.January {
			// 	return s.calculateNextFromTime(t)
			// }
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

func (s Schedule) matchesMonth(month time.Month) bool {
	return contains(s.expression.month.Enumerate(), int(month))
}

func (s Schedule) matchesDayOfMonth(dayOfMonth int) bool {
	return contains(s.expression.dayOfMonth.Enumerate(), dayOfMonth)
}

func (s Schedule) matchesDayOfWeek(dayOfWeek time.Weekday) bool {
	weekday := int(dayOfWeek)
	if weekday == 0 {
		weekday = 7
	}
	return contains(s.expression.dayOfWeek.Enumerate(), weekday)
}

func (s Schedule) matchesHour(hour int) bool {
	return contains(s.expression.hours.Enumerate(), hour)
}

func (s Schedule) matchesMinute(minute int) bool {
	return contains(s.expression.minutes.Enumerate(), minute)
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
