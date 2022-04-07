package check

import (
	"fmt"
	"os"
	"time"
)

// FileInfoSize returns a function that will check that the file size passes
// the test specified by the passed Int64 check
func FileInfoSize(cf ValCk[int64]) ValCk[os.FileInfo] {
	return func(fi os.FileInfo) error {
		err := cf(fi.Size())
		if err != nil {
			return fmt.Errorf("the check on the size of %q failed: %w",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoPerm returns a function that will check that the file permissions
// pass the test specified by the passed FilePerm check
func FileInfoPerm(cf ValCk[os.FileMode]) ValCk[os.FileInfo] {
	return func(fi os.FileInfo) error {
		err := cf(fi.Mode())
		if err != nil {
			return fmt.Errorf("the file permissions of %q are incorrect: %w",
				fi.Name(), err)
		}
		return nil
	}
}

// FileInfoName returns a function that will check that the file name
// passes the test specified by the passed String check
func FileInfoName(cf ValCk[string]) ValCk[os.FileInfo] {
	return func(fi os.FileInfo) error {
		err := cf(fi.Name())
		if err != nil {
			return fmt.Errorf("the file name %q is incorrect: %w",
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
	return fmt.Errorf("%q should be a directory", fi.Name())
}

// FileInfoIsRegular will check that the file info describes a regular file
func FileInfoIsRegular(fi os.FileInfo) error {
	if fi.Mode().IsRegular() {
		return nil
	}
	return fmt.Errorf("%q should be a regular file", fi.Name())
}

// FileInfoMode returns a function that will check that the file mode type
// matches the value passed. Typically this will be a single value as
// enumerated in the os package under FileMode but if more than one value is
// given (and'ed together as a bitmask) then if any of the bits is set this
// function will return nil if any of the bits is set. This allows you to
// check for several types at once. If a zero value is passed it will check
// for a regular file - none of the bits are set.
func FileInfoMode(m os.FileMode) ValCk[os.FileInfo] {
	return func(fi os.FileInfo) error {
		typeBits := fi.Mode() & os.ModeType
		if typeBits == m ||
			typeBits&m != 0 {
			return nil
		}
		return fmt.Errorf(
			"%q should have been %s but was %s",
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
		sep = _Or
	}
	if m&os.ModeSymlink == os.ModeSymlink {
		name += sep + "a symlink"
		sep = _Or
	}
	if m&os.ModeNamedPipe == os.ModeNamedPipe {
		name += sep + "a named pipe"
		sep = _Or
	}
	if m&os.ModeSocket == os.ModeSocket {
		name += sep + "a socket"
		sep = _Or
	}
	if m&os.ModeDevice == os.ModeDevice {
		name += sep + "a device"
		sep = _Or
	}
	if m&os.ModeIrregular == os.ModeIrregular {
		name += sep + "a non-regular file"
	}
	return name
}

// FileInfoModTime returns a function that will check that the file
// modification time passes the test specified by the passed Time check
func FileInfoModTime(cf ValCk[time.Time]) ValCk[os.FileInfo] {
	return func(fi os.FileInfo) error {
		err := cf(fi.ModTime())
		if err != nil {
			return fmt.Errorf(
				"the modification time of %q is incorrect: %s",
				fi.Name(), err)
		}
		return nil
	}
}
