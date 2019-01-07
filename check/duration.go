package check

import (
	"fmt"
	"time"
)

// Duration is the type of a check function which takes a time.Duration
// parameter and returns an error or nil if the check passes
type Duration func(d time.Duration) error

// DurationGT returns a function that will check that the value is
// greater than the limit
func DurationGT(limit time.Duration) Duration {
	return func(d time.Duration) error {
		if d > limit {
			return nil
		}
		return fmt.Errorf("the value (%s) must be greater than %s", d, limit)
	}
}

// DurationGE returns a function that will check that the value is
// greater than or equal to the limit
func DurationGE(limit time.Duration) Duration {
	return func(d time.Duration) error {
		if d >= limit {
			return nil
		}
		return fmt.Errorf("the value (%s) must be greater than or equal to %s",
			d, limit)
	}
}

// DurationLT returns a function that will check that the value is
// less than the limit
func DurationLT(limit time.Duration) Duration {
	return func(d time.Duration) error {
		if d < limit {
			return nil
		}
		return fmt.Errorf("the value (%s) must be less than %s", d, limit)
	}
}

// DurationLE returns a function that will check that the value is
// less than or equal to the limit
func DurationLE(limit time.Duration) Duration {
	return func(d time.Duration) error {
		if d <= limit {
			return nil
		}
		return fmt.Errorf("the value (%s) must be less than or equal to %s",
			d, limit)
	}
}

// DurationBetween  returns a function that will check that the
// value lies between the upper and lower limits (inclusive)
func DurationBetween(low, high time.Duration) Duration {
	if low >= high {
		panic(fmt.Sprintf(
			"Impossible checks passed to DurationBetween: "+
				"the lower limit (%s) should be less than the upper limit (%s)",
			low, high))
	}

	return func(d time.Duration) error {
		if d < low {
			return fmt.Errorf(
				"the value (%s) must be between %s and %s - too short",
				d, low, high)
		}
		if d > high {
			return fmt.Errorf(
				"the value (%s) must be between %s and %s - too long",
				d, low, high)
		}
		return nil
	}
}

// DurationOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func DurationOr(chkFuncs ...Duration) Duration {
	return func(d time.Duration) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(d)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// DurationAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func DurationAnd(chkFuncs ...Duration) Duration {
	return func(d time.Duration) error {
		for _, cf := range chkFuncs {
			err := cf(d)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// DurationNot returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the duration that fails
func DurationNot(c Duration, errMsg string) Duration {
	return func(d time.Duration) error {
		err := c(d)
		if err != nil {
			return nil
		}

		return fmt.Errorf("'%s' %s", d, errMsg)
	}
}
