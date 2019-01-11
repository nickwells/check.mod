package check_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestInt64(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.Int64
		val            int64
		errExpected    bool
		errMustContain []string
	}{
		{
			name:      "EQ: 1 == 1",
			checkFunc: check.Int64EQ(1),
			val:       1,
		},
		{
			name:           "EQ: 1 != 2",
			checkFunc:      check.Int64EQ(1),
			val:            2,
			errExpected:    true,
			errMustContain: []string{"must be equal to"},
		},
		{
			name:      "LT: 1 < 2",
			checkFunc: check.Int64LT(2),
			val:       1,
		},
		{
			name:           "LT: 1 !< 1",
			checkFunc:      check.Int64LT(1),
			val:            1,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:           "LT: 2 !< 1",
			checkFunc:      check.Int64LT(1),
			val:            2,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:      "LE: 1 < 2",
			checkFunc: check.Int64LE(2),
			val:       1,
		},
		{
			name:      "LE: 1 <= 1",
			checkFunc: check.Int64LE(1),
			val:       1,
		},
		{
			name:           "LE: 2 !<= 1",
			checkFunc:      check.Int64LE(1),
			val:            2,
			errExpected:    true,
			errMustContain: []string{"must be less than or equal to"},
		},
		{
			name:      "GT: 2 > 1",
			checkFunc: check.Int64GT(1),
			val:       2,
		},
		{
			name:           "GT: 1 !< 1",
			checkFunc:      check.Int64GT(1),
			val:            1,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:           "GT: 2 !< 1",
			checkFunc:      check.Int64GT(2),
			val:            1,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:      "GE: 2 >= 1",
			checkFunc: check.Int64GE(1),
			val:       2,
		},
		{
			name:      "GE: 1 >= 1",
			checkFunc: check.Int64GE(1),
			val:       1,
		},
		{
			name:           "GE: 1 !>= 2",
			checkFunc:      check.Int64GE(2),
			val:            1,
			errExpected:    true,
			errMustContain: []string{"must be greater than or equal to"},
		},
		{
			name:      "Between: 1 <= 2 <= 3",
			checkFunc: check.Int64Between(1, 3),
			val:       2,
		},
		{
			name:      "Between: 1 <= 1 <= 3",
			checkFunc: check.Int64Between(1, 3),
			val:       1,
		},
		{
			name:      "Between: 1 <= 3 <= 3",
			checkFunc: check.Int64Between(1, 3),
			val:       3,
		},
		{
			name:        "Between: 1 !<= 0 <= 3",
			checkFunc:   check.Int64Between(1, 3),
			val:         0,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too small",
			},
		},
		{
			name:        "Between: 1 <= 4 !<= 3",
			checkFunc:   check.Int64Between(1, 3),
			val:         4,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too big",
			},
		},
		{
			name:      "Divides: 2 divides 60",
			checkFunc: check.Int64Divides(60),
			val:       2,
		},
		{
			name:           "Divides: 7 does not divide 60",
			checkFunc:      check.Int64Divides(60),
			val:            7,
			errExpected:    true,
			errMustContain: []string{"must be a divisor of"},
		},
		{
			name:      "IsAMultiple: 20 is a multiple of 5",
			checkFunc: check.Int64IsAMultiple(5),
			val:       20,
		},
		{
			name:           "IsAMultiple: 21 is not a multiple of 5",
			checkFunc:      check.Int64IsAMultiple(5),
			val:            21,
			errExpected:    true,
			errMustContain: []string{"must be a multiple of"},
		},
		{
			name: "Or: 1 > 2 , 1 > 3, 1 < 3",
			checkFunc: check.Int64Or(
				check.Int64GT(2),
				check.Int64GT(3),
				check.Int64LT(3),
			),
			val: 1,
		},
		{
			name: "Or: 7 > 8, 7 < 6, 7 divides 60",
			checkFunc: check.Int64Or(
				check.Int64GT(8),
				check.Int64LT(6),
				check.Int64Divides(60),
			),
			val:         7,
			errExpected: true,
			errMustContain: []string{
				"must be greater than",
				"must be less than",
				"must be a divisor of",
				" OR ",
			},
		},
		{
			name: "And: 5 > 2 , 5 > 3, 5 < 6",
			checkFunc: check.Int64And(
				check.Int64GT(2),
				check.Int64GT(3),
				check.Int64LT(6),
			),
			val: 5,
		},
		{
			name: "And: 11 > 8, 11 < 10",
			checkFunc: check.Int64And(
				check.Int64GT(8),
				check.Int64LT(10),
			),
			val:         11,
			errExpected: true,
			errMustContain: []string{
				"must be less than",
			},
		},
		{
			name: "Not: 5 > 7",
			checkFunc: check.Int64Not(
				check.Int64GT(7),
				"should not be greater than 7"),
			val: 5,
		},
		{
			name: "Not: 11 > 7",
			checkFunc: check.Int64Not(
				check.Int64GT(7),
				"should not be greater than 7"),
			val:         11,
			errExpected: true,
			errMustContain: []string{
				"should not be greater than 7",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		err := tc.checkFunc(tc.val)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}

}

func panicSafeTestInt64Between(t *testing.T, lowerVal, upperVal int64) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.Int64Between(lowerVal, upperVal)
	return panicked, panicVal
}

func TestInt64BetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            int64
		upper            int64
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:          "Between: 1, 3",
			lower:         1,
			upper:         3,
			panicExpected: false,
		},
		{
			name:          "Between: 4, 3",
			lower:         4,
			upper:         3,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to Int64Between: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestInt64Between(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}
}
