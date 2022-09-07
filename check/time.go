package check

import (
	"fmt"
	"time"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/tempus.mod/tempus"
)

// TimeEQ returns a function that will check that the tested time is equal to
// the time.Time parameters
func TimeEQ(t time.Time) ValCk[time.Time] {
	return func(val time.Time) error {
		if val == t {
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
		if val.After(t) || val == t {
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
		if val.Before(t) || val == t {
			return nil
		}

		return fmt.Errorf("the time (%s) must be at or before %s", val, t)
	}
}

// TimeBetween returns a function that will check that the tested time is
// between the start and end times (inclusive)
func TimeBetween(start, end time.Time) ValCk[time.Time] {
	if start.After(end) || start == end {
		panic(fmt.Sprintf("Impossible checks passed to TimeBetween:"+
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
	return func(val time.Time) error {
		days := []time.Weekday{dow}
		days = append(days, otherDOW...)
		valDow := val.Weekday()
		for _, dow := range days {
			if valDow == dow {
				return nil
			}
		}
		sep := ""
		validDays := ""
		for i, d := range days {
			validDays += sep + d.String()
			sep = ", "
			if i == len(days)-2 {
				sep = " or "
			}
		}
		return fmt.Errorf("the day of the week (%s) must be a %s",
			valDow, validDays)
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
		panic(fmt.Sprintf(
			"Impossible check passed to TimeIsNthWeekdayOfMonth:"+
				" n (== %d) must be between 1 & 5 or -5 & -1",
			n))
	}

	return func(val time.Time) error {
		valDow := val.Weekday()
		if valDow != dow {
			return fmt.Errorf(
				"the day of the week is not %s (it is %s)",
				dow, valDow)
		}

		var valDom int
		var desc string
		if n > 0 {
			valDom = daysFromStartOfMonth(val)
			desc = "of the month"
		} else {
			n = -n
			valDom = daysFromEndOfMonth(val)
			desc = "from the end of the month"
		}
		wk := (valDom / 7) + 1
		if n != wk {
			return fmt.Errorf(
				"the day is not the %d%s %s %s",
				n, english.OrdinalSuffix(n), dow, desc)
		}
		return nil
	}
}
