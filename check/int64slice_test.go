package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestInt64Slice(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.Int64Slice
		val       []int64
	}{
		{
			ID:        testhelper.MkID("NoDups - none - empty"),
			checkFunc: check.Int64SliceNoDups,
			val:       []int64{},
		},
		{
			ID:        testhelper.MkID("NoDups - none - 1 entry"),
			checkFunc: check.Int64SliceNoDups,
			val:       []int64{1},
		},
		{
			ID:        testhelper.MkID("NoDups - none"),
			checkFunc: check.Int64SliceNoDups,
			val:       []int64{1, 2, 3, 4, 5},
		},
		{
			ID:        testhelper.MkID("NoDups - has duplicates"),
			checkFunc: check.Int64SliceNoDups,
			val:       []int64{1, 2, 3, 3, 5},
			ExpErr: testhelper.MkExpErr("list entries:",
				"are duplicates, both are: "),
		},
		{
			ID:        testhelper.MkID("Int64Check - LT - all good"),
			checkFunc: check.Int64SliceInt64Check(check.Int64LT(10)),
			val:       []int64{1, 2, 3, 3, 5},
		},
		{
			ID:        testhelper.MkID("Int64Check - LT - all bad"),
			checkFunc: check.Int64SliceInt64Check(check.Int64GT(10)),
			val:       []int64{1, 2, 3, 3, 5},
			ExpErr: testhelper.MkExpErr(
				"list entry:",
				"does not pass the test: ",
				"greater than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 1 < 2"),
			checkFunc: check.Int64SliceLenLT(2),
			val:       []int64{1},
		},
		{
			ID:        testhelper.MkID("LenLT: 1 !< 1"),
			checkFunc: check.Int64SliceLenLT(1),
			val:       []int64{1},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenLT: 2 !< 1"),
			checkFunc: check.Int64SliceLenLT(1),
			val:       []int64{1, 2},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be less than"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 == 1"),
			checkFunc: check.Int64SliceLenEQ(1),
			val:       []int64{1},
		},
		{
			ID:        testhelper.MkID("LenEQ: 1 != 2"),
			checkFunc: check.Int64SliceLenEQ(1),
			val:       []int64{1, 2},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenEQ: 0 !=1"),
			checkFunc: check.Int64SliceLenEQ(1),
			val:       []int64{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must equal"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 > 1"),
			checkFunc: check.Int64SliceLenGT(1),
			val:       []int64{1, 2},
		},
		{
			ID:        testhelper.MkID("LenGT: 1 !< 1"),
			checkFunc: check.Int64SliceLenGT(1),
			val:       []int64{1},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenGT: 2 !< 1"),
			checkFunc: check.Int64SliceLenGT(2),
			val:       []int64{1},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be greater than"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 2 <= 3"),
			checkFunc: check.Int64SliceLenBetween(1, 3),
			val:       []int64{1, 2},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 1 <= 3"),
			checkFunc: check.Int64SliceLenBetween(1, 3),
			val:       []int64{1},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 3 <= 3"),
			checkFunc: check.Int64SliceLenBetween(1, 3),
			val:       []int64{1, 2, 3},
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 !<= 0 <= 3"),
			checkFunc: check.Int64SliceLenBetween(1, 3),
			val:       []int64{},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too short"),
		},
		{
			ID:        testhelper.MkID("LenBetween: 1 <= 4 !<= 3"),
			checkFunc: check.Int64SliceLenBetween(1, 3),
			val:       []int64{1, 2, 3, 4},
			ExpErr: testhelper.MkExpErr(
				"the length of the list",
				"must be between",
				"too long"),
		},
		{
			ID: testhelper.MkID("Or: 1 > 2 , 1 > 3, 1 < 3"),
			checkFunc: check.Int64SliceOr(
				check.Int64SliceLenGT(2),
				check.Int64SliceLenGT(3),
				check.Int64SliceLenLT(3),
			),
			val: []int64{1},
		},
		{
			ID: testhelper.MkID("Or: 7 > 8, 7 < 6, 7 divides 60"),
			checkFunc: check.Int64SliceOr(
				check.Int64SliceLenGT(8),
				check.Int64SliceLenLT(6),
				check.Int64SliceLenEQ(60),
			),
			val: []int64{1, 2, 3, 4, 5, 6, 7},
			ExpErr: testhelper.MkExpErr(
				"must be greater than",
				"must be less than",
				"must equal",
				" or "),
		},
		{
			ID: testhelper.MkID("And: 5 > 2 , 5 > 3, 5 < 6"),
			checkFunc: check.Int64SliceAnd(
				check.Int64SliceLenGT(2),
				check.Int64SliceLenGT(3),
				check.Int64SliceLenLT(6),
			),
			val: []int64{1, 2, 3, 4, 5},
		},
		{
			ID: testhelper.MkID("And: 11 > 8, 11 < 10"),
			checkFunc: check.Int64SliceAnd(
				check.Int64SliceLenGT(8),
				check.Int64SliceLenLT(10),
			),
			val:    []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 5 < 4"),
			checkFunc: check.Int64SliceNot(
				check.Int64SliceLenLT(4),
				"should not be shorter than 4 elements",
			),
			val: []int64{1, 2, 3, 4, 5},
		},
		{
			ID: testhelper.MkID("Not: 11 > 4"),
			checkFunc: check.Int64SliceNot(
				check.Int64SliceLenGT(4),
				"should not be longer than 4 elements",
			),
			val:    []int64{1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11},
			ExpErr: testhelper.MkExpErr("should not be longer than 4 elements"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}

}

func panicSafeTestInt64SliceLenBetween(t *testing.T, lowerVal, upperVal int) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.Int64SliceLenBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestInt64SliceLenBetweenPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		lower int
		upper int
	}{
		{
			ID:    testhelper.MkID("LenBetween: 1, 3"),
			lower: 1,
			upper: 3,
		},
		{
			ID:    testhelper.MkID("LenBetween: 4, 3"),
			lower: 4,
			upper: 3,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to Int64SliceLenBetween: ",
				"the lower limit",
				"should be less than the upper limit"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := panicSafeTestInt64SliceLenBetween(t,
			tc.lower, tc.upper)
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}

}
