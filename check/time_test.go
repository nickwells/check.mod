package check_test

import (
	"testing"
	"time"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
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
		cf  check.ValCk[time.Time]
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
			ExpErr: testhelper.MkExpErr("the time (", ") must equal "),
		},
		{
			ID:     testhelper.MkID("TimeEQ - bad - after"),
			cf:     check.TimeEQ(testTime),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") must equal "),
		},
		{
			ID:     testhelper.MkID("TimeNE - bad"),
			cf:     check.TimeNE(testTime),
			val:    testTime,
			ExpErr: testhelper.MkExpErr("the time must not equal "),
		},
		{
			ID:  testhelper.MkID("TimeNE - good - before"),
			cf:  check.TimeNE(testTime),
			val: timeMinus60s,
		},
		{
			ID:  testhelper.MkID("TimeNE - good - after"),
			cf:  check.TimeNE(testTime),
			val: timePlus60s,
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
			ExpErr: testhelper.MkExpErr("the time (", ") must be after "),
		},
		{
			ID:     testhelper.MkID("TimeGT - bad - before"),
			cf:     check.TimeGT(testTime),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") must be after "),
		},
		{
			ID:  testhelper.MkID("TimeGE - good"),
			cf:  check.TimeGE(testTime),
			val: timePlus60s,
		},
		{
			ID:  testhelper.MkID("TimeGE - good - equal"),
			cf:  check.TimeGE(testTime),
			val: testTime,
		},
		{
			ID:     testhelper.MkID("TimeGE - bad - before"),
			cf:     check.TimeGE(testTime),
			val:    timeMinus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") must be at or after "),
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
			ExpErr: testhelper.MkExpErr("the time (", ") must be before "),
		},
		{
			ID:     testhelper.MkID("TimeLT - bad - after"),
			cf:     check.TimeLT(testTime),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") must be before "),
		},
		{
			ID:  testhelper.MkID("TimeLE - good"),
			cf:  check.TimeLE(testTime),
			val: timeMinus60s,
		},
		{
			ID:  testhelper.MkID("TimeLE - good - equal"),
			cf:  check.TimeLE(testTime),
			val: testTime,
		},
		{
			ID:     testhelper.MkID("TimeLE - bad - after"),
			cf:     check.TimeLE(testTime),
			val:    timePlus60s,
			ExpErr: testhelper.MkExpErr("the time (", ") must be at or before "),
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
			ExpErr: testhelper.MkExpErr("the time (", ") must be between "),
		},
		{
			ID:     testhelper.MkID("TimeBetween - bad - after end"),
			cf:     check.TimeBetween(timeMinus60s, timePlus60s),
			val:    timePlus120s,
			ExpErr: testhelper.MkExpErr("the time (", ") must be between "),
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
				"the day of the week (Tuesday) must be a Wednesday"),
		},
		{
			ID:  testhelper.MkID("TimeIsOnDOW - multi-day - bad"),
			cf:  check.TimeIsOnDOW(time.Saturday, time.Sunday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day of the week (Tuesday) must be a Saturday or Sunday"),
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
				"must be before the end time"),
		},
		{
			ID:    testhelper.MkID("Bad: Between: bad: start == end"),
			start: testTime,
			end:   testTime,
			ExpPanic: testhelper.MkExpPanic(
				"Impossible checks passed to TimeBetween: ",
				"the start time",
				"must be before the end time"),
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.TimeBetween(tc.start, tc.end)
		})
		testhelper.CheckExpPanic(t, panicked, panicVal, tc)
	}
}
