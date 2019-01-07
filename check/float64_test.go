package check_test

import (
	"fmt"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestFloat64(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.Float64
		val            float64
		errExpected    bool
		errMustContain []string
	}{
		{
			name:      "LT: 1.0 < 2.0",
			checkFunc: check.Float64LT(2.0),
			val:       1.0,
		},
		{
			name:           "LT: 1.0 !< 1.0",
			checkFunc:      check.Float64LT(1.0),
			val:            1.0,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:           "LT: 2.0 !< 1.0",
			checkFunc:      check.Float64LT(1.0),
			val:            2.0,
			errExpected:    true,
			errMustContain: []string{"must be less than"},
		},
		{
			name:      "LE: 1.0 <= 2.0",
			checkFunc: check.Float64LE(2.0),
			val:       1.0,
		},
		{
			name:      "LE: 1.0 <= 1.0",
			checkFunc: check.Float64LE(1.0),
			val:       1.0,
		},
		{
			name:           "LE: 2.0 !<= 1.0",
			checkFunc:      check.Float64LE(1.0),
			val:            2.0,
			errExpected:    true,
			errMustContain: []string{"must be less than or equal to"},
		},
		{
			name:      "GT: 2.0 > 1.0",
			checkFunc: check.Float64GT(1.0),
			val:       2.0,
		},
		{
			name:           "GT: 1.0 !< 1.0",
			checkFunc:      check.Float64GT(1.0),
			val:            1.0,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:           "GT: 2.0 !< 1.0",
			checkFunc:      check.Float64GT(2.0),
			val:            1.0,
			errExpected:    true,
			errMustContain: []string{"must be greater than"},
		},
		{
			name:      "GE: 2.0 >= 1.0",
			checkFunc: check.Float64GE(1.0),
			val:       2.0,
		},
		{
			name:      "GE: 1.0 >= 1.0",
			checkFunc: check.Float64GE(1.0),
			val:       1.0,
		},
		{
			name:           "GE: 2.0 !<= 1.0",
			checkFunc:      check.Float64GE(2.0),
			val:            1.0,
			errExpected:    true,
			errMustContain: []string{"must be greater than or equal to"},
		},
		{
			name:      "Between: 1.0 <= 2.0 <= 3.0",
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       2.0,
		},
		{
			name:      "Between: 1.0 <= 1.0 <= 3.0",
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       1.0,
		},
		{
			name:      "Between: 1.0 <= 3.0 <= 3.0",
			checkFunc: check.Float64Between(1.0, 3.0),
			val:       3.0,
		},
		{
			name:        "Between: 1.0 !<= 0 <= 3.0",
			checkFunc:   check.Float64Between(1.0, 3.0),
			val:         0.0,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too small",
			},
		},
		{
			name:        "Between: 1.0 <= 4.0 !<= 3.0",
			checkFunc:   check.Float64Between(1.0, 3.0),
			val:         4.0,
			errExpected: true,
			errMustContain: []string{
				"the value",
				"must be between",
				" - too big",
			},
		},
		{
			name: "Or: 1.0 > 2.0 , 1.0 > 3.0, 1.0 < 3.0",
			checkFunc: check.Float64Or(
				check.Float64GT(2),
				check.Float64GT(3),
				check.Float64LT(3),
			),
			val: 1.0,
		},
		{
			name: "Or: 7.0 > 8.0, 7.0 < 6.0, 7.0 divides 60.0",
			checkFunc: check.Float64Or(
				check.Float64GT(8),
				check.Float64LT(6),
			),
			val:         7.0,
			errExpected: true,
			errMustContain: []string{
				"must be greater than",
				"must be less than",
				" OR ",
			},
		},
		{
			name: "And: 5.0 > 2.0 , 5.0 > 3.0, 5.0 < 6.0",
			checkFunc: check.Float64And(
				check.Float64GT(2),
				check.Float64GT(3),
				check.Float64LT(6),
			),
			val: 5.0,
		},
		{
			name: "And: 11.0 > 8.0, 11.0 < 10.0",
			checkFunc: check.Float64And(
				check.Float64GT(8),
				check.Float64LT(10),
			),
			val:         11.0,
			errExpected: true,
			errMustContain: []string{
				"must be less than",
			},
		},
		{
			name: "Not: 5.0 > 7.0",
			checkFunc: check.Float64Not(
				check.Float64GT(7),
				"should not be greater than 7"),
			val: 5.0,
		},
		{
			name: "Not: 11.0 > 7.0",
			checkFunc: check.Float64Not(
				check.Float64GT(7),
				"should not be greater than 7"),
			val:         11.0,
			errExpected: true,
			errMustContain: []string{
				"should not be greater than 7",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)

		err := tc.checkFunc(tc.val)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}

}

func panicSafeTestFloat64Between(t *testing.T, lowerVal, upperVal float64) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.Float64Between(lowerVal, upperVal)
	return panicked, panicVal
}

func TestFloat64BetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            float64
		upper            float64
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:          "Between: 1.0, 3.0",
			lower:         1.0,
			upper:         3.0,
			panicExpected: false,
		},
		{
			name:          "Between: 4.0, 3.0",
			lower:         4.0,
			upper:         3.0,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to Float64Between: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		testName := fmt.Sprintf("%d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestFloat64Between(t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, testName,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
