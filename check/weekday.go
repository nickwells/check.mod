package check

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nickwells/english.mod/english"
)

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
		if err := WeekdayIsValid(dow); err != nil {
			return err
		}
	}

	return nil
}

// checkDays returns a non-nil error if the days slice contains either an
// invalid day or a duplicate day.
func checkDays(name string, days []time.Weekday) error {
	if err := findBadDOW(days); err != nil {
		return fmt.Errorf("impossible check passed to %s: %w", name, err)
	}

	if dupVals := findDupDOW(days); len(dupVals) > 0 {
		slices.Sort(dupVals) // sort to make tests reproducible

		return fmt.Errorf(
			"bad check passed to %s: Duplicate days-of-week: %s",
			name, strings.Join(dupVals, ", "))
	}

	return nil
}

// daysToString returns a string with the Weekdays in string form
func daysToString(days []time.Weekday) string {
	dayNames := []string{}

	for _, d := range days {
		dayNames = append(dayNames, d.String())
	}

	return english.Join(dayNames, ", ", " or ")
}

// WeekdayIsOneOf returns a func that will return a non-nil error if the
// supplied day is not one of the supplied days
func WeekdayIsOneOf(dow time.Weekday, others ...time.Weekday,
) ValCk[time.Weekday] {
	days := []time.Weekday{dow}
	days = append(days, others...)

	if err := checkDays("WeekdayIsOneOf", days); err != nil {
		panic(err)
	}

	return func(w time.Weekday) error {
		if slices.Contains(days, w) {
			return nil
		}

		return fmt.Errorf("the weekday (%s) must be a %s",
			w, daysToString(days))
	}
}

// WeekdayIsValid returns a non-nil error if the passed weekday is not
// between Sunday and Saturday inclusive, nil otherwise.
func WeekdayIsValid(w time.Weekday) error {
	if w >= time.Sunday && w <= time.Saturday {
		return nil
	}

	return fmt.Errorf(
		"the weekday (%d) is invalid, it must be in the range [%d - %d]",
		w, time.Sunday, time.Saturday)
}
