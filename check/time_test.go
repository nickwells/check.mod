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

	testTimeLeapYearRegular := time.Date(2016, time.December, 18,
		9, 16, 30, 0, time.UTC)
	testTimeLeapYearCentury := time.Date(2000, time.December, 18,
		9, 16, 30, 0, time.UTC)
	testTimeLeapYearCenturyNot := time.Date(2100, time.December, 18,
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

		{
			ID:  testhelper.MkID("TimeIsALeapYear - regular, good"),
			cf:  check.TimeIsALeapYear,
			val: testTimeLeapYearRegular,
		},
		{
			ID:  testhelper.MkID("TimeIsALeapYear - century, good"),
			cf:  check.TimeIsALeapYear,
			val: testTimeLeapYearCentury,
		},
		{
			ID:  testhelper.MkID("TimeIsALeapYear - century, bad"),
			cf:  check.TimeIsALeapYear,
			val: testTimeLeapYearCenturyNot,
			ExpErr: testhelper.MkExpErr(
				"the year (2100) is not a leap year"),
		},
		{
			ID:  testhelper.MkID("TimeIsALeapYear - regular, bad"),
			cf:  check.TimeIsALeapYear,
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the year (2018) is not a leap year"),
		},

		{
			ID:  testhelper.MkID("TimeIsNthWeekdayOfMonth - good"),
			cf:  check.TimeIsNthWeekdayOfMonth(3, time.Tuesday),
			val: testTime,
		},
		{
			ID:  testhelper.MkID("TimeIsNthWeekdayOfMonth - good (from end)"),
			cf:  check.TimeIsNthWeekdayOfMonth(-2, time.Tuesday),
			val: testTime,
		},
		{
			ID:  testhelper.MkID("TimeIsNthWeekdayOfMonth - bad - wrong DOW"),
			cf:  check.TimeIsNthWeekdayOfMonth(3, time.Monday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day of the week is not Monday (it is Tuesday)"),
		},
		{
			ID:  testhelper.MkID("TimeIsNthWeekdayOfMonth - bad - wrong n"),
			cf:  check.TimeIsNthWeekdayOfMonth(1, time.Tuesday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day is not the 1st Tuesday of the month (it is the 3rd)"),
		},
		{
			ID: testhelper.MkID(
				"TimeIsNthWeekdayOfMonth - bad - wrong n (from end)"),
			cf:  check.TimeIsNthWeekdayOfMonth(-1, time.Tuesday),
			val: testTime,
			ExpErr: testhelper.MkExpErr(
				"the day is not the 1st Tuesday" +
					" from the end of the month (it is the 2nd)"),
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
		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
	}
}

func TestTimeIsOnDOWPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		dow      time.Weekday
		otherDOW []time.Weekday
	}{
		{
			ID:  testhelper.MkID("good dow (Sunday) - no panic"),
			dow: time.Sunday,
		},
		{
			ID:  testhelper.MkID("good dow (Saturday) - no panic"),
			dow: time.Saturday,
		},
		{
			ID: testhelper.MkID("bad dow - too small"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsOnDOW:" +
					" The day-of-week (-1) is invalid" +
					" it must be in the range [0 - 6]"),
			dow: -1,
		},
		{
			ID: testhelper.MkID("bad dow - too big"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsOnDOW:" +
					" The day-of-week (7) is invalid" +
					" it must be in the range [0 - 6]"),
			dow: 7,
		},
		{
			ID: testhelper.MkID("bad otherDOW - first"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsOnDOW:" +
					" The day-of-week (7) is invalid" +
					" it must be in the range [0 - 6]"),
			dow:      time.Monday,
			otherDOW: []time.Weekday{7, time.Tuesday, time.Wednesday},
		},
		{
			ID: testhelper.MkID("bad otherDOW - last"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsOnDOW:" +
					" The day-of-week (99) is invalid" +
					" it must be in the range [0 - 6]"),
			dow:      time.Monday,
			otherDOW: []time.Weekday{time.Tuesday, time.Wednesday, 99},
		},
		{
			ID: testhelper.MkID("duplicate days"),
			ExpPanic: testhelper.MkExpPanic(
				"Bad check passed to TimeIsOnDOW:" +
					" Duplicate days-of-week:" +
					" Sunday appears 2 times"),
			dow:      time.Sunday,
			otherDOW: []time.Weekday{time.Tuesday, time.Wednesday, time.Sunday},
		},
		{
			ID: testhelper.MkID("multiple duplicate days"),
			ExpPanic: testhelper.MkExpPanic(
				"Bad check passed to TimeIsOnDOW:" +
					" Duplicate days-of-week:" +
					" Sunday appears 2 times," +
					" Tuesday appears 3 times"),
			dow: time.Sunday,
			otherDOW: []time.Weekday{
				time.Tuesday, time.Wednesday, time.Sunday,
				time.Tuesday,
				time.Tuesday,
			},
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.TimeIsOnDOW(tc.dow, tc.otherDOW...)
		})
		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
	}
}

