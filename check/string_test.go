package check_test

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestString(t *testing.T) {
	testCases := []struct {
		name           string
		checkFunc      check.String
		val            string
		errExpected    bool
		errMustContain []string
	}{
		{
			name:      "LenEQ: 2 == 2",
			checkFunc: check.StringLenEQ(2),
			val:       "ab",
		},
		{
			name:        "LenEQ: 1 != 2",
			checkFunc:   check.StringLenEQ(2),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must equal",
			},
		}, {
			name:      "LenLT: 1 < 2",
			checkFunc: check.StringLenLT(2),
			val:       "a",
		},
		{
			name:        "LenLT: 1 !< 1",
			checkFunc:   check.StringLenLT(1),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be less than",
			},
		},
		{
			name:        "LenLT: 2 !< 1",
			checkFunc:   check.StringLenLT(1),
			val:         "ab",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be less than",
			},
		},
		{
			name:      "LenGT: 2 > 1",
			checkFunc: check.StringLenGT(1),
			val:       "ab",
		},
		{
			name:        "LenGT: 1 !< 1",
			checkFunc:   check.StringLenGT(1),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be greater than",
			},
		},
		{
			name:        "LenGT: 2 !< 1",
			checkFunc:   check.StringLenGT(2),
			val:         "a",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be greater than",
			},
		},
		{
			name:      "Between: 1 <= 2 <= 3",
			checkFunc: check.StringLenBetween(1, 3),
			val:       "ab",
		},
		{
			name:      "Between: 1 <= 1 <= 3",
			checkFunc: check.StringLenBetween(1, 3),
			val:       "a",
		},
		{
			name:      "Between: 1 <= 3 <= 3",
			checkFunc: check.StringLenBetween(1, 3),
			val:       "abc",
		},
		{
			name:        "Between: 1 !<= 0 <= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be between",
				" - too short",
			},
		},
		{
			name:        "Between: 1 <= 4 !<= 3",
			checkFunc:   check.StringLenBetween(1, 3),
			val:         "abcd",
			errExpected: true,
			errMustContain: []string{
				"the length of the value",
				"must be between",
				" - too long",
			},
		},
		{
			name: "Matches - good",
			checkFunc: check.StringMatchesPattern(
				regexp.MustCompile("^a[a-z]+d$"),
				"3 or more letters starting with an a and ending with d"),
			val: "abcd",
		},
		{
			name: "Matches - bad",
			checkFunc: check.StringMatchesPattern(
				regexp.MustCompile("^a[a-z]+d$"),
				"3 or more letters starting with an a and ending with d"),
			val:         "xxx",
			errExpected: true,
			errMustContain: []string{
				"does not match the pattern",
			},
		},
		{
			name:      "HasPrefix - good",
			checkFunc: check.StringHasPrefix("abc"),
			val:       "abcd",
		},
		{
			name:        "HasPrefix - bad - completely different",
			checkFunc:   check.StringHasPrefix("abc"),
			val:         "xxxx",
			errExpected: true,
			errMustContain: []string{
				"should have 'abc' as a prefix",
			},
		},
		{
			name:        "HasPrefix - bad - only partly",
			checkFunc:   check.StringHasPrefix("abc"),
			val:         "abxxxx",
			errExpected: true,
			errMustContain: []string{
				"should have 'abc' as a prefix",
			},
		},
		{
			name:      "HasSuffix - good",
			checkFunc: check.StringHasSuffix("abc"),
			val:       "xabc",
		},
		{
			name:        "HasSuffix - bad - completely different",
			checkFunc:   check.StringHasSuffix("abc"),
			val:         "xxxx",
			errExpected: true,
			errMustContain: []string{
				"should have 'abc' as a suffix",
			},
		},
		{
			name:        "HasSuffix - bad - only partly",
			checkFunc:   check.StringHasSuffix("abc"),
			val:         "xxxxab",
			errExpected: true,
			errMustContain: []string{
				"should have 'abc' as a suffix",
			},
		},
		{
			name: "Or: len(\"a\") > 2 , len(\"a\") > 3, len(\"a\") < 3",
			checkFunc: check.StringOr(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(3),
			),
			val: "a",
		},
		{
			name: "Or: len(\"ab\") > 3, len(\"ab\") < 1",
			checkFunc: check.StringOr(
				check.StringLenGT(3),
				check.StringLenLT(1),
			),
			val:         "ab",
			errExpected: true,
			errMustContain: []string{
				"must be greater than",
				"must be less than",
				" OR ",
			},
		},
		{
			name: "And: len(\"abcd\") > 2 , len(\"abcd\") > 3, len(\"abcd\") < 6",
			checkFunc: check.StringAnd(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(6),
			),
			val: "abcd",
		},
		{
			name: "And: len(\"abcd\") > 1, len(\"abcd\") < 3",
			checkFunc: check.StringAnd(
				check.StringLenGT(1),
				check.StringLenLT(3),
			),
			val:         "abcd",
			errExpected: true,
			errMustContain: []string{
				"must be less than",
			},
		},
		{
			name: "Not: len(\"abcd\") > 6",
			checkFunc: check.StringNot(
				check.StringLenGT(6),
				"should not be longer than 6 characters",
			),
			val: "abcd",
		},
		{
			name: "Not: len(\"abcd\") > 3",
			checkFunc: check.StringNot(
				check.StringLenGT(3),
				"should not be longer than 3 characters",
			),
			val:         "abcd",
			errExpected: true,
			errMustContain: []string{
				"should not be longer than 3 characters",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		err := tc.checkFunc(tc.val)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}

}

func panicSafeTestStringLenBetween(t *testing.T, lowerVal, upperVal int) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.StringLenBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestStringLenBetweenPanic(t *testing.T) {
	testCases := []struct {
		name             string
		lower            int
		upper            int
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
				"Impossible checks passed to StringLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestStringLenBetween(
			t, tc.lower, tc.upper)
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}

}
