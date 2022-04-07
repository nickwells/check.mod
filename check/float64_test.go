package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFloat64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.ValCk[float64]
		val       float64
	}{
		{
			ID:        testhelper.MkID("EQ: 2.0 == 2.0"),
			checkFunc: check.ValEQ(2.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("EQ: 2.0 != 1.0"),
			checkFunc: check.ValEQ(1.0),
			val:       2.0,
			ExpErr:    testhelper.MkExpErr("must equal"),
		},
		{
			ID:        testhelper.MkID("NE: 1.0 == 2.0"),
			checkFunc: check.ValNE(2.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("NE: 1.0 != 1.0"),
			checkFunc: check.ValNE(1.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must not equal"),
		},
		{
			ID:        testhelper.MkID("LT: 1.0 < 2.0"),
			checkFunc: check.ValLT(2.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LT: 1.0 !< 1.0"),
			checkFunc: check.ValLT(1.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LT: 2.0 !< 1.0"),
			checkFunc: check.ValLT(1.0),
			val:       2.0,
			ExpErr:    testhelper.MkExpErr("must be less than"),
		},
		{
			ID:        testhelper.MkID("LE: 1.0 <= 2.0"),
			checkFunc: check.ValLE(2.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LE: 1.0 <= 1.0"),
			checkFunc: check.ValLE(1.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("LE: 2.0 !<= 1.0"),
			checkFunc: check.ValLE(1.0),
			val:       2.0,
			ExpErr:    testhelper.MkExpErr("must be less than or equal to"),
		},
		{
			ID:        testhelper.MkID("GT: 2.0 > 1.0"),
			checkFunc: check.ValGT(1.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("GT: 1.0 !< 1.0"),
			checkFunc: check.ValGT(1.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GT: 2.0 !< 1.0"),
			checkFunc: check.ValGT(2.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than"),
		},
		{
			ID:        testhelper.MkID("GE: 2.0 >= 1.0"),
			checkFunc: check.ValGE(1.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("GE: 1.0 >= 1.0"),
			checkFunc: check.ValGE(1.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("GE: 2.0 !<= 1.0"),
			checkFunc: check.ValGE(2.0),
			val:       1.0,
			ExpErr:    testhelper.MkExpErr("must be greater than or equal to"),
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 2.0 <= 3.0"),
			checkFunc: check.ValBetween(1.0, 3.0),
			val:       2.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 1.0 <= 3.0"),
			checkFunc: check.ValBetween(1.0, 3.0),
			val:       1.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 3.0 <= 3.0"),
			checkFunc: check.ValBetween(1.0, 3.0),
			val:       3.0,
		},
		{
			ID:        testhelper.MkID("Between: 1.0 !<= 0 <= 3.0"),
			checkFunc: check.ValBetween(1.0, 3.0),
			val:       0.0,
			ExpErr: testhelper.MkExpErr("the value",
				"must be between",
				" - too small"),
		},
		{
			ID:        testhelper.MkID("Between: 1.0 <= 4.0 !<= 3.0"),
			checkFunc: check.ValBetween(1.0, 3.0),
			val:       4.0,
			ExpErr: testhelper.MkExpErr("the value",
				"must be between",
				" - too big"),
		},
		{
			ID: testhelper.MkID("Or: 1.0 > 2.0 , 1.0 > 3.0, 1.0 < 3.0"),
			checkFunc: check.Or(
				check.ValGT[float64](2),
				check.ValGT[float64](3),
				check.ValLT[float64](3),
			),
			val: 1.0,
		},
		{
			ID: testhelper.MkID("Or: 7.0 > 8.0, 7.0 < 6.0, 7.0 divides 60.0"),
			checkFunc: check.Or(
				check.ValGT[float64](8),
				check.ValLT[float64](6),
			),
			val: 7.0,
			ExpErr: testhelper.MkExpErr("must be greater than",
				"must be less than",
				" or "),
		},
		{
			ID: testhelper.MkID("And: 5.0 > 2.0 , 5.0 > 3.0, 5.0 < 6.0"),
			checkFunc: check.And(
				check.ValGT[float64](2),
				check.ValGT[float64](3),
				check.ValLT[float64](6),
			),
			val: 5.0,
		},
		{
			ID: testhelper.MkID("And: 11.0 > 8.0, 11.0 < 10.0"),
			checkFunc: check.And(
				check.ValGT[float64](8),
				check.ValLT[float64](10),
			),
			val:    11.0,
			ExpErr: testhelper.MkExpErr("must be less than"),
		},
		{
			ID: testhelper.MkID("Not: 5.0 > 7.0"),
			checkFunc: check.Not(
				check.ValGT[float64](7),
				"should not be greater than 7"),
			val: 5.0,
		},
		{
			ID: testhelper.MkID("Not: 11.0 > 7.0"),
			checkFunc: check.Not(
				check.ValGT[float64](7),
				"should not be greater than 7"),
			val:    11.0,
			ExpErr: testhelper.MkExpErr("should not be greater than 7"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}
