package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestMap(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.ValCk[map[int]int]
		val       map[int]int
	}{
		{
			ID:        testhelper.MkID("MapLength - ok"),
			checkFunc: check.MapLength[map[int]int](check.ValEQ(3)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapLength - fail"),
			ExpErr: testhelper.MkExpErr(
				"the length of the map",
				" is incorrect: ",
				"must equal"),
			checkFunc: check.MapLength[map[int]int](check.ValEQ(99)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapKeyAggregate - ok"),
			checkFunc: check.MapKeyAggregate[map[int]int, int](
				check.NewCounter(check.ValGT(1), check.ValEQ(2))),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID:     testhelper.MkID("MapKeyAggregate - fail"),
			ExpErr: testhelper.MkExpErr("the value", "must equal"),
			checkFunc: check.MapKeyAggregate[map[int]int, int](
				check.NewCounter(check.ValGT(1), check.ValEQ(99))),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapValAggregate - ok"),
			checkFunc: check.MapValAggregate[map[int]int, int, int](
				check.NewCounter(check.ValGT(2), check.ValEQ(2))),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID:     testhelper.MkID("MapValAggregate - fail"),
			ExpErr: testhelper.MkExpErr("the value", "must equal"),
			checkFunc: check.MapValAggregate[map[int]int, int, int](
				check.NewCounter(check.ValGT(2), check.ValEQ(99))),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID:        testhelper.MkID("MapValAll - ok"),
			checkFunc: check.MapValAll[map[int]int](check.ValGT(0)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapValAll - fail"),
			ExpErr: testhelper.MkExpErr("bad value: the value",
				"must be greater than or equal to"),
			checkFunc: check.MapValAll[map[int]int](check.ValGE(3)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID:        testhelper.MkID("MapKeyAll - ok"),
			checkFunc: check.MapKeyAll[map[int]int](check.ValGT(0)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapKeyAll - fail"),
			ExpErr: testhelper.MkExpErr("bad key: the value",
				"must be greater than or equal to"),
			checkFunc: check.MapKeyAll[map[int]int](check.ValGE(3)),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapValAny - ok"),
			checkFunc: check.MapValAny[map[int]int](check.ValGE(3),
				"the value must be greater than or equal to 3"),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapValAny - fail"),
			ExpErr: testhelper.MkExpErr("no map values pass the test: " +
				"the value must be greater than or equal to"),
			checkFunc: check.MapValAny[map[int]int](check.ValGE(7),
				"the value must be greater than or equal to 7"),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapKeyAny - ok"),
			checkFunc: check.MapKeyAny[map[int]int](check.ValGE(3),
				"the value must be greater than or equal to 3"),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
		{
			ID: testhelper.MkID("MapKeyAny - fail"),
			ExpErr: testhelper.MkExpErr("no map keys pass the test: " +
				"the value must be greater than "),
			checkFunc: check.MapKeyAny[map[int]int](check.ValGT(3),
				"the value must be greater than 3"),
			val: map[int]int{
				1: 2,
				2: 4,
				3: 6,
			},
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}
