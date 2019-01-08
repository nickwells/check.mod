// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package check

import (
	"fmt"
	"os"
	"syscall"
)

// FileInfoOwnedBySelf tests that the file is owned by the calling user
func FileInfoOwnedBySelf(fi os.FileInfo) error {
	stat := fi.Sys().(*syscall.Stat_t) // will panic if Sys() doesn't
	// return a Stat_t
	if stat.Uid == uint32(os.Getuid()) {
		return nil
	}
	return fmt.Errorf("'%s' - should have been owned by the user -"+
		" the user ID should have been %d but was %d",
		fi.Name(), os.Getuid(), stat.Uid)
}

// FileInfoUidEQ returns a FileInfo (func) that tests that the file is owned
// by the specified user
func FileInfoUidEQ(uid uint32) FileInfo {
	return func(fi os.FileInfo) error {
		stat := fi.Sys().(*syscall.Stat_t) // will panic if Sys() doesn't
		// return a Stat_t
		if stat.Uid == uid {
			return nil
		}
		return fmt.Errorf("'%s' - the user ID should have been %d but was %d",
			fi.Name(), uid, stat.Uid)
	}
}

// FileInfoGidEQ returns a FileInfo (func) that tests that the file is owned
// by the specified user
func FileInfoGidEQ(gid uint32) FileInfo {
	return func(fi os.FileInfo) error {
		stat := fi.Sys().(*syscall.Stat_t) // will panic if Sys() doesn't
		// return a Stat_t
		if stat.Gid == gid {
			return nil
		}
		return fmt.Errorf("'%s' - the group ID should have been %d but was %d",
			fi.Name(), gid, stat.Gid)
	}
}
