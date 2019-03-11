package check

import "fmt"

// Float64 is the type of a check function for a float64 value. It takes a
// float64 value and returns an error or nil if the check passes
type Float64 func(f float64) error

// Float64GT returns a function that will check that the value is greater
// than the limit
func Float64GT(limit float64) Float64 {
	return func(f float64) error {
		if f > limit {
			return nil
		}
		return fmt.Errorf("the value (%f) must be greater than %f", f, limit)
	}
}

// Float64GE returns a function that will check that the value is greater
// than or equal to the limit
func Float64GE(limit float64) Float64 {
	return func(f float64) error {
		if f >= limit {
			return nil
		}
		return fmt.Errorf(
			"the value (%f) must be greater than or equal to %f", f, limit)
	}
}

// Float64LT returns a function that will check that the value is less than
// the limit
func Float64LT(limit float64) Float64 {
	return func(f float64) error {
		if f < limit {
			return nil
		}
		return fmt.Errorf("the value (%f) must be less than %f", f, limit)
	}
}

// Float64LE returns a function that will check that the value is less than
// or equal to the limit
func Float64LE(limit float64) Float64 {
	return func(f float64) error {
		if f <= limit {
			return nil
		}
		return fmt.Errorf(
			"the value (%f) must be less than or equal to %f", f, limit)
	}
}

// Float64Between returns a function that will check that the value lies
// between the upper and lower limits (inclusive)
func Float64Between(low, high float64) Float64 {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to Float64Between:"+
			" the lower limit (%f) should be less than the upper limit (%f)",
			low, high))
	}

	return func(f float64) error {
		if f < low {
			return fmt.Errorf(
				"the value (%f) must be between %f and %f - too small",
				f, low, high)
		}
		if f > high {
			return fmt.Errorf(
				"the value (%f) must be between %f and %f - too big",
				f, low, high)
		}
		return nil
	}
}

// Float64Or returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func Float64Or(chkFuncs ...Float64) Float64 {
	return func(f float64) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(f)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = _Or
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// Float64And returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func Float64And(chkFuncs ...Float64) Float64 {
	return func(f float64) error {
		for _, cf := range chkFuncs {
			err := cf(f)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// Float64Not returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails. This error text should be a string
// that describes the quality that the value should not have. So, for
// instance, if the function being Not'ed was
//     check.Float64GT(5.0)
// then the errMsg parameter should be
//     "greater than 5.0".
func Float64Not(c Float64, errMsg string) Float64 {
	return func(v float64) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("%f should not be %s", v, errMsg)
	}
}
