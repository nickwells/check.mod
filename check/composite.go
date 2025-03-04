package check

import "fmt"

// Or returns a function that will check that the value, when passed to each
// of the check funcs in turn, passes at least one of them. If any check
// passes a nil error is returned. The error returned (if any) will show all
// the failing checks.
func Or[T any](chkFuncs ...ValCk[T]) ValCk[T] {
	return func(v T) error {
		compositeErr := ""
		sep := "either ["

		for _, cf := range chkFuncs {
			err := cf(v)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = "] or ["
		}

		return fmt.Errorf("%s]", compositeErr)
	}
}

// And returns a function that will check that the value, when passed to each
// of the check funcs in turn, passes all of them. The error from the first
// check to fail is returned.
func And[T any](chkFuncs ...ValCk[T]) ValCk[T] {
	return func(v T) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

// Not returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails. This error text should be a string
// that describes the quality that the number should not have. So, for
// instance, if the function being Not'ed was
//
//	check.ValGT[T any](5)
//
// then the errMsg parameter should be
//
//	"a number greater than 5".
func Not[T any](c ValCk[T], errMsg string) ValCk[T] {
	return func(v T) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("%v should not be %s", v, errMsg)
	}
}
