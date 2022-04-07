package check_test

import (
	"regexp"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestStringSlice(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.ValCk[[]string]
		val       []string
	}{
		{
			ID:        testhelper.MkID("NoDups - none - empty"),
			checkFunc: check.SliceHasNoDups[[]string],
			val:       []string{},
		},
		{
			ID:        testhelper.MkID("NoDups - none - 1 entry"),
			checkFunc: check.SliceHasNoDups[[]string],
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("NoDups - none"),
			checkFunc: check.SliceHasNoDups[[]string],
			val:       []string{"a", "b", "c", "d", "e"},
		},
		{
			ID:        testhelper.MkID("NoDups - has duplicates"),
			checkFunc: check.SliceHasNoDups[[]string],
			val:       []string{"a", "b", "b", "c", "d"},
			ExpErr: testhelper.MkExpErr(
				"duplicate list entries:",
				"are both: b"),
		},
		{
			ID: testhelper.MkID("SliceAll - LT - all good"),
			checkFunc: check.SliceAll[[]string](
				check.StringLength[string](check.ValLT(10))),
			val: []string{"a", "b", "c", "c", "e"},
		},
		{
			ID: testhelper.MkID("SliceAll - LT - empty"),
			checkFunc: check.SliceAll[[]string](
				check.StringLength[string](check.ValLT(10))),
			val: []string{},
		},
		{
			ID: testhelper.MkID("SliceAll - LT - all bad"),
			checkFunc: check.SliceAll[[]string](
				check.StringLength[string](check.ValGT(10))),
			val: []string{"a", "b", "c", "c", "e"},
			ExpErr: testhelper.MkExpErr(
				"list entry:",
				"does not pass the test: ",
				"greater than"),
		},
		{
			ID: testhelper.MkID("SliceAny - EQ - one good"),
			checkFunc: check.SliceAny[[]string](
				check.StringLength[string](check.ValEQ(10)),
				"The length must equal 10"),
			val: []string{"a", "1234567890", "c", "c", "e"},
		},
		{
			ID:     testhelper.MkID("SliceAny - EQ - bad - empty slice"),
			ExpErr: testhelper.MkExpErr("The length must equal 10"),
			checkFunc: check.SliceAny[[]string](
				check.StringLength[string](check.ValEQ(10)),
				"The length must equal 10"),
			val: []string{},
		},
		{
			ID:     testhelper.MkID("SliceAny - EQ - bad - no valid entries"),
			ExpErr: testhelper.MkExpErr("The length must equal 10"),
			checkFunc: check.SliceAny[[]string](
				check.StringLength[string](check.ValEQ(10)),
				"The length must equal 10"),
			val: []string{"a", "b", "c"},
		},
		{
			ID:        testhelper.MkID("LenLT: 1 < 2"),
			checkFunc: check.SliceLength[[]string](check.ValLT(2)),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenLT: 1 !< 1"),
			checkFunc: check.SliceLength[[]string](check.ValLT(1)),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 2 !< 1"),
			checkFunc: check.SliceLength[[]string](check.ValLT(1)),
			val:       []string{"a", "b"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 == 1"),
			checkFunc: check.SliceLength[[]string](check.ValEQ(1)),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 != 2"),
			checkFunc: check.SliceLength[[]string](check.ValEQ(1)),
			val:       []string{"a", "b"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 0 !=1"),
			checkFunc: check.SliceLength[[]string](check.ValEQ(1)),
			val:       []string{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 > 1"),
			checkFunc: check.SliceLength[[]string](check.ValGT(1)),
			val:       []string{"a", "b"},
		},
		{
			ID:        testhelper.MkID("LenGT: 1 !< 1"),
			checkFunc: check.SliceLength[[]string](check.ValGT(1)),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 !< 1"),
			checkFunc: check.SliceLength[[]string](check.ValGT(2)),
			val:       []string{"a"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 2 <= 3"),
			checkFunc: check.SliceLength[[]string](check.ValBetween(1, 3)),
			val:       []string{"a", "b"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 1 <= 3"),
			checkFunc: check.SliceLength[[]string](check.ValBetween(1, 3)),
			val:       []string{"a"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 3 <= 3"),
			checkFunc: check.SliceLength[[]string](check.ValBetween(1, 3)),
			val:       []string{"a", "b", "c"},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 0 <= 3"),
			checkFunc: check.SliceLength[[]string](check.ValBetween(1, 3)),
			val:       []string{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too small"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 4 !<= 3"),
			checkFunc: check.SliceLength[[]string](check.ValBetween(1, 3)),
			val:       []string{"a", "b", "c", "4"},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too big"),
		},
		{
			ID: testhelper.MkID("Or: 1 > 2 , 1 > 3, 1 < 3"),
			checkFunc: check.Or(
				check.SliceLength[[]string](check.ValGT(2)),
				check.SliceLength[[]string](check.ValGT(3)),
				check.SliceLength[[]string](check.ValLT(3)),
			),
			val: []string{"a"},
		},
		{
			ID: testhelper.MkID("Or: 7 > 8, 7 < 6, 7 divides 60"),
			checkFunc: check.Or(
				check.SliceLength[[]string](check.ValGT(8)),
				check.SliceLength[[]string](check.ValLT(6)),
				check.SliceLength[[]string](check.ValEQ(60)),
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
			checkFunc: check.And(
				check.SliceLength[[]string](check.ValGT(2)),
				check.SliceLength[[]string](check.ValGT(3)),
				check.SliceLength[[]string](check.ValLT(6)),
			),
			val: []string{"a", "b", "c", "d", "e"},
		},
		{
			ID: testhelper.MkID("And: 11 > 8, 11 < 10"),
			checkFunc: check.And(
				check.SliceLength[[]string](check.ValGT(8)),
				check.SliceLength[[]string](check.ValLT(10)),
			),
			val: []string{
				"a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k",
			},
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 10 > 8"),
			checkFunc: check.Not(
				check.SliceLength[[]string](check.ValGT(8)),
				"the slice length must not be greater than 8"),
			val: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
			ExpErr: testhelper.MkExpErr(
				"the slice length must not be greater than 8"),
		},
		{
			ID: testhelper.MkID("Not: 10 > 12"),
			checkFunc: check.Not(
				check.SliceLength[[]string](check.ValGT(12)),
				"the slice length must not be greater than 12"),
			val: []string{"a", "b", "c", "d", "e", "f", "g", "h", "i", "j"},
		},
		{
			ID: testhelper.MkID("SliceByPos - all good"),
			checkFunc: check.SliceByPos[[]string](
				check.ValEQ("RC"),
				check.StringMatchesPattern[string](
					regexp.MustCompile("[1-9][0-9]*"), "a non-zero number"),
				check.ValOK[string]),
			val: []string{"RC", "1", "xxx", "yyy"},
		},
		{
			ID:        testhelper.MkID("SliceByPos - no checks - all good"),
			checkFunc: check.SliceByPos[[]string](),
			val:       []string{"RC", "1", "xxx", "yyy"},
		},
		{
			ID: testhelper.MkID("SliceByPos - bad"),
			checkFunc: check.SliceByPos[[]string](
				check.ValEQ("RC"),
				check.StringMatchesPattern[string](
					regexp.MustCompile("[1-9][0-9]*"), "a non-zero number")),
			val: []string{"XXX", "1", "2", "3"},
			ExpErr: testhelper.MkExpErr(
				"does not pass the test",
				"must equal RC"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestSliceLengthValBetweenPanic(t *testing.T) {
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
				"must be less than the upper limit"),
		},
		{
			ID:    testhelper.MkID("Bad: Length Between: 2, 2"),
			lower: 2,
			upper: 2,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to ValBetween: ",
				"the lower limit",
				"must be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.SliceLength[[]string](check.ValBetween(tc.lower, tc.upper))
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
