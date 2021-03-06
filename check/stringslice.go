package check

import "fmt"

// StringSlice is the type of a check function for a slice of strings. It
// takes a slice of strings as a parameter and returns an error or nil if
// there is no error
type StringSlice func(v []string) error

// StringSliceStringCheck returns a check function that checks that every
// entry in the list passes the supplied String check func
func StringSliceStringCheck(sc String) StringSlice {
	return func(v []string) error {
		for i, s := range v {
			if err := sc(s); err != nil {
				return fmt.Errorf(
					"list entry: %d (%s) does not pass the test: %s",
					i, s, err)
			}
		}
		return nil
	}
}

// StringSliceStringCheckByPos returns a check function that checks that a
// given entry in the list passes the corresponding supplied String check
// func. If there are more entries in the slice than check.String functions
// then the final supplied check is applied. If there are fewer entries than
// functions then the excess checks are not applied. Note that you could
// choose to combine this check with a test of the length of the slice by
// means of a StringSliceAnd check.
func StringSliceStringCheckByPos(scs ...String) StringSlice {
	return func(v []string) error {
		if len(scs) == 0 {
			return nil
		}

		for i, s := range v {
			if i >= len(scs) {
				i = len(scs) - 1
			}
			if err := scs[i](s); err != nil {
				return fmt.Errorf(
					"list entry: %d (%s) does not pass the test: %s",
					i, s, err)
			}
		}
		return nil
	}
}

// StringSliceContains returns a check function that checks that at least one
// entry in the list matches the supplied String check func. The condition
// parameter should describe the check that is being performed. For instance,
// if the check is that the string length must be greater than 5 characters
// then the condition parameter should be
//     "the string should be greater than 5 characters"
func StringSliceContains(sc String, condition string) StringSlice {
	return func(v []string) error {
		for _, s := range v {
			if err := sc(s); err == nil {
				return nil
			}
		}
		return fmt.Errorf("none of the list entries passes the test: %s",
			condition)
	}
}

// StringSliceLenEQ returns a check function that checks that the length of
// the list equals the supplied value
func StringSliceLenEQ(i int) StringSlice {
	return func(v []string) error {
		if len(v) != i {
			return fmt.Errorf("the length of the list (%d) must equal %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenGT returns a check function that checks that the length of
// the list is greater than the supplied value
func StringSliceLenGT(i int) StringSlice {
	return func(v []string) error {
		if len(v) <= i {
			return fmt.Errorf(
				"the length of the list (%d) must be greater than %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenLT returns a check function that checks that the length of
// the list is less than the supplied value
func StringSliceLenLT(i int) StringSlice {
	return func(v []string) error {
		if len(v) >= i {
			return fmt.Errorf(
				"the length of the list (%d) must be less than %d",
				len(v), i)
		}
		return nil
	}
}

// StringSliceLenBetween returns a check function that checks that the length
// of the list is between the two supplied values (inclusive)
func StringSliceLenBetween(low, high int) StringSlice {
	if low >= high {
		panic(fmt.Sprintf("Impossible checks passed to StringSliceLenBetween: "+
			"the lower limit (%d) should be less than the upper limit (%d)",
			low, high))
	}

	return func(v []string) error {
		if len(v) < low {
			return fmt.Errorf(
				"the length of the list (%d) must be between %d and %d"+
					" - too short",
				len(v), low, high)
		}
		if len(v) > high {
			return fmt.Errorf(
				"the length of the list (%d) must be between %d and %d"+
					" - too long",
				len(v), low, high)
		}
		return nil
	}
}

// StringSliceNoDups checks that the list contains no duplicates
func StringSliceNoDups(v []string) error {
	dupMap := make(map[string]int)
	for i, s := range v {
		if dup, ok := dupMap[s]; ok {
			return fmt.Errorf(
				"list entries: %d and %d are duplicates, both are: '%s'",
				dup, i, s)
		}
		dupMap[s] = i
	}
	return nil
}

// StringSliceOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func StringSliceOr(chkFuncs ...StringSlice) StringSlice {
	return func(v []string) error {
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

// StringSliceAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func StringSliceAnd(chkFuncs ...StringSlice) StringSlice {
	return func(v []string) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// StringSliceNot returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails. This error text should be a string
// that describes the quality that the slice of strings should not have.
func StringSliceNot(c StringSlice, errMsg string) StringSlice {
	return func(v []string) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("%v should not be %s", v, errMsg)
	}
}
