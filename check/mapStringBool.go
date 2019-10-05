package check

import "fmt"

// MapStringBool is the type of a check function for a slice of strings. It
// takes a slice of strings as a parameter and returns an error or nil if
// there is no error
type MapStringBool func(v map[string]bool) error

// MapStringBoolStringCheck returns a check function that checks that every
// key in the map passes the supplied String check func
func MapStringBoolStringCheck(sc String) MapStringBool {
	return func(v map[string]bool) error {
		for k := range v {
			if err := sc(k); err != nil {
				return fmt.Errorf(
					"map entry: %q - the key does not pass the test: %s",
					k, err)
			}
		}
		return nil
	}
}

// MapStringBoolContains returns a check function that checks that at least
// one entry in the list matches the supplied String check func. The
// condition parameter should describe the check that is being performed. For
// instance, if the check is that the string length must be greater than 5
// characters then the condition parameter should be "the string should be
// greater than 5 characters"
func MapStringBoolContains(sc String, condition string) MapStringBool {
	return func(v map[string]bool) error {
		for k := range v {
			if err := sc(k); err == nil {
				return nil
			}
		}
		return fmt.Errorf("none of the list entries passes the test: %s",
			condition)
	}
}

// MapStringBoolLenEQ returns a check function that checks that the length of
// the list equals the supplied value
func MapStringBoolLenEQ(i int) MapStringBool {
	return func(v map[string]bool) error {
		if len(v) != i {
			return fmt.Errorf("the number of entries (%d)"+" must equal %d",
				len(v), i)
		}
		return nil
	}
}

// MapStringBoolLenGT returns a check function that checks that the length of
// the list is greater than the supplied value
func MapStringBoolLenGT(i int) MapStringBool {
	return func(v map[string]bool) error {
		if len(v) <= i {
			return fmt.Errorf(
				"the number of entries (%d) must be greater than %d",
				len(v), i)
		}
		return nil
	}
}

// MapStringBoolLenLT returns a check function that checks that the length of
// the list is less than the supplied value
func MapStringBoolLenLT(i int) MapStringBool {
	return func(v map[string]bool) error {
		if len(v) >= i {
			return fmt.Errorf("the number of entries (%d) must be less than %d",
				len(v), i)
		}
		return nil
	}
}

// MapStringBoolLenBetween returns a check function that checks that the length
// of the list is between the two supplied values (inclusive)
func MapStringBoolLenBetween(low, high int) MapStringBool {
	if low >= high {
		panic(fmt.Sprintf(
			"Impossible checks passed to MapStringBoolLenBetween: "+
				"the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(v map[string]bool) error {
		if len(v) < low {
			return fmt.Errorf(
				"the number of entries in the map (%d)"+
					" must be between %d and %d - too short",
				len(v), low, high)
		}
		if len(v) > high {
			return fmt.Errorf(
				"the number of entries in the map (%d)"+
					" must be between %d and %d - too long",
				len(v), low, high)
		}
		return nil
	}
}

// MapStringBoolTrueCountEQ returns a check function that checks that the
// number of entries in the map set to true equals the supplied value
func MapStringBoolTrueCountEQ(i int) MapStringBool {
	return func(v map[string]bool) error {
		trueCount := 0
		for _, b := range v {
			if b {
				trueCount++
			}
		}

		if trueCount != i {
			return fmt.Errorf(
				"the number of entries set to true (%d) must equal %d",
				trueCount, i)
		}
		return nil
	}
}

// MapStringBoolTrueCountGT returns a check function that checks that the
// number of entries in the map set to true is greater than the supplied value
func MapStringBoolTrueCountGT(i int) MapStringBool {
	return func(v map[string]bool) error {
		trueCount := 0
		for _, b := range v {
			if b {
				trueCount++
			}
		}

		if trueCount <= i {
			return fmt.Errorf(
				"the number of entries set to true (%d)"+
					" must be greater than %d",
				trueCount, i)
		}
		return nil
	}
}

// MapStringBoolTrueCountLT returns a check function that checks that the
// number of entries in the map set to true is less than the supplied value
func MapStringBoolTrueCountLT(i int) MapStringBool {
	return func(v map[string]bool) error {
		trueCount := 0
		for _, b := range v {
			if b {
				trueCount++
			}
		}

		if trueCount >= i {
			return fmt.Errorf(
				"the number of entries set to true (%d) must be less than %d",
				trueCount, i)
		}
		return nil
	}
}

// MapStringBoolTrueCountBetween returns a check function that checks that
// the number of entries in the map set to true is between the two supplied
// values (inclusive)
func MapStringBoolTrueCountBetween(low, high int) MapStringBool {
	if low >= high {
		panic(fmt.Sprintf(
			"Impossible checks passed to MapStringBoolTrueCountBetween: "+
				"the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(v map[string]bool) error {
		trueCount := 0
		for _, b := range v {
			if b {
				trueCount++
			}
		}

		if trueCount < low {
			return fmt.Errorf(
				"the number of entries set to true (%d)"+
					" must be between %d and %d - too short",
				trueCount, low, high)
		}
		if trueCount > high {
			return fmt.Errorf(
				"the number of entries set to true (%d)"+
					" must be between %d and %d - too long",
				trueCount, low, high)
		}
		return nil
	}
}

// MapStringBoolOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func MapStringBoolOr(chkFuncs ...MapStringBool) MapStringBool {
	return func(v map[string]bool) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(v)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = _Or
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// MapStringBoolAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func MapStringBoolAnd(chkFuncs ...MapStringBool) MapStringBool {
	return func(v map[string]bool) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// MapStringBoolNot returns a function that will check that the value, when
// passed to the check func, does not pass it. You must also supply the error
// text to appear after the value that fails. This error text should be a
// string that describes the quality that the slice of strings should not
// have.
func MapStringBoolNot(c MapStringBool, errMsg string) MapStringBool {
	return func(v map[string]bool) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("%v should not be %s", v, errMsg)
	}
}
