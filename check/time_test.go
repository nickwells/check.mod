package check_test

import (
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestTime(t *testing.T) {
	testTime := time.Date(2018, time.December, 18,
		9, 16, 30, 0, time.UTC)

	timeMinus60s := testTime.Add(-60 * time.Second)
	timeMinus120s := testTime.Add(-120 * time.Second)

	timePlus60s := testTime.Add(60 * time.Second)
	timePlus120s := testTime.Add(120 * time.Second)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		cf  check.Time
		val time.Time
	}{
		{
			ID:  testhelper.MkID("TimeEQ - good"),
			cf:  check.TimeEQ(testTime),
			val: testTime,
		},
		{
			ID:     testhelper.MkID("TimeEQ - bad - before"),
			cf:     check.TimeEQ(testTime),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be equal to "),
		},
		{
			ID:     testhelper.MkID("TimeEQ - bad - after"),
			cf:     check.TimeEQ(testTime),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be equal to "),
		},
		{
			ID:  testhelper.MkID("TimeGT - good"),
			cf:  check.TimeGT(testTime),
			val: timePlus60s,
		},
		{
			ID:     testhelper.MkID("TimeGT - bad - equal"),
			cf:     check.TimeGT(testTime),
			val:    testTime,
			ExpErr: testhelper.MkExpErr("the time (", ") should be after "),
		},
		{
			ID:     testhelper.MkID("TimeGT - bad - before"),
			cf:     check.TimeGT(testTime),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be after "),
		},
		{
			ID:  testhelper.MkID("TimeLT - good"),
			cf:  check.TimeLT(testTime),
			val: timeMinus60s,
		},
		{
			ID:     testhelper.MkID("TimeLT - bad - equal"),
			cf:     check.TimeLT(testTime),
			val:    testTime,
			ExpErr: testhelper.MkExpErr("the time (", ") should be before "),
		},
		{
			ID:     testhelper.MkID("TimeLT - bad - after"),
			cf:     check.TimeLT(testTime),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be before "),
		},
		{
			ID:  testhelper.MkID("TimeBetween - good"),
			cf:  check.TimeBetween(timeMinus60s, timePlus60s),
			val: testTime,
		},
		{
			ID:  testhelper.MkID("TimeBetween - good - equal start"),
			cf:  check.TimeBetween(timeMinus60s, timePlus60s),
			val: timeMinus60s,
		},
		{
			ID:  testhelper.MkID("TimeBetween - good - equal end"),
			cf:  check.TimeBetween(timeMinus60s, timePlus60s),
			val: timePlus60s,
		},
		{
			ID:     testhelper.MkID("TimeBetween - bad - before start"),
			cf:     check.TimeBetween(timeMinus60s, timePlus60s),
			val:    timeMinus120s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be between "),
		},
		{
			ID:     testhelper.MkID("TimeBetween - bad - after end"),
			cf:     check.TimeBetween(timeMinus60s, timePlus60s),
			val:    timePlus120s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be between "),
		},
		{
			ID:  testhelper.MkID("TimeIsOnDOW - good"),
			cf:  check.TimeIsOnDOW(time.Tuesday),
			val: testTime,
		},
		{
			ID: testhelper.MkID("TimeIsOnDOW - multi-day - good"),
			cf: check.TimeIsOnDOW(
				time.Monday,
				time.Tuesday,
				time.Wednesday,
				time.Thursday,
				time.Friday),
			val: testTime,
		},
		{
			ID:  testhelper.MkID("TimeIsOnDOW - bad"),
			cf:  check.TimeIsOnDOW(time.Wednesday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day of the week (Tuesday) should be a Wednesday"),
		},
		{
			ID:  testhelper.MkID("TimeIsOnDOW - multi-day - bad"),
			cf:  check.TimeIsOnDOW(time.Saturday, time.Sunday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day of the week (Tuesday) should be a Saturday or Sunday"),
		},
		{
			ID: testhelper.MkID("TimeOr - good - passes first test"),
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val: timeMinus120s,
		},
		{
			ID: testhelper.MkID("TimeOr - good - passes second test"),
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val: timePlus120s,
		},
		{
			ID: testhelper.MkID("TimeOr - bad"),
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the time (", ") should be before ",
				" OR the time (", ") should be after "),
		},
		{
			ID: testhelper.MkID("TimeAnd - good"),
			cf: check.TimeAnd(
				check.TimeLT(timeMinus60s),
				check.TimeLT(testTime)),
			val: timeMinus120s,
		},
		{
			ID: testhelper.MkID("TimeAnd - bad - fails first"),
			cf: check.TimeAnd(
				check.TimeLT(testTime),
				check.TimeLT(timeMinus120s),
			),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be before "),
		},
		{
			ID: testhelper.MkID("TimeAnd - bad - fails second"),
			cf: check.TimeAnd(
				check.TimeLT(testTime),
				check.TimeLT(timeMinus120s),
			),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") should be before "),
		},
		{
			ID: testhelper.MkID("TimeNot - good"),
			cf: check.TimeNot(
				check.TimeLT(timeMinus120s),
				"should not be less than the given time",
			),
			val: timeMinus60s,
		},
		{
			ID: testhelper.MkID("TimeNot - bad"),
			cf: check.TimeNot(
				check.TimeLT(testTime),
				"should not be before the given time",
			),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("should not be before the given time"),
		},
	}

	for _, tc := range testCases {
		err := tc.cf(tc.val)
		testhelper.CheckExpErr(t, err, tc)
	}
}

func TestTimeBetweenPanic(t *testing.T) {
	testTime := time.Date(2018, time.December, 18,
		9, 16, 30, 0, time.UTC)

	timeMinus60s := testTime.Add(-60 * time.Second)
	timePlus60s := testTime.Add(60 * time.Second)

	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		start time.Time
		end   time.Time
	}{
		{
			ID:    testhelper.MkID("Good: Between: start before end"),
			start: timeMinus60s,
			end:   timePlus60s,
		},
		{
			ID:    testhelper.MkID("Bad: Between: bad: start after end"),
			start: timePlus60s,
			end:   timeMinus60s,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to TimeBetween: ",
				"the start time",
				"should be before the end time"),
		},
		{
			ID:    testhelper.MkID("Bad: Between: bad: start == end"),
			start: testTime,
			end:   testTime,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to TimeBetween: ",
				"the start time",
				"should be before the end time"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.TimeBetween(tc.start, tc.end)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
