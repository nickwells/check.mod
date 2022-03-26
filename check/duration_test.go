package check_test

import (
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestDuration(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.Duration
		d         time.Duration
	}{
		{
			ID:        testhelper.MkID("LT: 1 < 2"),
			checkFunc: check.DurationLT(2 * time.Second),
			d:         1 * time.Second,
		},
		{
			ID:        testhelper.MkID("LT: 1 !< 1"),
			checkFunc: check.DurationLT(1 * time.Second),
			d:         1 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LT: 2 !< 1"),
			checkFunc: check.DurationLT(1 * time.Second),
			d:         2 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LE: 1 <= 2"),
			checkFunc: check.DurationLE(2 * time.Second),
			d:         1 * time.Second,
		},
		{
			ID:        testhelper.MkID("LE: 1 <= 1"),
			checkFunc: check.DurationLE(1 * time.Second),
			d:         1 * time.Second,
		},
		{
			ID:        testhelper.MkID("LE: 2 !<= 1"),
			checkFunc: check.DurationLE(1 * time.Second),
			d:         2 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be less than or equal to"),
		},
		{
			ID:        testhelper.MkID("GT: 2 > 1"),
			checkFunc: check.DurationGT(1 * time.Second),
			d:         2 * time.Second,
		},
		{
			ID:        testhelper.MkID("GT: 1 !> 1"),
			checkFunc: check.DurationGT(1 * time.Second),
			d:         1 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GT: 1 !> 2"),
			checkFunc: check.DurationGT(2 * time.Second),
			d:         1 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GE: 2 >= 1"),
			checkFunc: check.DurationGE(1 * time.Second),
			d:         2 * time.Second,
		},
		{
			ID:        testhelper.MkID("GE: 1 >= 1"),
			checkFunc: check.DurationGE(1 * time.Second),
			d:         1 * time.Second,
		},
		{
			ID:        testhelper.MkID("GE: 1 !>= 2"),
			checkFunc: check.DurationGE(2 * time.Second),
			d:         1 * time.Second,
			ExpErr:    testhelper.MkExpErr("must be greater than or equal to"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 2 <= 3"),
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         2 * time.Second,
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 1 <= 3"),
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         1 * time.Second,
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 3 <= 3"),
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         3 * time.Second,
		},
		{
			ID:        testhelper.MkID("Between: 1 !<= 0 <= 3"),
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         0 * time.Second,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be between",
				" - too short"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 4 !<= 3"),
			checkFunc: check.DurationBetween(1*time.Second, 3*time.Second),
			d:         4 * time.Second,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be between",
				" - too long"),
		},
		{
			ID: testhelper.MkID("And: 2 < 4 and 2 > 1"),
			checkFunc: check.DurationAnd(
				check.DurationLT(4*time.Second),
				check.DurationGT(1*time.Second),
			),
			d: 2 * time.Second,
		},
		{
			ID: testhelper.MkID("And: 2 < 4 and 1 !> 1"),
			checkFunc: check.DurationAnd(
				check.DurationLT(4*time.Second),
				check.DurationGT(1*time.Second),
			),
			d: 1 * time.Second,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be greater than"),
		},
		{
			ID: testhelper.MkID("Or: 2 < 4 or 2 < 1"),
			checkFunc: check.DurationOr(
				check.DurationLT(1*time.Second),
				check.DurationLT(4*time.Second),
			),
			d: 2 * time.Second,
		},
		{
			ID: testhelper.MkID("Or: 3 < 2 or 3 > 4"),
			checkFunc: check.DurationOr(
				check.DurationLT(2*time.Second),
				check.DurationGT(4*time.Second),
			),
			d: 3 * time.Second,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be greater than",
				"must be less than",
				" or "),
		},
		{
			ID: testhelper.MkID("Not: 2 !< 1"),
			checkFunc: check.DurationNot(check.DurationLT(1*time.Second),
				"should not be less than 1"),
			d: 2 * time.Second,
		},
		{
			ID: testhelper.MkID("Not: 3 !< 4"),
			checkFunc: check.DurationNot(check.DurationLT(4*time.Second),
				"should not be less than 4"),
			d:      3 * time.Second,
			ExpErr: testhelper.MkExpErr("should not be less than 4"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.d)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func panicSafeTestDurationBetween(t *testing.T, ld, ud time.Duration) (
	panicked bool, panicVal any,
) {
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
		testhelper.ID
		testhelper.ExpPanic
		lower time.Duration
		upper time.Duration
	}{
		{
			ID:    testhelper.MkID("Between: 1, 3"),
			lower: 1 * time.Second,
			upper: 3 * time.Second,
		},
		{
			ID:    testhelper.MkID("Between: 4, 3"),
			lower: 4 * time.Second,
			upper: 3 * time.Second,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to DurationBetween: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := panicSafeTestDurationBetween(
			t, tc.lower, tc.upper)
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
