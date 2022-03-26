package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFloat64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.Float64
		val       float64
	}{
		{
			ID:        testhelper.MkID("LT: 1.0 < 2.0"),
			checkFunc: check.Float64LT(2.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LT: 1.0 !< 1.0"),
			checkFunc: check.Float64LT(1.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LT: 2.0 !< 1.0"),
			checkFunc: check.Float64LT(1.0),
			val:       2.0,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LE: 1.0 <= 2.0"),
			checkFunc: check.Float64LE(2.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LE: 1.0 <= 1.0"),
			checkFunc: check.Float64LE(1.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LE: 2.0 !<= 1.0"),
			checkFunc: check.Float64LE(1.0),
			val:       2.0,
			ExpErr:    testhelper.MkExpErr("must be less than or equal to"),
		},
		{
			ID:        testhelper.MkID("GT: 2.0 > 1.0"),
			checkFunc: check.Float64GT(1.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("GT: 1.0 !< 1.0"),
			checkFunc: check.Float64GT(1.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GT: 2.0 !< 1.0"),
			checkFunc: check.Float64GT(2.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GE: 2.0 >= 1.0"),
			checkFunc: check.Float64GE(1.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("GE: 1.0 >= 1.0"),
			checkFunc: check.Float64GE(1.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("GE: 2.0 !<= 1.0"),
			checkFunc: check.Float64GE(2.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than or equal to"),
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 2.0 <= 3.0"),
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 1.0 <= 3.0"),
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 3.0 <= 3.0"),
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       3.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 !<= 0 <= 3.0"),
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       0.0,
			ExpErr: testhelper.MkExpErr("the value",
				"must be between",
				" - too small"),
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 4.0 !<= 3.0"),
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       4.0,
			ExpErr: testhelper.MkExpErr("the value",
				"must be between",
				" - too big"),
		},
		{
			ID: testhelper.MkID("Or: 1.0 > 2.0 , 1.0 > 3.0, 1.0 < 3.0"),
			checkFunc: check.Float64Or(
				check.Float64GT(2),
				check.Float64GT(3),
				check.Float64LT(3),
			),
			val: 1.0,
		},
		{
			ID: testhelper.MkID("Or: 7.0 > 8.0, 7.0 < 6.0, 7.0 divides 60.0"),
			checkFunc: check.Float64Or(
				check.Float64GT(8),
				check.Float64LT(6),
			),
			val: 7.0,
			ExpErr: testhelper.MkExpErr("must be greater than",
				"must be less than",
				" or "),
		},
		{
			ID: testhelper.MkID("And: 5.0 > 2.0 , 5.0 > 3.0, 5.0 < 6.0"),
			checkFunc: check.Float64And(
				check.Float64GT(2),
				check.Float64GT(3),
				check.Float64LT(6),
			),
			val: 5.0,
		},
		{
			ID: testhelper.MkID("And: 11.0 > 8.0, 11.0 < 10.0"),
			checkFunc: check.Float64And(
				check.Float64GT(8),
				check.Float64LT(10),
			),
			val:    11.0,
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 5.0 > 7.0"),
			checkFunc: check.Float64Not(
				check.Float64GT(7),
				"should not be greater than 7"),
			val: 5.0,
		},
		{
			ID: testhelper.MkID("Not: 11.0 > 7.0"),
			checkFunc: check.Float64Not(
				check.Float64GT(7),
				"should not be greater than 7"),
			val:    11.0,
			ExpErr: testhelper.MkExpErr("should not be greater than 7"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestFloat64BetweenPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		lower float64
		upper float64
	}{
		{
			ID:    testhelper.MkID("Good: Between: 1.0, 3.0"),
			lower: 1.0,
			upper: 3.0,
		},
		{
			ID:    testhelper.MkID("Bad: Between: 4.0, 3.0"),
			lower: 4.0,
			upper: 3.0,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to Float64Between: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
		{
			ID:    testhelper.MkID("Bad: Between: 2.0, 2.0"),
			lower: 2.0,
			upper: 2.0,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to Float64Between: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.Float64Between(tc.lower, tc.upper)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
