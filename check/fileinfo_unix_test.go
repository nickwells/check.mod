// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package check_test

import (
	"fmt"
	"os"
	"syscall"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestFileInfoUnix(t *testing.T) {
	filename := "testdata/IsAFile"
	fi, err := os.Stat(filename)
	if err != nil {
		t.Fatalf("\t: cannot get the FileInfo from file: %s err: %s",
			filename, err)
	}
	stat := fi.Sys().(*syscall.Stat_t)

	testCases := []struct {
		name           string
		cf             check.FileInfo
		errExpected    bool
		errMustContain []string
	}{
		{
			name: "good uid",
			cf:   check.FileInfoUidEQ(stat.Uid),
		},
		{
			name:           "bad uid",
			cf:             check.FileInfoUidEQ(stat.Uid + 1),
			errExpected:    true,
			errMustContain: []string{"the user ID should have been"},
		},
		{
			name: "good gid",
			cf:   check.FileInfoGidEQ(stat.Gid),
		},
		{
			name:           "bad gid",
			cf:             check.FileInfoGidEQ(stat.Gid + 1),
			errExpected:    true,
			errMustContain: []string{"the group ID should have been"},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s :", i, tc.name)

		err = tc.cf(fi)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}
}
