package check

import (
	"fmt"
	"slices"
	"strings"
	"time"

	"github.com/nickwells/english.mod/english"
)

// findDupMonth will return a slice (possibly empty) describing all the
// mopnths that appear multiple times in the supplied slice.
func findDupMonth(months []time.Month) []string {
	dupChk := map[time.Month]int{}

	for _, month := range months {
		dupChk[month]++
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

// findBadMonth returns an error for the first bad entry in the slice of
// Months or nil if they are all good
func findBadMonth(months []time.Month) error {
	for _, month := range months {
		if err := MonthIsValid(month); err != nil {
			return err
		}
	}

	return nil
}

// checkMonths returns a non-nil error if the months slice contains either an
// invalid month or a duplicate month.
func checkMonths(name string, months []time.Month) error {
	if err := findBadMonth(months); err != nil {
		return fmt.Errorf("impossible check passed to %s: %w", name, err)
	}

	if dupVals := findDupMonth(months); len(dupVals) > 0 {
		slices.Sort(dupVals) // sort to make tests reproducible

		return fmt.Errorf(
			"bad check passed to %s: Duplicate months: %s",
			name, strings.Join(dupVals, ", "))
	}

	return nil
}

// monthsToString returns a string with the Months in string form
func monthsToString(months []time.Month) string {
	monthNames := []string{}

	for _, m := range months {
		monthNames = append(monthNames, m.String())
	}

	return english.Join(monthNames, ", ", " or ")
}

// MonthIsOneOf returns a func that will return a non-nil error if the
// supplied month is not one of the supplied months
func MonthIsOneOf(month time.Month, others ...time.Month) ValCk[time.Month] {
	months := []time.Month{month}
	months = append(months, others...)

	if err := checkMonths("MonthIsOneOf", months); err != nil {
		panic(err)
	}

	return func(m time.Month) error {
		if slices.Contains(months, m) {
			return nil
		}

		return fmt.Errorf("the month (%s) must be a %s",
			m, monthsToString(months))
	}
}

// MonthIsValid returns a non-nil error if the passed month is not between
// January and December inclusive, nil otherwise.
func MonthIsValid(month time.Month) error {
	if month >= time.January && month <= time.December {
		return nil
	}

	return fmt.Errorf(
		"the month (%d) is invalid, it must be in the range [%d - %d]",
		month, time.January, time.December)
}
