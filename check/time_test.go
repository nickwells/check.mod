package check_test

import (
	"fmt"
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
		name           string
		cf             check.Time
		val            time.Time
		errExpected    bool
		errMustContain []string
	}{
		{
			name: "TimeEQ - good",
			cf:   check.TimeEQ(testTime),
			val:  testTime,
		},
		{
			name:           "TimeEQ - bad - before",
			cf:             check.TimeEQ(testTime),
			val:            timeMinus60s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be equal to "},
		},
		{
			name:           "TimeEQ - bad - after",
			cf:             check.TimeEQ(testTime),
			val:            timePlus60s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be equal to "},
		},
		{
			name: "TimeGT - good",
			cf:   check.TimeGT(testTime),
			val:  timePlus60s,
		},
		{
			name:           "TimeGT - bad - equal",
			cf:             check.TimeGT(testTime),
			val:            testTime,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be after "},
		},
		{
			name:           "TimeGT - bad - before",
			cf:             check.TimeGT(testTime),
			val:            timeMinus60s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be after "},
		},
		{
			name: "TimeLT - good",
			cf:   check.TimeLT(testTime),
			val:  timeMinus60s,
		},
		{
			name:           "TimeLT - bad - equal",
			cf:             check.TimeLT(testTime),
			val:            testTime,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be before "},
		},
		{
			name:           "TimeLT - bad - after",
			cf:             check.TimeLT(testTime),
			val:            timePlus60s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be before "},
		},
		{
			name: "TimeBetween - good",
			cf:   check.TimeBetween(timeMinus60s, timePlus60s),
			val:  testTime,
		},
		{
			name: "TimeBetween - good - equal start",
			cf:   check.TimeBetween(timeMinus60s, timePlus60s),
			val:  timeMinus60s,
		},
		{
			name: "TimeBetween - good - equal end",
			cf:   check.TimeBetween(timeMinus60s, timePlus60s),
			val:  timePlus60s,
		},
		{
			name:           "TimeBetween - bad - before start",
			cf:             check.TimeBetween(timeMinus60s, timePlus60s),
			val:            timeMinus120s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be between "},
		},
		{
			name:           "TimeBetween - bad - after end",
			cf:             check.TimeBetween(timeMinus60s, timePlus60s),
			val:            timePlus120s,
			errExpected:    true,
			errMustContain: []string{"the time (", ") should be between "},
		},
		{
			name: "TimeIsOnDOW - good",
			cf:   check.TimeIsOnDOW(time.Tuesday),
			val:  testTime,
		},
		{
			name:        "TimeIsOnDOW - bad",
			cf:          check.TimeIsOnDOW(time.Wednesday),
			val:         testTime,
			errExpected: true,
			errMustContain: []string{
				"the day of week (Tuesday) should be Wednesday",
			},
		},
		{
			name: "TimeOr - good - passes first test",
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val: timeMinus120s,
		},
		{
			name: "TimeOr - good - passes second test",
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val: timePlus120s,
		},
		{
			name: "TimeOr - bad",
			cf: check.TimeOr(
				check.TimeLT(timeMinus60s),
				check.TimeGT(timePlus60s)),
			val:         testTime,
			errExpected: true,
			errMustContain: []string{
				"the time (", ") should be before ",
				" OR the time (", ") should be after ",
			},
		},
		{
			name: "TimeAnd - good",
			cf: check.TimeAnd(
				check.TimeLT(timeMinus60s),
				check.TimeLT(testTime)),
			val: timeMinus120s,
		},
		{
			name: "TimeAnd - bad - fails first",
			cf: check.TimeAnd(
				check.TimeLT(testTime),
				check.TimeLT(timeMinus120s),
			),
			val:         timePlus60s,
			errExpected: true,
			errMustContain: []string{
				"the time (", ") should be before ",
			},
		},
		{
			name: "TimeAnd - bad - fails second",
			cf: check.TimeAnd(
				check.TimeLT(testTime),
				check.TimeLT(timeMinus120s),
			),
			val:         timeMinus60s,
			errExpected: true,
			errMustContain: []string{
				"the time (", ") should be before ",
			},
		},
		{
			name: "TimeNot - good",
			cf: check.TimeNot(
				check.TimeLT(timeMinus120s),
				"should not be less than the given time",
			),
			val: timeMinus60s,
		},
		{
			name: "TimeNot - bad",
			cf: check.TimeNot(
				check.TimeLT(testTime),
				"should not be before the given time",
			),
			val:         timeMinus60s,
			errExpected: true,
			errMustContain: []string{
				"should not be before the given time",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		err := tc.cf(tc.val)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}
}

func panicSafeTestTimeBetween(t *testing.T, lowerVal, upperVal time.Time) (panicked bool, panicVal interface{}) {
	t.Helper()
	defer func() {
		if r := recover(); r != nil {
			panicked = true
			panicVal = r
		}
	}()
	check.TimeBetween(lowerVal, upperVal)
	return panicked, panicVal
}

func TestTimeBetweenPanic(t *testing.T) {
	testTime := time.Date(2018, time.December, 18,
		9, 16, 30, 0, time.UTC)

	timeMinus60s := testTime.Add(-60 * time.Second)
	timePlus60s := testTime.Add(60 * time.Second)

	testCases := []struct {
		name             string
		start            time.Time
		end              time.Time
		panicExpected    bool
		panicMustContain []string
	}{
		{
			name:  "Between: good",
			start: timeMinus60s,
			end:   timePlus60s,
		},
		{
			name:          "Between: bad: start after end",
			start:         timePlus60s,
			end:           timeMinus60s,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to TimeBetween: ",
				"the start time",
				"should be before the end time",
			},
		},
		{
			name:          "Between: bad: start == end",
			start:         testTime,
			end:           testTime,
			panicExpected: true,
			panicMustContain: []string{
				"Impossible checks passed to TimeBetween: ",
				"the start time",
				"should be before the end time",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		panicked, panicVal := panicSafeTestTimeBetween(t, tc.start, tc.end)
		testhelper.PanicCheckString(t, tcID,
			panicked, tc.panicExpected,
			panicVal, tc.panicMustContain)
	}
}
