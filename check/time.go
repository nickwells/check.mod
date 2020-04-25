package check

import (
	"fmt"
	"time"
)

// Time is the type of a check function which takes an time.Time
// parameter and returns an error or nil if the check passes
type Time func(d time.Time) error

// TimeEQ returns a function that will check that the time is equal to the
// value of the t parameter
func TimeEQ(t time.Time) Time {
	return func(val time.Time) error {
		if val == t {
			return nil
		}
		return fmt.Errorf("the time (%s) should be equal to %s", val, t)
	}
}

// TimeLT returns a function that will check that the time is before the
// value of the t parameter
func TimeLT(t time.Time) Time {
	return func(val time.Time) error {
		if val.Before(t) {
			return nil
		}
		return fmt.Errorf("the time (%s) should be before %s", val, t)
	}
}

// TimeGT returns a function that will check that the time is after the
// value of the t parameter
func TimeGT(t time.Time) Time {
	return func(val time.Time) error {
		if val.After(t) {
			return nil
		}
		return fmt.Errorf("the time (%s) should be after %s", val, t)
	}
}

// TimeBetween returns a function that will check that the time is between
// the start and end times
func TimeBetween(start, end time.Time) Time {
	if start.After(end) || start == end {
		panic(fmt.Sprintf("Impossible checks passed to TimeBetween:"+
			" the start time (%s) should be before the end time (%s)",
			start, end))
	}
	return func(val time.Time) error {
		if val.Before(start) {
			return fmt.Errorf(
				"the time (%s) should be between %s and %s (too early)",
				val, start, end)
		}
		if val.After(end) {
			return fmt.Errorf(
				"the time (%s) should be between %s and %s (too late)",
				val, start, end)
		}
		return nil
	}
}

// TimeIsOnDOW returns a function that will check that the time is on the day
// of the week given by one of the parameters
func TimeIsOnDOW(dow time.Weekday, otherDOW ...time.Weekday) Time {
	return func(val time.Time) error {
		days := []time.Weekday{dow}
		days = append(days, otherDOW...)
		valDow := val.Weekday()
		for _, dow := range days {
			if valDow == dow {
				return nil
			}
		}
		sep := ""
		validDays := ""
		for i, d := range days {
			validDays += sep + d.String()
			sep = ", "
			if i == len(days)-2 {
				sep = " or "
			}
		}
		return fmt.Errorf("the day of the week (%s) should be a %s",
			valDow, validDays)
	}
}

// TimeOr returns a function that will check that the value, when passed to
// each of the check funcs in turn, passes at least one of them
func TimeOr(chkFuncs ...Time) Time {
	return func(t time.Time) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(t)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// TimeAnd returns a function that will check that the value, when passed to
// each of the check funcs in turn, passes all of them
func TimeAnd(chkFuncs ...Time) Time {
	return func(t time.Time) error {
		for _, cf := range chkFuncs {
			err := cf(t)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// TimeNot returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails. This error text should be a string
// that describes the quality that the value should not have. So, for
// instance, if the function being Not'ed was
//     check.TimeIsOnDOW(time.Monday)
// then the errMsg parameter should be
//     "on a Monday.
func TimeNot(c Time, errMsg string) Time {
	return func(v time.Time) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("'%s' should not be %s", v, errMsg)
	}
}
