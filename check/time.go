package check

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nickwells/english.mod/english"
	"github.com/nickwells/tempus.mod/tempus"
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

// dowValid returns an error if the passed weekday is not between Sunday and
// Saturday inclusive, nil otherwise.
func dowValid(dow time.Weekday) error {
	if dow >= time.Sunday && dow <= time.Saturday {
		return nil
	}

	return fmt.Errorf(
		"the day-of-week (%d) is invalid it must be in the range [%d - %d]",
		dow, time.Sunday, time.Saturday)
}

// findDupDOW will return a slice (possibly empty) describing all the days of
// the week that appear multiple times in the supplied slice.
func findDupDOW(dows []time.Weekday) []string {
	dupChk := map[time.Weekday]int{}

	for _, dow := range dows {
		dupChk[dow]++
	}

	dupVals := []string{}

	for k, count := range dupChk {
		if count > 1 {
			dupVals = append(dupVals,
				fmt.Sprintf("%s appears %d times", k, count))
		}
	}

	return dupVals
}

// findBadDOW returns an error for the first bad entry in the slice of
// Weekdays or nil if they are all good
func findBadDOW(dows []time.Weekday) error {
	for _, dow := range dows {
		if err := dowValid(dow); err != nil {
			return err
		}
	}

	return nil
}

// TimeIsOnDOW returns a function that will check that the time is on the day
// of the week given by one of the parameters
func TimeIsOnDOW(dow time.Weekday, otherDOW ...time.Weekday) ValCk[time.Time] {
	if err := findBadDOW(append(otherDOW, dow)); err != nil {
		panic(fmt.Errorf("impossible check passed to TimeIsOnDOW: %w", err))
	}

	if dupVals := findDupDOW(append(otherDOW, dow)); len(dupVals) > 0 {
		slices.Sort(dupVals) // sort to make tests reproducible
		panic(fmt.Errorf(
			"bad check passed to TimeIsOnDOW: Duplicate days-of-week: %s",
			strings.Join(dupVals, ", ")))
	}

	days := []time.Weekday{dow}
	days = append(days, otherDOW...)

	return func(val time.Time) error {
		valDow := val.Weekday()

		if slices.Contains(days, valDow) {
			return nil
		}

		dayNames := []string{}

		for _, d := range days {
			dayNames = append(dayNames, d.String())
		}

		return fmt.Errorf("the day of the week (%s) must be a %s",
			valDow, english.Join(dayNames, ", ", " or "))
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

	if err := dowValid(dow); err != nil {
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