func TestTimeIsNthWeekdayOfMonthPanic(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpPanic
		n   int
		dow time.Weekday
	}{
		{
			ID:  testhelper.MkID("good n/dow (1/Sunday) - no panic"),
			n:   1,
			dow: time.Sunday,
		},
		{
			ID:  testhelper.MkID("good n/dow (1/Saturday) - no panic"),
			n:   1,
			dow: time.Saturday,
		},
		{
			ID:  testhelper.MkID("good n/dow (5/Sunday) - no panic"),
			n:   5,
			dow: time.Sunday,
		},
		{
			ID:  testhelper.MkID("good n/dow (5/Saturday) - no panic"),
			n:   5,
			dow: time.Saturday,
		},
		{
			ID:  testhelper.MkID("good n/dow (-1/Sunday) - no panic"),
			n:   -1,
			dow: time.Sunday,
		},
		{
			ID:  testhelper.MkID("good n/dow (-1/Saturday) - no panic"),
			n:   -1,
			dow: time.Saturday,
		},
		{
			ID:  testhelper.MkID("good n/dow (-5/Sunday) - no panic"),
			n:   -5,
			dow: time.Sunday,
		},
		{
			ID:  testhelper.MkID("good n/dow (-5/Saturday) - no panic"),
			n:   -5,
			dow: time.Saturday,
		},
		{
			ID: testhelper.MkID("bad n/dow (0/Saturday) - bad n"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsNthWeekdayOfMonth:" +
					" n (== 0) must be between 1 & 5 or -5 & -1"),
			n:   0,
			dow: time.Saturday,
		},
		{
			ID: testhelper.MkID("bad n/dow (6/Saturday) - bad n"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsNthWeekdayOfMonth:" +
					" n (== 6) must be between 1 & 5 or -5 & -1"),
			n:   6,
			dow: time.Saturday,
		},
		{
			ID: testhelper.MkID("bad n/dow (-6/Saturday) - bad n"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsNthWeekdayOfMonth:" +
					" n (== -6) must be between 1 & 5 or -5 & -1"),
			n:   -6,
			dow: time.Saturday,
		},
		{
			ID: testhelper.MkID("bad n/dow (1/-1) - bad dow"),
			ExpPanic: testhelper.MkExpPanic(
				"Impossible check passed to TimeIsNthWeekdayOfMonth:" +
					" The day-of-week (-1) is invalid" +
					" it must be in the range [0 - 6]"),
			n:   1,
			dow: -1,
		},
	}

	for _, tc := range testCases {
		panicked, panicVal := testhelper.PanicSafe(func() {
			check.TimeIsNthWeekdayOfMonth(tc.n, tc.dow)
		})
		testhelper.CheckExpPanicError(t, panicked, panicVal, tc)
	}
}
