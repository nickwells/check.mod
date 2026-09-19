package check

import (
	"fmt"
	"slices"
	"time"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/tempus.mod/v2/tempus"
)

// TimeEQ returns a function that will check that the tested time is equal to
// the time.Time parameters
func TimeEQ(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val.Equal(t) {
			return nil
		}

		return fmt.Errorf("the time (%s) must equal %s", val, t)
	}
}

// TimeNE returns a function that will check that the tested time is not
// equal to the time.Time parameters
func TimeNE(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val != t {
			return nil
		}

		return fmt.Errorf("the time must not equal %s", t)
	}
}

// TimeGT returns a function that will check that the tested time is after
// the time.Time parameter
func TimeGT(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val.After(t) {
			return nil
		}

		return fmt.Errorf("the time (%s) must be after %s", val, t)
	}
}

// TimeGE returns a function that will check that the tested time is after
// or equal to the time.Time parameter
func TimeGE(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val.After(t) || val.Equal(t) {
			return nil
		}

		return fmt.Errorf("the time (%s) must be at or after %s", val, t)
	}
}

// TimeLT returns a function that will check that the tested time is before
// the time.Time parameter
func TimeLT(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val.Before(t) {
			return nil
		}

		return fmt.Errorf("the time (%s) must be before %s", val, t)
	}
}

// TimeLE returns a function that will check that the tested time is before
// or equal to the time.Time parameter
func TimeLE(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val.Before(t) || val.Equal(t) {
			return nil
		}

		return fmt.Errorf("the time (%s) must be at or before %s", val, t)
	}
}

// TimeBetween returns a function that will check that the tested time is
// between the start and end times (inclusive)
func TimeBetween(start, end time.Time) ValCk[time.Time] {
	if start.After(end) || start.Equal(end) {
		panic(fmt.Errorf("impossible checks passed to TimeBetween:"+
			" the start time (%v) must be before the end time (%v)",
			start, end))
	}

	return func(val time.Time) error {
		if val.Before(start) {
			return fmt.Errorf(
				"the time (%s) must be between %v and %v (too early)",
				val, start, end)
		}

		if val.After(end) {
			return fmt.Errorf(
				"the time (%s) must be between %v and %v (too late)",
				val, start, end)
		}

		return nil
	}
}

// TimeIsOnDOW returns a function that will check that the time is on the day
// of the week given by one of the parameters
func TimeIsOnDOW(dow time.Weekday, otherDOW ...time.Weekday) ValCk[time.Time] {
	days := []time.Weekday{dow}
	days = append(days, otherDOW...)

	if err := checkDays("TimeIsOnDOW", days); err != nil {
		panic(err)
	}

	return func(val time.Time) error {
		w := val.Weekday()

		if slices.Contains(days, w) {
			return nil
		}

		return fmt.Errorf("the day of the week (%s) must be a %s",
			w, daysToString(days))
	}
}

// TimeIsALeapYear checks that the time value falls on a leap year
func TimeIsALeapYear(t time.Time) error {
	if tempus.IsLeapYear(t) {
		return nil
	}

	return fmt.Errorf("the year (%d) is not a leap year", t.Year())
}

// daysFromStartOfMonth returns the number of days from the start of the
// month (0-n)
func daysFromStartOfMonth(t time.Time) int {
	return t.Day() - 1
}

// daysFromEndOfMonth returns the number of days from the end of the
// month (0-n)
func daysFromEndOfMonth(t time.Time) int {
	return tempus.DaysInMonth(t) - t.Day()
}

// TimeIsNthWeekdayOfMonth returns a function that will check that the time
// is on the nth day of the week of the month. Negative values for n mean
// that the check is from the end of the month.
func TimeIsNthWeekdayOfMonth(n int, dow time.Weekday) ValCk[time.Time] {
	if n == 0 || n > 5 || n < -5 {
		panic(fmt.Errorf(
			"impossible check passed to TimeIsNthWeekdayOfMonth:"+
				" n (== %d) must be between 1 & 5 or -5 & -1",
			n))
	}

	if err := WeekdayIsValid(dow); err != nil {
		panic(fmt.Errorf(
			"impossible check passed to TimeIsNthWeekdayOfMonth: %w", err))
	}

	return func(val time.Time) error {
		valDow := val.Weekday()
		if valDow != dow {
			return fmt.Errorf(
				"the day of the week is not %s (it is %s)",
				dow, valDow)
		}

		var valDom int

		var fromEnd bool

		if n > 0 {
			valDom = daysFromStartOfMonth(val)
		} else {
			n = -n
			valDom = daysFromEndOfMonth(val)
			fromEnd = true
		}

		wk := (valDom / tempus.DaysPerWeek) + 1
		if n != wk {
			return fmt.Errorf(
				"the day is not the %s of the month (it is the %s)",
				expectedDowDesc(n, fromEnd, dow),
				actualDowDesc(wk, fromEnd))
		}

		return nil
	}
}

// expectedDowDesc returns a description of the day of the week within a month
func expectedDowDesc(n int, fromEnd bool, dow time.Weekday) string {
	if fromEnd {
		if n == 1 {
			return "last " + dow.String()
		}

		return fmt.Sprintf("%d%s %s from the end",
			n, english.OrdinalSuffix(n), dow)
	}

	return fmt.Sprintf("%d%s %s", n, english.OrdinalSuffix(n), dow)
}

// actualDowDesc returns a description of the day of the week within a month
func actualDowDesc(n int, fromEnd bool) string {
	if fromEnd {
		if n == 1 {
			return "last"
		}

		return fmt.Sprintf("%d%s from the end", n, english.OrdinalSuffix(n))
	}

	return fmt.Sprintf("%d%s", n, english.OrdinalSuffix(n))
}
