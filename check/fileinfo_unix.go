//go:build darwin || dragonfly || freebsd || linux || netbsd || openbsd || solaris

package check

import (
	"fmt"
	"io/fs"
	"os"
	"syscall"
)

// FileInfoOwnedBySelf tests that the file is owned by the calling user
func FileInfoOwnedBySelf(fi fs.FileInfo) error {
	// the next call will panic if Sys() doesn't return a Stat_t
	stat := fi.Sys().(*syscall.Stat_t)
	if stat.Uid == uint32(os.Getuid()) { //nolint:gosec
		return nil
	}

	return fmt.Errorf("%q - should have been owned by the user -"+
		" the user ID should have been %d but was %d",
		fi.Name(), os.Getuid(), stat.Uid)
}

// FileInfoUidEQ returns a function that tests that the file is owned
// by the specified user
func FileInfoUidEQ(uid uint32) ValCk[fs.FileInfo] { //nolint:revive
	return func(fi fs.FileInfo) error {
		// the next call will panic if Sys() doesn't return a Stat_t
		stat := fi.Sys().(*syscall.Stat_t)
		if stat.Uid == uid {
			return nil
		}

		return fmt.Errorf("%q - the user ID should have been %d but was %d",
			fi.Name(), uid, stat.Uid)
	}
}

// FileInfoGidEQ returns a function that tests that the file is owned
// by the specified user
func FileInfoGidEQ(gid uint32) ValCk[fs.FileInfo] {
	return func(fi fs.FileInfo) error {
		stat := fi.Sys().(*syscall.Stat_t) // will panic if Sys() doesn't
		// return a Stat_t
		if stat.Gid == gid {
			return nil
		}

		return fmt.Errorf("%q - the group ID should have been %d but was %d",
			fi.Name(), gid, stat.Gid)
	}
}
