package check

import (
	"fmt"
	"io/fs"
)

// FilePermEQ returns a function that will check that the file permission is
// set to the value of the perms parameter
func FilePermEQ(perms fs.FileMode) ValCk[fs.FileMode] {
	return func(fm fs.FileMode) error {
		if fm.Perm() == perms {
			return nil
		}

		return fmt.Errorf("the permissions (%04o) should equal %04o",
			fm.Perm(), perms)
	}
}

// FilePermHasAll returns a function that will check that the file permission
// has all of the permissions set in the perms parameter
func FilePermHasAll(perms fs.FileMode) ValCk[fs.FileMode] {
	return func(fm fs.FileMode) error {
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
func FilePermHasNone(perms fs.FileMode) ValCk[fs.FileMode] {
	return func(fm fs.FileMode) error {
		if (fm.Perm() & perms) == 0 {
			return nil
		}

		return fmt.Errorf(
			"the permissions (%04o)"+
				" should have none of the permissions in %04o",
			fm.Perm(), perms)
	}
}
