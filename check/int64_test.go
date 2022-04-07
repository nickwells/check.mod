package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestInt64(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		checkFunc check.ValCk[int64]
		val       int64
	}{
		{
			ID:        testhelper.MkID("Divides: 2 divides 60"),
			checkFunc: check.ValDivides(int64(60)),
			val:       2,
		},
		{
			ID:        testhelper.MkID("Divides: 7 does not divide 60"),
			checkFunc: check.ValDivides(int64(60)),
			val:       7,
			ExpErr:    testhelper.MkExpErr("must be a divisor of"),
		},
		{
			ID:        testhelper.MkID("IsAMultiple: 20 is a multiple of 5"),
			checkFunc: check.ValIsAMultiple(int64(5)),
			val:       20,
		},
		{
			ID:        testhelper.MkID("IsAMultiple: 21 is not a multiple of 5"),
			checkFunc: check.ValIsAMultiple(int64(5)),
			val:       21,
			ExpErr:    testhelper.MkExpErr("must be a multiple of"),
		},
	}

	for _, tc := range testCases {
		err := tc.checkFunc(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}
