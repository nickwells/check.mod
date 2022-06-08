package check_test

import (
	"regexp"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestString(t *testing.T) {
	reStr := "^a[a-z]+d$"
	reVal := regexp.MustCompile(reStr)
	reDesc := "3 or more letters starting with an a and ending with d"
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.ValCk[string]
		val       string
	}{
		{
			ID:        testhelper.MkID("LenEQ: 2 == 2"),
			checkFunc: check.StringLength[string](check.ValEQ(2)),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 != 2"),
			checkFunc: check.StringLength[string](check.ValEQ(2)),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenLT: 1 < 2"),
			checkFunc: check.StringLength[string](check.ValLT(2)),
			val:       "a",
		},
		{
			ID:        testhelper.MkID("LenLT: 1 !< 1"),
			checkFunc: check.StringLength[string](check.ValLT(1)),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 2 !< 1"),
			checkFunc: check.StringLength[string](check.ValLT(1)),
			val:       "ab",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 > 1"),
			checkFunc: check.StringLength[string](check.ValGT(1)),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("LenGT: 1 !< 1"),
			checkFunc: check.StringLength[string](check.ValGT(1)),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 !< 1"),
			checkFunc: check.StringLength[string](check.ValGT(2)),
			val:       "a",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 2 <= 3"),
			checkFunc: check.StringLength[string](check.ValBetween(1, 3)),
			val:       "ab",
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 1 <= 3"),
			checkFunc: check.StringLength[string](check.ValBetween(1, 3)),
			val:       "a",
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 3 <= 3"),
			checkFunc: check.StringLength[string](check.ValBetween(1, 3)),
			val:       "abc",
		},
		{
			ID:        testhelper.MkID("Between: 1 !<= 0 <= 3"),
			checkFunc: check.StringLength[string](check.ValBetween(1, 3)),
			val:       "",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be between",
				" - too small"),
		},
		{
			ID:        testhelper.MkID("Between: 1 <= 4 !<= 3"),
			checkFunc: check.StringLength[string](check.ValBetween(1, 3)),
			val:       "abcd",
			ExpErr: testhelper.MkExpErr(
				"the length of the string",
				"must be between",
				" - too big"),
		},
		{
			ID:        testhelper.MkID("Matches - good"),
			checkFunc: check.StringMatchesPattern[string](reVal, reDesc),
			val:       "abcd",
		},
		{
			ID:        testhelper.MkID("Matches - bad"),
			checkFunc: check.StringMatchesPattern[string](reVal, reDesc),
			val:       "xxx",
			ExpErr:    testhelper.MkExpErr(`"xxx" should be:`, reDesc),
		},
		{
			ID:        testhelper.MkID("Equals - good"),
			checkFunc: check.ValEQ("abc"),
			val:       "abc",
		},
		{
			ID:        testhelper.MkID("Equals - bad - different"),
			checkFunc: check.ValEQ("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr(`must equal abc`),
		},
		{
			ID:        testhelper.MkID("HasPrefix - good"),
			checkFunc: check.StringHasPrefix[string]("abc"),
			val:       "abcx",
		},
		{
			ID:        testhelper.MkID("HasPrefix - bad - different"),
			checkFunc: check.StringHasPrefix[string]("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr(`should have "abc" as a prefix`),
		},
		{
			ID:        testhelper.MkID("HasPrefix - bad - partly different"),
			checkFunc: check.StringHasPrefix[string]("abc"),
			val:       "abxxxx",
			ExpErr:    testhelper.MkExpErr(`should have "abc" as a prefix`),
		},
		{
			ID:        testhelper.MkID("HasSuffix - good"),
			checkFunc: check.StringHasSuffix[string]("abc"),
			val:       "xabc",
		},
		{
			ID:        testhelper.MkID("HasSuffix - bad - different"),
			checkFunc: check.StringHasSuffix[string]("abc"),
			val:       "xxxx",
			ExpErr:    testhelper.MkExpErr(`should have "abc" as a suffix`),
		},
		{
			ID:        testhelper.MkID("HasSuffix - bad - partly different"),
			checkFunc: check.StringHasSuffix[string]("abc"),
			val:       "xxxxab",
			ExpErr:    testhelper.MkExpErr(`should have "abc" as a suffix`),
		},
		{
			ID:        testhelper.MkID("StringContains - good"),
			checkFunc: check.StringContains[string]("abc"),
			val:       "xabcy",
		},
		{
			ID:        testhelper.MkID("StringContains - good"),
			checkFunc: check.StringContains[string]("abc"),
			val:       "abcy",
		},
		{
			ID:        testhelper.MkID("StringContains - good"),
			checkFunc: check.StringContains[string]("abc"),
			val:       "xabc",
		},
		{
			ID:        testhelper.MkID("StringContains - bad"),
			checkFunc: check.StringContains[string]("abc"),
			val:       "xabyc",
			ExpErr:    testhelper.MkExpErr(`should contain "abc"`),
		},
		{
			ID:        testhelper.MkID("StringFoldedEQ - good"),
			checkFunc: check.StringFoldedEQ[string]("Hello, World!"),
			val:       "Hello, World!",
		},
		{
			ID:        testhelper.MkID("StringFoldedEQ - good"),
			checkFunc: check.StringFoldedEQ[string]("Hello, World!"),
			val:       "HELLO, WORLD!",
		},
		{
			ID:        testhelper.MkID("StringFoldedEQ - good"),
			checkFunc: check.StringFoldedEQ[string]("Hello, World!"),
			val:       "HeLlO, WoRlD!",
		},
		{
			ID:        testhelper.MkID("StringFoldedEQ - bad"),
			checkFunc: check.StringFoldedEQ[string]("Hello, World!"),
			val:       "Goodbye, World!",
			ExpErr: testhelper.MkExpErr(
				`should equal "Hello, World!" when ignoring case`),
		},
		{
			ID: testhelper.MkID(`Or: len("a") > 2, len("a") > 3, len("a") < 3`),
			checkFunc: check.Or(
				check.StringLength[string](check.ValGT(2)),
				check.StringLength[string](check.ValGT(3)),
				check.StringLength[string](check.ValLT(3)),
			),
			val: "a",
		},
		{
			ID: testhelper.MkID(`Or: len("ab") > 3, len("ab") < 1 - bad`),
			checkFunc: check.Or(
				check.StringLength[string](check.ValGT(3)),
				check.StringLength[string](check.ValLT(1)),
			),
			val: "ab",
			ExpErr: testhelper.MkExpErr(
				"must be greater than",
				"must be less than",
				" or "),
		},
		{
			ID: testhelper.MkID(`And: len(\"abcd\") > 2 ,` +
				` len("abcd") > 3, len("abcd") < 6`),
			checkFunc: check.And(
				check.StringLength[string](check.ValGT(2)),
				check.StringLength[string](check.ValGT(3)),
				check.StringLength[string](check.ValLT(6)),
			),
			val: "abcd",
		},
		{
			ID: testhelper.MkID(`And: len("abcd") > 1, len("abcd") < 3 - bad`),
			checkFunc: check.And(
				check.StringLength[string](check.ValGT(1)),
				check.StringLength[string](check.ValLT(3)),
			),
			val:    "abcd",
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID(`Not: len("abcd") > 6`),
			checkFunc: check.Not(
				check.StringLength[string](check.ValGT(6)),
				"should not be longer than 6 characters",
			),
			val: "abcd",
		},
		{
			ID: testhelper.MkID(`Not: len("abcd") > 3 - bad`),
			checkFunc: check.Not(
				check.StringLength[string](check.ValGT(3)),
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

func TestValLenBetweenPanic(t *testing.T) {
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
				"Impossible checks passed to ValBetween: ",
				"the lower limit",
				"must be less than the upper limit",
			),
		},
		{
			ID:    testhelper.MkID("Bad: Length Between: 2, 2"),
			lower: 2,
			upper: 2,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to ValBetween: ",
				"the lower limit",
				"must be less than the upper limit",
			),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.StringLength[string](check.ValBetween(tc.lower, tc.upper))
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
