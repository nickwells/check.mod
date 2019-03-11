package check_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestDuration(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.Duration
		d              time.Duration
		errExpected    bool
		errMustContain []string
	}{
		{
			name:      "LT: 1 < 2",
			checkFunc: check.DurationLT(2 * time.Second),
			d:         1 * time.Second,
		},
		{
			name:           "LT: 1 !< 1",
			checkFunc:      check.DurationLT(1 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:           "LT: 2 !< 1",
			checkFunc:      check.DurationLT(1 * time.Second),
			d:              2 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:      "LE: 1 <= 2",
			checkFunc: check.DurationLE(2 * time.Second),
			d:         1 * time.Second,
		},
		{
			name:      "LE: 1 <= 1",
			checkFunc: check.DurationLE(1 * time.Second),
			d:         1 * time.Second,
		},
		{
			name:           "LE: 2 !<= 1",
			checkFunc:      check.DurationLE(1 * time.Second),
			d:              2 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be less than or equal to"},
		},
		{
			name:      "GT: 2 > 1",
			checkFunc: check.DurationGT(1 * time.Second),
			d:         2 * time.Second,
		},
		{
			name:           "GT: 1 !> 1",
			checkFunc:      check.DurationGT(1 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:           "GT: 1 !> 2",
			checkFunc:      check.DurationGT(2 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:      "GE: 2 >= 1",
			checkFunc: check.DurationGE(1 * time.Second),
			d:         2 * time.Second,
		},
		{
			name:      "GE: 1 >= 1",
			checkFunc: check.DurationGE(1 * time.Second),
			d:         1 * time.Second,
		},
		{
			name:           "GE: 1 !>= 2",
			checkFunc:      check.DurationGE(2 * time.Second),
			d:              1 * time.Second,
			errExpected:    true,
			errMustContain: []string{"must be greater than or equal to"},
		},
		{
			name:      "Between: 1 <= 2 <= 3",
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         2 * time.Second,
		},
		{
			name:      "Between: 1 <= 1 <= 3",
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         1 * time.Second,
		},
		{
			name:      "Between: 1 <= 3 <= 3",
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         3 * time.Second,
		},
		{
			name:        "Between: 1 !<= 0 <= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           0 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too short",
			},
		},
		{
			name:        "Between: 1 <= 4 !<= 3",
			checkFunc:   check.DurationBetween(1*time.Second, 3*time.Second),
			d:           4 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too long",
			},
		},
		{
			name: "And: 2 < 4 and 2 > 1",
			checkFunc: check.DurationAnd(
				check.DurationLT(4*time.Second),
				check.DurationGT(1*time.Second),
			),
			d: 2 * time.Second,
		},
		{
			name: "And: 2 < 4 and 1 !> 1",
			checkFunc: check.DurationAnd(
				check.DurationLT(4*time.Second),
				check.DurationGT(1*time.Second),
			),
			d:           1 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be greater than",
			},
		},
		{
			name: "Or: 2 < 4 or 2 < 1",
			checkFunc: check.DurationOr(
				check.DurationLT(1*time.Second),
				check.DurationLT(4*time.Second),
			),
			d: 2 * time.Second,
		},
		{
			name: "Or: 3 < 2 or 3 > 4",
			checkFunc: check.DurationOr(
				check.DurationLT(2*time.Second),
				check.DurationGT(4*time.Second),
			),
			d:           3 * time.Second,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be greater than",
				"must be less than",
				" or ",
			},
		},
		{
			name: "Not: 2 !< 1",
			checkFunc: check.DurationNot(
				check.DurationLT(1*time.Second),
				"should not be less than 1",
			),
			d: 2 * time.Second,
		},
		{
			name: "Not: 3 !< 4",
			checkFunc: check.DurationNot(
				check.DurationLT(4*time.Second),
				"should not be less than 4",
			),
			d:              3 * time.Second,
			errExpected:    true,
			errMustContain: []string{"should not be less than 4"},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		err := tc.checkFunc(tc.d)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}

}

func panicSafeTestDurationBetween(t *testing.T, ld, ud time.Duration) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.DurationBetween(ld, ud)
	return panicked, panicVal
}

func TestDurationBetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            time.Duration
		upper            time.Duration
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:  "Between: 1, 3",
			lower: 1 * time.Second,
			upper: 3 * time.Second,
		},
		{
			name:          "Between: 4, 3",
			lower:         4 * time.Second,
			upper:         3 * time.Second,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to DurationBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestDurationBetween(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
