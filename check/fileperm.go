package check

import (
	"fmt"
	"os"
)

// FilePerm is the type of a check function which takes an os.FileMode
// parameter and returns an error or nil if the check passes
type FilePerm func(d os.FileMode) error

// FilePermEQ returns a function that will check that the file permission is
// set to the value of the perms parameter
func FilePermEQ(perms os.FileMode) FilePerm {
	return func(fm os.FileMode) error {
		if fm.Perm() == perms {
			return nil
		}
		return fmt.Errorf("the permissions (%04o) should be equal to %04o",
			fm.Perm(), perms)
	}
}

// FilePermHasAll returns a function that will check that the file permission
// has all of the permissions set in the perms parameter
func FilePermHasAll(perms os.FileMode) FilePerm {
	return func(fm os.FileMode) error {
		if (fm.Perm() & perms) == perms {
			return nil
		}
		return fmt.Errorf(
			"the permissions (%04o) should have all of the permissions in %04o",
			fm.Perm(), perms)
	}
}

// FilePermHasNone returns a function that will check that the file permission
// has none of the permissions set in the perms parameter
func FilePermHasNone(perms os.FileMode) FilePerm {
	return func(fm os.FileMode) error {
		if (fm.Perm() & perms) == 0 {
			return nil
		}
		return fmt.Errorf(
			"the permissions (%04o)"+
				" should have none of the permissions in %04o",
			fm.Perm(), perms)
	}
}

// FilePermOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func FilePermOr(chkFuncs ...FilePerm) FilePerm {
	return func(v os.FileMode) error {
		compositeErr := ""
		sep := "("

		for _, cf := range chkFuncs {
			err := cf(v)
			if err == nil {
				return nil
			}

			compositeErr += sep + err.Error()
			sep = " OR "
		}
		return fmt.Errorf("%s)", compositeErr)
	}
}

// FilePermAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func FilePermAnd(chkFuncs ...FilePerm) FilePerm {
	return func(v os.FileMode) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// FilePermNot returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails
func FilePermNot(c FilePerm, errMsg string) FilePerm {
	return func(v os.FileMode) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("the permissions (%04o) %s", v.Perm(), errMsg)
	}
}
