package check_test

import (
	"testing"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestMonthIsValid(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		val time.Month
	}{
		{
			ID:  testhelper.MkID("MonthIsValid - good, January"),
			val: time.January,
		},
		{
			ID:  testhelper.MkID("MonthIsValid - good, December"),
			val: time.December,
		},
		{
			ID: testhelper.MkID("MonthIsValid - bad, month too small"),
			ExpErr: testhelper.MkExpErr("the month (0)" +
				" is invalid, it must be in the range [1 - 12]"),
			val: time.Month(0),
		},
		{
			ID: testhelper.MkID("MonthIsValid - bad, month too big"),
			ExpErr: testhelper.MkExpErr("the month (13)" +
				" is invalid, it must be in the range [1 - 12]"),
			val: time.Month(13),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := check.MonthIsValid(tc.val)

			testhelper.CheckExpErr(t, err, tc)
		})
	}
}

func TestMonthIsOneOf(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		cf  check.ValCk[time.Month]
		val time.Month
	}{
		{
			ID:  testhelper.MkID("MonthIsOneOf - good, 1 month"),
			cf:  check.MonthIsOneOf(time.January),
			val: time.January,
		},
		{
			ID:  testhelper.MkID("MonthIsOneOf - good, 3 months"),
			cf:  check.MonthIsOneOf(time.January, time.February, time.May),
			val: time.May,
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - bad, 1 month"),
			ExpErr: testhelper.MkExpErr("the month (February)" +
				" must be a January"),
			cf:  check.MonthIsOneOf(time.January),
			val: time.February,
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - good, 3 months"),
			ExpErr: testhelper.MkExpErr("the month (March)" +
				" must be a January, February or May"),
			cf:  check.MonthIsOneOf(time.January, time.February, time.May),
			val: time.March,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.cf(tc.val)

			testhelper.CheckExpErr(t, err, tc)
		})
	}
}

func TestMonthIsOneOfPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		dow    time.Month
		others []time.Month
	}{
		{
			ID:  testhelper.MkID("MonthIsOneOf - no panic, 1 month"),
			dow: time.January,
		},
		{
			ID:     testhelper.MkID("MonthIsOneOf - no panic, 2 month"),
			dow:    time.January,
			others: []time.Month{time.February},
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - panic, bad month"),
			ExpPanic: testhelper.MkExpPanic(
				"impossible check passed to MonthIsOneOf:" +
					" the month (99) is invalid," +
					" it must be in the range [1 - 12]"),
			dow:    time.Month(99),
			others: []time.Month{time.February},
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - panic, bad other month"),
			ExpPanic: testhelper.MkExpPanic(
				"impossible check passed to MonthIsOneOf:" +
					" the month (99) is invalid," +
					" it must be in the range [1 - 12]"),
			dow:    time.January,
			others: []time.Month{time.Month(99)},
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - panic, duplicate month"),
			ExpPanic: testhelper.MkExpPanic(
				"bad check passed to MonthIsOneOf:" +
					" Duplicate months: January appears 2 times"),
			dow:    time.January,
			others: []time.Month{time.February, time.March, time.January},
		},
		{
			ID: testhelper.MkID("MonthIsOneOf - panic, 2 dup. months"),
			ExpPanic: testhelper.MkExpPanic(
				"bad check passed to MonthIsOneOf:" +
					" Duplicate months:" +
					" January appears 2 times," +
					" May appears 3 times"),
			dow: time.January,
			others: []time.Month{
				time.May,
				time.March,
				time.May,
				time.January,
				time.May,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			panicked, panicVal := testhelper.PanicSafe(func() {
				check.MonthIsOneOf(tc.dow, tc.others...)
			})
			testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
		})
	}
}
