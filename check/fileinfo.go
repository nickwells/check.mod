package check

import (
	"fmt"
	"os"
)

// FileInfo is the type of a check function which takes an os.FileInfo
// parameter and returns an error or nil if the check passes
type FileInfo func(d os.FileInfo) error

// FileInfoSize returns a function that will check that the file size passes
// the test specified by the passed Int64 check
func FileInfoSize(c Int64) FileInfo {
	return func(fi os.FileInfo) error {
		err := c(fi.Size())
		if err != nil {
			return fmt.Errorf("the check on the size of '%s' failed: %s",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoPerm returns a function that will check that the file permissions
// pass the test specified by the passed FilePerm check
func FileInfoPerm(c FilePerm) FileInfo {
	return func(fi os.FileInfo) error {
		err := c(fi.Mode())
		if err != nil {
			return fmt.Errorf("the check on the permissions of '%s' failed: %s",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoName returns a function that will check that the file name
// passes the test specified by the passed String check
func FileInfoName(c String) FileInfo {
	return func(fi os.FileInfo) error {
		err := c(fi.Name())
		if err != nil {
			return fmt.Errorf("the check on the name of '%s' failed: %s",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoIsDir will check that the file info describes a directory
func FileInfoIsDir(fi os.FileInfo) error {
	if fi.IsDir() {
		return nil
	}
	return fmt.Errorf("'%s' should be a directory", fi.Name())
}

// FileInfoIsRegular will check that the file info describes a regular file
func FileInfoIsRegular(fi os.FileInfo) error {
	if fi.Mode().IsRegular() {
		return nil
	}
	return fmt.Errorf("'%s' should be a regular file", fi.Name())
}

// FileInfoMode returns a function that will check that the file mode type
// matches the value passed. Typically this will be a single value as
// enumerated in the os package under FileMode but if more than one value is
// given (and'ed together as a bitmask) then if any of the bits is set this
// function will return nil if any of the bits is set. This allows you to
// check for several types at once. If a zero value is passed it will check
// for a regular file - none of the bits are set.
func FileInfoMode(m os.FileMode) FileInfo {
	return func(fi os.FileInfo) error {
		typeBits := fi.Mode() & os.ModeType
		if typeBits == m ||
			typeBits&m != 0 {
			return nil
		}
		return fmt.Errorf(
			"'%s' should have been %s but was %s",
			fi.Name(), modeName(m), modeName(typeBits))
	}
}

// modeName reports the bits set in the file mode value in a human-readable
// form
func modeName(m os.FileMode) string {
	if m == 0 {
		return "a regular file"
	}

	sep := ""
	name := ""

	if m&os.ModeDir == os.ModeDir {
		name += sep + "a directory"
		sep = " or "
	}
	if m&os.ModeSymlink == os.ModeSymlink {
		name += sep + "a symlink"
		sep = " or "
	}
	if m&os.ModeNamedPipe == os.ModeNamedPipe {
		name += sep + "a named pipe"
		sep = " or "
	}
	if m&os.ModeSocket == os.ModeSocket {
		name += sep + "a socket"
		sep = " or "
	}
	if m&os.ModeDevice == os.ModeDevice {
		name += sep + "a device"
		sep = " or "
	}
	if m&os.ModeIrregular == os.ModeIrregular {
		name += sep + "a non-regular file"
		sep = " or "
	}
	return name
}

// FileInfoModTime returns a function that will check that the file
// modification time passes the test specified by the passed Time check
func FileInfoModTime(c Time) FileInfo {
	return func(fi os.FileInfo) error {
		err := c(fi.ModTime())
		if err != nil {
			return fmt.Errorf(
				"the check on the modification time of '%s' failed: %s",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoOr returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes at least one of them
func FileInfoOr(chkFuncs ...FileInfo) FileInfo {
	return func(v os.FileInfo) error {
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

// FileInfoAnd returns a function that will check that the value, when
// passed to each of the check funcs in turn, passes all of them
func FileInfoAnd(chkFuncs ...FileInfo) FileInfo {
	return func(v os.FileInfo) error {
		for _, cf := range chkFuncs {
			err := cf(v)
			if err != nil {
				return err
			}
		}
		return nil
	}
}

// FileInfoNot returns a function that will check that the value, when passed
// to the check func, does not pass it. You must also supply the error text
// to appear after the value that fails
func FileInfoNot(c FileInfo, errMsg string) FileInfo {
	return func(v os.FileInfo) error {
		err := c(v)
		if err != nil {
			return nil
		}

		return fmt.Errorf("'%s' %s", v.Name(), errMsg)
	}
}
