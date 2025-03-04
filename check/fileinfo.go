package check

import (
	"fmt"
	"io/fs"
	"time"
)

// FileInfoSize returns a function that will check that the file size passes
// the test specified by the passed Int64 check
func FileInfoSize(cf ValCk[int64]) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
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
func FileInfoPerm(cf ValCk[fs.FileMode]) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
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
func FileInfoName(cf ValCk[string]) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
		err := cf(fi.Name())
		if err != nil {
			return fmt.Errorf("the file name %q is incorrect: %w",
				fi.Name(), err)
		}

		return nil
	}
}

// FileInfoIsDir will check that the file info describes a directory
func FileInfoIsDir(fi fs.FileInfo) error {
	if fi.IsDir() {
		return nil
	}

	return fmt.Errorf("%q should be a directory", fi.Name())
}

// FileInfoIsRegular will check that the file info describes a regular file
func FileInfoIsRegular(fi fs.FileInfo) error {
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
func FileInfoMode(m fs.FileMode) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
		typeBits := fi.Mode() & fs.ModeType
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
func modeName(m fs.FileMode) string {
	if m == 0 {
		return "a regular file"
	}

	sep := ""
	name := ""

	if m&fs.ModeDir == fs.ModeDir {
		name += sep + "a directory"
		sep = " or "
	}

	if m&fs.ModeSymlink == fs.ModeSymlink {
		name += sep + "a symlink"
		sep = " or "
	}

	if m&fs.ModeNamedPipe == fs.ModeNamedPipe {
		name += sep + "a named pipe"
		sep = " or "
	}

	if m&fs.ModeSocket == fs.ModeSocket {
		name += sep + "a socket"
		sep = " or "
	}

	if m&fs.ModeDevice == fs.ModeDevice {
		name += sep + "a device"
		sep = " or "
	}

	if m&fs.ModeIrregular == fs.ModeIrregular {
		name += sep + "a non-regular file"
	}

	return name
}

// FileInfoModTime returns a function that will check that the file
// modification time passes the test specified by the passed Time check
func FileInfoModTime(cf ValCk[time.Time]) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
		err := cf(fi.ModTime())
		if err != nil {
			return fmt.Errorf(
				"the modification time of %q is incorrect: %s",
				fi.Name(), err)
		}

		return nil
	}
}
