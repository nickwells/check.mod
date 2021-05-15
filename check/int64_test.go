package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestInt64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.Int64
		val       int64
	}{
		{
			ID:        testhelper.MkID("EQ: 1 == 1"),
			checkFunc: check.Int64EQ(1),
			val:       1,
		},
		{
			ID:        testhelper.MkID("EQ: 1 != 2"),
			checkFunc: check.Int64EQ(1),
			val:       2,
			ExpErr:    testhelper.MkExpErr("must be equal to"),
		},
		{
			ID:        testhelper.MkID("LT: 1 < 2"),
			checkFunc: check.Int64LT(2),
			val:       1,
		},
		{
			ID:        testhelper.MkID("LT: 1 !< 1"),
			checkFunc: check.Int64LT(1),
			val:       1,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LT: 2 !< 1"),
			checkFunc: check.Int64LT(1),
			val:       2,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LE: 1 < 2"),
			checkFunc: check.Int64LE(2),
			val:       1,
		},
		{
			ID:        testhelper.MkID("LE: 1 <= 1"),
			checkFunc: check.Int64LE(1),
			val:       1,
		},
		{
			ID:        testhelper.MkID("LE: 2 !<= 1"),
			checkFunc: check.Int64LE(1),
			val:       2,
			ExpErr:    testhelper.MkExpErr("must be less than or equal to"),
		},
		{
			ID:        testhelper.MkID("GT: 2 > 1"),
			checkFunc: check.Int64GT(1),
			val:       2,
		},
		{
			ID:        testhelper.MkID("GT: 1 !< 1"),
			checkFunc: check.Int64GT(1),
			val:       1,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GT: 2 !< 1"),
			checkFunc: check.Int64GT(2),
			val:       1,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GE: 2 >= 1"),
			checkFunc: check.Int64GE(1),
			val:       2,
		},
		{
			ID:        testhelper.MkID("GE: 1 >= 1"),
			checkFunc: check.Int64GE(1),
			val:       1,
		},
		{
			ID:        testhelper.MkID("GE: 1 !>= 2"),
			checkFunc: check.Int64GE(2),
			val:       1,
			ExpErr:    testhelper.MkExpErr("must be greater than or equal to"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 2 <= 3"),
			checkFunc: check.Int64Between(1, 3),
			val:       2,
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 1 <= 3"),
			checkFunc: check.Int64Between(1, 3),
			val:       1,
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 3 <= 3"),
			checkFunc: check.Int64Between(1, 3),
			val:       3,
		},
		{
			ID:        testhelper.MkID("Between: 1 !<= 0 <= 3"),
			checkFunc: check.Int64Between(1, 3),
			val:       0,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be between",
				" - too small"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 4 !<= 3"),
			checkFunc: check.Int64Between(1, 3),
			val:       4,
			ExpErr: testhelper.MkExpErr(
				"the value",
				"must be between",
				" - too big"),
		},
		{
			ID:        testhelper.MkID("Divides: 2 divides 60"),
			checkFunc: check.Int64Divides(60),
			val:       2,
		},
		{
			ID:        testhelper.MkID("Divides: 7 does not divide 60"),
			checkFunc: check.Int64Divides(60),
			val:       7,
			ExpErr:    testhelper.MkExpErr("must be a divisor of"),
		},
		{
			ID:        testhelper.MkID("IsAMultiple: 20 is a multiple of 5"),
			checkFunc: check.Int64IsAMultiple(5),
			val:       20,
		},
		{
			ID:        testhelper.MkID("IsAMultiple: 21 is not a multiple of 5"),
			checkFunc: check.Int64IsAMultiple(5),
			val:       21,
			ExpErr:    testhelper.MkExpErr("must be a multiple of"),
		},
		{
			ID: testhelper.MkID("Or: 1 > 2 , 1 > 3, 1 < 3"),
			checkFunc: check.Int64Or(
				check.Int64GT(2),
				check.Int64GT(3),
				check.Int64LT(3),
			),
			val: 1,
		},
		{
			ID: testhelper.MkID("Or: 7 > 8, 7 < 6, 7 divides 60"),
			checkFunc: check.Int64Or(
				check.Int64GT(8),
				check.Int64LT(6),
				check.Int64Divides(60),
			),
			val: 7,
			ExpErr: testhelper.MkExpErr(
				"must be greater than",
				"must be less than",
				"must be a divisor of",
				" or "),
		},
		{
			ID: testhelper.MkID("And: 5 > 2 , 5 > 3, 5 < 6"),
			checkFunc: check.Int64And(
				check.Int64GT(2),
				check.Int64GT(3),
				check.Int64LT(6),
			),
			val: 5,
		},
		{
			ID: testhelper.MkID("And: 11 > 8, 11 < 10"),
			checkFunc: check.Int64And(
				check.Int64GT(8),
				check.Int64LT(10),
			),
			val: 11,
			ExpErr: testhelper.MkExpErr(
				"must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 5 > 7"),
			checkFunc: check.Int64Not(
				check.Int64GT(7),
				"should not be greater than 7"),
			val: 5,
		},
		{
			ID: testhelper.MkID("Not: 11 > 7"),
			checkFunc: check.Int64Not(
				check.Int64GT(7),
				"should not be greater than 7"),
			val:    11,
			ExpErr: testhelper.MkExpErr("should not be greater than 7"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestInt64BetweenPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		lower int64
		upper int64
	}{
		{
			ID:    testhelper.MkID("Good: Between: 1, 3"),
			lower: 1,
			upper: 3,
		},
		{
			ID:    testhelper.MkID("Bad: Between: 4, 3"),
			lower: 4,
			upper: 3,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to Int64Between: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
		{
			ID:    testhelper.MkID("Bad: Between: 2, 2"),
			lower: 2,
			upper: 2,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to Int64Between: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.Int64Between(tc.lower, tc.upper)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
