package check_test

import (
	"testing"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestWeekdayIsValid(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		val time.Weekday
	}{
		{
			ID:  testhelper.MkID("WeekdayIsValid - good, Sunday"),
			val: time.Sunday,
		},
		{
			ID:  testhelper.MkID("WeekdayIsValid - good, Saturday"),
			val: time.Saturday,
		},
		{
			ID: testhelper.MkID("WeekdayIsValid - bad, day too low"),
			ExpErr: testhelper.MkExpErr("the weekday (-1)" +
				" is invalid, it must be in the range [0 - 6]"),
			val: time.Sunday - 1,
		},
		{
			ID: testhelper.MkID("WeekdayIsValid - good, day too high"),
			ExpErr: testhelper.MkExpErr("the weekday (7)" +
				" is invalid, it must be in the range [0 - 6]"),
			val: time.Saturday + 1,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := check.WeekdayIsValid(tc.val)

			testhelper.CheckExpErr(t, err, tc)
		})
	}
}

func TestWeekdayIsOneOf(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		cf  check.ValCk[time.Weekday]
		val time.Weekday
	}{
		{
			ID:  testhelper.MkID("WeekdayIsOneOf - good, 1 day"),
			cf:  check.WeekdayIsOneOf(time.Monday),
			val: time.Monday,
		},
		{
			ID:  testhelper.MkID("WeekdayIsOneOf - good, 3 days"),
			cf:  check.WeekdayIsOneOf(time.Monday, time.Tuesday, time.Friday),
			val: time.Friday,
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - bad, 1 day"),
			ExpErr: testhelper.MkExpErr("the weekday (Tuesday)" +
				" must be a Monday"),
			cf:  check.WeekdayIsOneOf(time.Monday),
			val: time.Tuesday,
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - good, 3 days"),
			ExpErr: testhelper.MkExpErr("the weekday (Wednesday)" +
				" must be a Monday, Tuesday or Friday"),
			cf:  check.WeekdayIsOneOf(time.Monday, time.Tuesday, time.Friday),
			val: time.Wednesday,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			err := tc.cf(tc.val)

			testhelper.CheckExpErr(t, err, tc)
		})
	}
}

func TestWeekdayIsOneOfPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		dow    time.Weekday
		others []time.Weekday
	}{
		{
			ID:  testhelper.MkID("WeekdayIsOneOf - no panic, 1 day"),
			dow: time.Monday,
		},
		{
			ID:     testhelper.MkID("WeekdayIsOneOf - no panic, 2 day"),
			dow:    time.Monday,
			others: []time.Weekday{time.Tuesday},
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - panic, bad day"),
			ExpPanic: testhelper.MkExpPanic(
				"impossible check passed to WeekdayIsOneOf:" +
					" the weekday (99) is invalid," +
					" it must be in the range [0 - 6]"),
			dow:    time.Weekday(99),
			others: []time.Weekday{time.Tuesday},
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - panic, bad other day"),
			ExpPanic: testhelper.MkExpPanic(
				"impossible check passed to WeekdayIsOneOf:" +
					" the weekday (99) is invalid," +
					" it must be in the range [0 - 6]"),
			dow:    time.Monday,
			others: []time.Weekday{time.Weekday(99)},
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - panic, duplicate day"),
			ExpPanic: testhelper.MkExpPanic(
				"bad check passed to WeekdayIsOneOf:" +
					" Duplicate days-of-week: Monday appears 2 times"),
			dow:    time.Monday,
			others: []time.Weekday{time.Tuesday, time.Wednesday, time.Monday},
		},
		{
			ID: testhelper.MkID("WeekdayIsOneOf - panic, 2 dup. days"),
			ExpPanic: testhelper.MkExpPanic(
				"bad check passed to WeekdayIsOneOf:" +
					" Duplicate days-of-week:" +
					" Monday appears 2 times," +
					" Tuesday appears 3 times"),
			dow: time.Monday,
			others: []time.Weekday{
				time.Tuesday,
				time.Wednesday,
				time.Tuesday,
				time.Monday,
				time.Tuesday,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			panicked, panicVal := testhelper.PanicSafe(func() {
				check.WeekdayIsOneOf(tc.dow, tc.others...)
			})
			testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
		})
	}
}
