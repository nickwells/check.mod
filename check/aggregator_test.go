package check_test

import (
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestSliceAggregateCounter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		counter *check.Counter[bool]
		slc     []bool
	}{
		{
			ID:      testhelper.MkID("ok"),
			counter: check.NewCounter(check.ValEQ(true), check.ValEQ(5)),
			slc:     []bool{true, true, false, true, false, true, false, true},
		},
		{
			ID:      testhelper.MkID("fail"),
			ExpErr:  testhelper.MkExpErr("the value", "must equal"),
			counter: check.NewCounter(check.ValEQ(true), check.ValEQ(5)),
			slc:     []bool{true, true, false, true, false, true, false},
		},
	}

	for _, tc := range testCases {
		vc := check.SliceAggregate[[]bool, bool](tc.counter)
		err := vc(tc.slc)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestMapAggregateCounter(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		counter *check.Counter[bool]
		m       map[string]bool
	}{
		{
			ID:      testhelper.MkID("ok"),
			counter: check.NewCounter(check.ValEQ(true), check.ValEQ(5)),
			m: map[string]bool{
				"one":   true,
				"two":   true,
				"three": true,
				"four":  true,
				"five":  false,
				"six":   true,
			},
		},
		{
			ID:      testhelper.MkID("bad"),
			ExpErr:  testhelper.MkExpErr("the value", "must equal"),
			counter: check.NewCounter(check.ValEQ(true), check.ValEQ(5)),
			m: map[string]bool{
				"one":   true,
				"two":   true,
				"three": true,
				"four":  true,
				"five":  false,
			},
		},
	}

	for _, tc := range testCases {
		vc := check.MapValAggregate[map[string]bool, string, bool](tc.counter)
		err := vc(tc.m)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestNewCounterPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		vtf check.ValCk[string]
		ctf check.ValCk[int]
	}{
		{
			ID:  testhelper.MkID("ok"),
			vtf: check.ValEQ("xxx"),
			ctf: check.ValEQ(99),
		},
		{
			ID: testhelper.MkID("bad - no value test function"),
			ExpPanic: testhelper.MkExpPanic(
				"no value test function has been given"),
			vtf: nil,
			ctf: check.ValEQ(99),
		},
		{
			ID: testhelper.MkID("bad - no count test function"),
			ExpPanic: testhelper.MkExpPanic(
				"no count test function has been given"),
			vtf: check.ValEQ("xxx"),
			ctf: nil,
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.NewCounter(tc.vtf, tc.ctf)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
