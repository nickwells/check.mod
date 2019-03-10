package check_test

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestFileInfo(t *testing.T) {
	fileWithSetPerms := "testdata/IsAFile.PBits0600"
	_ = os.Chmod(fileWithSetPerms, 0600) // force the file mode

	fileWithKnownInfo := "testdata/IsAFile"
	fi, err := os.Stat(fileWithKnownInfo)
	if err != nil {
		t.Fatalf("pre-test setup:"+
			" Cannot get the FileInfo from file: %s err: %s",
			fileWithKnownInfo, err)
		return
	}
	modTime := fi.ModTime()
	timeBeforeModTime := modTime.Add(-60 * time.Second)
	timeAfterModTime := modTime.Add(60 * time.Second)

	testCases := []struct {
		name           string
		cf             check.FileInfo
		fileName       string
		errExpected    bool
		errMustContain []string
	}{
		{
			name:     "FileInfoIsRegular - good",
			cf:       check.FileInfoIsRegular,
			fileName: "testdata/IsAFile",
		},
		{
			name:           "FileInfoIsRegular - bad",
			cf:             check.FileInfoIsRegular,
			fileName:       "testdata/IsADirectory",
			errExpected:    true,
			errMustContain: []string{"should be a regular file"},
		},
		{
			name:     "FileInfoIsDir - good",
			cf:       check.FileInfoIsDir,
			fileName: "testdata/IsADirectory",
		},
		{
			name:           "FileInfoIsDir - bad",
			cf:             check.FileInfoIsDir,
			fileName:       "testdata/IsAFile",
			errExpected:    true,
			errMustContain: []string{"should be a directory"},
		},
		{
			name:     "FileInfoName - good",
			cf:       check.FileInfoName(check.StringHasPrefix("IsA")),
			fileName: "testdata/IsAFile",
		},
		{
			name:        "FileInfoName - bad",
			cf:          check.FileInfoName(check.StringHasPrefix("XXX")),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix",
			},
		},
		{
			name:     "FileInfoSize - good",
			cf:       check.FileInfoSize(check.Int64EQ(0)),
			fileName: "testdata/IsAFile",
		},
		{
			name:        "FileInfoSize - bad",
			cf:          check.FileInfoSize(check.Int64EQ(99)),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"the check on the size of '",
				"' failed: ",
				"the value ",
				" must be equal to 99"},
		},
		{
			name:     "FileInfoPerm - good",
			cf:       check.FileInfoPerm(check.FilePermEQ(0600)),
			fileName: fileWithSetPerms,
		},
		{
			name:        "FileInfoPerm - bad",
			cf:          check.FileInfoPerm(check.FilePermEQ(0666)),
			fileName:    fileWithSetPerms,
			errExpected: true,
			errMustContain: []string{
				"the check on the permissions of '",
				"' failed: ",
				"the permissions ",
				" should be equal to 0666"},
		},
		{
			name:     "FileInfoMode - good",
			cf:       check.FileInfoMode(os.ModeDir),
			fileName: "testdata/IsADirectory",
		},
		{
			name:        "FileInfoMode - bad",
			cf:          check.FileInfoMode(os.ModeDir),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"should have been a directory but was a regular file"},
		},
		{
			name:     "FileInfoModTime - good",
			cf:       check.FileInfoModTime(check.TimeLT(timeAfterModTime)),
			fileName: fileWithKnownInfo,
		},
		{
			name:        "FileInfoModTime - bad",
			cf:          check.FileInfoModTime(check.TimeLT(timeBeforeModTime)),
			fileName:    fileWithKnownInfo,
			errExpected: true,
			errMustContain: []string{
				"the check on the modification time of '",
				"' failed: ",
				"the time ",
				" should be before"},
		},
		{
			name: "FileInfoAnd - good",
			cf: check.FileInfoAnd(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				check.FileInfoName(check.StringHasSuffix("File")),
			),
			fileName: "testdata/IsAFile",
		},
		{
			name: "FileInfoAnd - bad",
			cf: check.FileInfoAnd(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				check.FileInfoName(check.StringHasPrefix("XXX")),
			),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix",
			},
		},
		{
			name: "FileInfoOr - good",
			cf: check.FileInfoOr(
				check.FileInfoName(check.StringHasPrefix("XXX")),
				check.FileInfoName(check.StringHasSuffix("File")),
			),
			fileName: "testdata/IsAFile",
		},
		{
			name: "FileInfoOr - bad",
			cf: check.FileInfoOr(
				check.FileInfoName(check.StringHasPrefix("YYY")),
				check.FileInfoName(check.StringHasPrefix("XXX")),
			),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix",
				" or ",
				"should have 'YYY' as a prefix",
			},
		},
		{
			name: "FileInfoNot - good",
			cf: check.FileInfoNot(
				check.FileInfoName(check.StringHasPrefix("XXX")),
				"test"),
			fileName: "testdata/IsAFile",
		},
		{
			name: "FileInfoNot - bad",
			cf: check.FileInfoNot(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				"should not have a prefix of 'IsA'"),
			fileName:    "testdata/IsAFile",
			errExpected: true,
			errMustContain: []string{
				"should not have a prefix of 'IsA'",
			},
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		fi, err := os.Stat(tc.fileName)
		if err != nil {
			t.Log(tcID)
			t.Fatalf("\t: cannot get the FileInfo from file: %s err: %s",
				tc.fileName, err)
		}

		err = tc.cf(fi)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}
}
