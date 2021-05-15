package check_test

import (
	"regexp"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestString(t *testing.T) {
	regexpStr := "^a[a-z]+d$"
	regexpVal := regexp.MustCompile(regexpStr)
	regexpDesc := "3 or more letters starting with an a and ending with d"
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.String
		val       string
	}{
		{
			ID:        testhelper.MkID("LenEQ: 2 == 2"),
			checkFunc: check.StringLenEQ(2),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 != 2"),
			checkFunc: check.StringLenEQ(2),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenLT: 1 < 2"),
			checkFunc: check.StringLenLT(2),
			val:       "a",
		},
		{
			ID:        testhelper.MkID("LenLT: 1 !< 1"),
			checkFunc: check.StringLenLT(1),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 2 !< 1"),
			checkFunc: check.StringLenLT(1),
			val:       "ab",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 > 1"),
			checkFunc: check.StringLenGT(1),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("LenGT: 1 !< 1"),
			checkFunc: check.StringLenGT(1),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 !< 1"),
			checkFunc: check.StringLenGT(2),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 2 <= 3"),
			checkFunc: check.StringLenBetween(1, 3),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 1 <= 3"),
			checkFunc: check.StringLenBetween(1, 3),
			val:       "a",
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 3 <= 3"),
			checkFunc: check.StringLenBetween(1, 3),
			val:       "abc",
		},
		{
			ID:        testhelper.MkID("Between: 1 !<= 0 <= 3"),
			checkFunc: check.StringLenBetween(1, 3),
			val:       "",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be between",
				" - too short"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 4 !<= 3"),
			checkFunc: check.StringLenBetween(1, 3),
			val:       "abcd",
			ExpErr: testhelper.MkExpErr(
				"the length of the value",
				"must be between",
				" - too long"),
		},
		{
			ID:        testhelper.MkID("Matches - good"),
			checkFunc: check.StringMatchesPattern(regexpVal, regexpDesc),
			val:       "abcd",
		},
		{
			ID:        testhelper.MkID("Matches - bad"),
			checkFunc: check.StringMatchesPattern(regexpVal, regexpDesc),
			val:       "xxx",
			ExpErr:    testhelper.MkExpErr("'xxx' should be:", regexpDesc),
		},
		{
			ID:        testhelper.MkID("Equals - good"),
			checkFunc: check.StringEquals("abc"),
			val:       "abc",
		},
		{
			ID:        testhelper.MkID("Equals - bad - different"),
			checkFunc: check.StringEquals("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr("should equal 'abc'"),
		},
		{
			ID:        testhelper.MkID("HasPrefix - good"),
			checkFunc: check.StringHasPrefix("abc"),
			val:       "abcx",
		},
		{
			ID:        testhelper.MkID("HasPrefix - bad - different"),
			checkFunc: check.StringHasPrefix("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr("should have 'abc' as a prefix"),
		},
		{
			ID:        testhelper.MkID("HasPrefix - bad - partly different"),
			checkFunc: check.StringHasPrefix("abc"),
			val:       "abxxxx",
			ExpErr:    testhelper.MkExpErr("should have 'abc' as a prefix"),
		},
		{
			ID:        testhelper.MkID("HasSuffix - good"),
			checkFunc: check.StringHasSuffix("abc"),
			val:       "xabc",
		},
		{
			ID:        testhelper.MkID("HasSuffix - bad - different"),
			checkFunc: check.StringHasSuffix("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr("should have 'abc' as a suffix"),
		},
		{
			ID:        testhelper.MkID("HasSuffix - bad - partly different"),
			checkFunc: check.StringHasSuffix("abc"),
			val:       "xxxxab",
			ExpErr:    testhelper.MkExpErr("should have 'abc' as a suffix"),
		},
		{
			ID: testhelper.MkID(") < 3"),
			checkFunc: check.StringOr(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(3),
			),
			val: "a",
		},
		{
			ID: testhelper.MkID(") < 1"),
			checkFunc: check.StringOr(
				check.StringLenGT(3),
				check.StringLenLT(1),
			),
			val: "ab",
			ExpErr: testhelper.MkExpErr(
				"must be greater than",
				"must be less than",
				" or "),
		},
		{
			ID: testhelper.MkID("And: len(\"abcd\") > 2 ," +
				" len(\"abcd\") > 3, len(\"abcd\") < 6"),
			checkFunc: check.StringAnd(
				check.StringLenGT(2),
				check.StringLenGT(3),
				check.StringLenLT(6),
			),
			val: "abcd",
		},
		{
			ID: testhelper.MkID(") < 3"),
			checkFunc: check.StringAnd(
				check.StringLenGT(1),
				check.StringLenLT(3),
			),
			val:    "abcd",
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID(") > 6"),
			checkFunc: check.StringNot(
				check.StringLenGT(6),
				"should not be longer than 6 characters",
			),
			val: "abcd",
		},
		{
			ID: testhelper.MkID(") > 3"),
			checkFunc: check.StringNot(
				check.StringLenGT(3),
				"should not be longer than 3 characters",
			),
			val: "abcd",
			ExpErr: testhelper.MkExpErr(
				"should not be longer than 3 characters"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestStringLenBetweenPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		lower int
		upper int
	}{
		{
			ID:    testhelper.MkID("Good: Length Between: 1, 3"),
			lower: 1,
			upper: 3,
		},
		{
			ID:    testhelper.MkID("Bad: Length Between: 4, 3"),
			lower: 4,
			upper: 3,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to StringLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			),
		},
		{
			ID:    testhelper.MkID("Bad: Length Between: 2, 2"),
			lower: 2,
			upper: 2,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to StringLenBetween: ",
				"the lower limit",
				"should be less than the upper limit",
			),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.StringLenBetween(tc.lower, tc.upper)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
