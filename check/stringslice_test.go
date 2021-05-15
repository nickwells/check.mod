package check_test

import (
	"regexp"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestStringSlice(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.StringSlice
		val       []string
	}{
		{
			ID:        testhelper.MkID("NoDups - none - empty"),
			checkFunc: check.StringSliceNoDups,
			val:       []string{},
		},
		{
			ID:        testhelper.MkID("NoDups - none - 1 entry"),
			checkFunc: check.StringSliceNoDups,
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("NoDups - none"),
			checkFunc: check.StringSliceNoDups,
			val:       []string{"a", "b", "c", "d", "e"},
		},
		{
			ID:        testhelper.MkID("NoDups - has duplicates"),
			checkFunc: check.StringSliceNoDups,
			val:       []string{"a", "b", "b", "c", "d"},
			ExpErr: testhelper.MkExpErr(
				"list entries:",
				"are duplicates, both are: 'b'"),
		},
		{
			ID:        testhelper.MkID("StringCheck - LT - all good"),
			checkFunc: check.StringSliceStringCheck(check.StringLenLT(10)),
			val:       []string{"a", "b", "c", "c", "e"},
		},
		{
			ID:        testhelper.MkID("StringCheck - LT - all bad"),
			checkFunc: check.StringSliceStringCheck(check.StringLenGT(10)),
			val:       []string{"a", "b", "c", "c", "e"},
			ExpErr: testhelper.MkExpErr(
				"list entry:",
				"does not pass the test: ",
				"greater than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 1 < 2"),
			checkFunc: check.StringSliceLenLT(2),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenLT: 1 !< 1"),
			checkFunc: check.StringSliceLenLT(1),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 2 !< 1"),
			checkFunc: check.StringSliceLenLT(1),
			val:       []string{"a", "b"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 == 1"),
			checkFunc: check.StringSliceLenEQ(1),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 != 2"),
			checkFunc: check.StringSliceLenEQ(1),
			val:       []string{"a", "b"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 0 !=1"),
			checkFunc: check.StringSliceLenEQ(1),
			val:       []string{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 > 1"),
			checkFunc: check.StringSliceLenGT(1),
			val:       []string{"a", "b"},
		},
		{
			ID:        testhelper.MkID("LenGT: 1 !< 1"),
			checkFunc: check.StringSliceLenGT(1),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 !< 1"),
			checkFunc: check.StringSliceLenGT(2),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 2 <= 3"),
			checkFunc: check.StringSliceLenBetween(1, 3),
			val:       []string{"a", "b"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 1 <= 3"),
			checkFunc: check.StringSliceLenBetween(1, 3),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 3 <= 3"),
			checkFunc: check.StringSliceLenBetween(1, 3),
			val:       []string{"a", "b", "c"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 !<= 0 <= 3"),
			checkFunc: check.StringSliceLenBetween(1, 3),
			val:       []string{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too short"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 4 !<= 3"),
			checkFunc: check.StringSliceLenBetween(1, 3),
			val:       []string{"a", "b", "c", "4"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too long"),
		},
		{
			ID: testhelper.MkID("Or: 1 > 2 , 1 > 3, 1 < 3"),
			checkFunc: check.StringSliceOr(
				check.StringSliceLenGT(2),
				check.StringSliceLenGT(3),
				check.StringSliceLenLT(3),
			),
			val: []string{"a"},
		},
		{
			ID: testhelper.MkID("Or: 7 > 8, 7 < 6, 7 divides 60"),
			checkFunc: check.StringSliceOr(
				check.StringSliceLenGT(8),
				check.StringSliceLenLT(6),
				check.StringSliceLenEQ(60),
			),
			val: []string{"a", "b", "c", "d", "e", "f", "g"},
			ExpErr: testhelper.MkExpErr(
				"must be greater than",
				"must be less than",
				"must equal",
				" or "),
		},
		{
			ID: testhelper.MkID("And: 5 > 2 , 5 > 3, 5 < 6"),
			checkFunc: check.StringSliceAnd(
				check.StringSliceLenGT(2),
				check.StringSliceLenGT(3),
				check.StringSliceLenLT(6),
			),
			val: []string{"a", "b", "c", "d", "e"},
		},
		{
			ID: testhelper.MkID("And: 11 > 8, 11 < 10"),
			checkFunc: check.StringSliceAnd(
				check.StringSliceLenGT(8),
				check.StringSliceLenLT(10),
			),
			val: []string{
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
			},
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 10 > 8"),
			checkFunc: check.StringSliceNot(
				check.StringSliceLenGT(8),
				"the slice length must not be greater than 8"),
			val: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			ExpErr: testhelper.MkExpErr(
				"the slice length must not be greater than 8"),
		},
		{
			ID: testhelper.MkID("Not: 10 > 12"),
			checkFunc: check.StringSliceNot(
				check.StringSliceLenGT(12),
				"the slice length must not be greater than 12"),
			val: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			ID: testhelper.MkID("StringCheckByPos - all good"),
			checkFunc: check.StringSliceStringCheckByPos(
				check.StringEquals("RC"),
				check.StringMatchesPattern(regexp.MustCompile("[1-9][0-9]*"),
					"a non-zero number"),
				check.StringOK),
			val: []string{"RC", "1", "xxx", "yyy"},
		},
		{
			ID: testhelper.MkID(
				"StringCheckByPos - no checks - all good"),
			checkFunc: check.StringSliceStringCheckByPos(),
			val:       []string{"RC", "1", "xxx", "yyy"},
		},
		{
			ID: testhelper.MkID("StringCheckByPos - bad"),
			checkFunc: check.StringSliceStringCheckByPos(
				check.StringEquals("RC"),
				check.StringMatchesPattern(regexp.MustCompile("[1-9][0-9]*"),
					"a non-zero number")),
			val: []string{"XXX", "1", "2", "3"},
			ExpErr: testhelper.MkExpErr(
				"does not pass the test",
				"should equal 'RC'"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestStringSliceLenBetweenPanic(t *testing.T) {
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
				"Impossible checks passed to StringSliceLenBetween: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
		{
			ID:    testhelper.MkID("Bad: Length Between: 2, 2"),
			lower: 2,
			upper: 2,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to StringSliceLenBetween: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.StringSliceLenBetween(tc.lower, tc.upper)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
