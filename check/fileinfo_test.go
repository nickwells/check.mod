package check_test

import (
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
		testhelper.ID
		testhelper.ExpErr
		cf       check.FileInfo
		fileName string
	}{
		{
			ID:       testhelper.MkID("FileInfoIsRegular - good"),
			cf:       check.FileInfoIsRegular,
			fileName: "testdata/IsAFile",
		},
		{
			ID:       testhelper.MkID("FileInfoIsRegular - bad"),
			cf:       check.FileInfoIsRegular,
			fileName: "testdata/IsADirectory",
			ExpErr:   testhelper.MkExpErr("should be a regular file"),
		},
		{
			ID:       testhelper.MkID("FileInfoIsDir - good"),
			cf:       check.FileInfoIsDir,
			fileName: "testdata/IsADirectory",
		},
		{
			ID:       testhelper.MkID("FileInfoIsDir - bad"),
			cf:       check.FileInfoIsDir,
			fileName: "testdata/IsAFile",
			ExpErr:   testhelper.MkExpErr("should be a directory"),
		},
		{
			ID:       testhelper.MkID("FileInfoName - good"),
			cf:       check.FileInfoName(check.StringHasPrefix("IsA")),
			fileName: "testdata/IsAFile",
		},
		{
			ID:       testhelper.MkID("FileInfoName - bad"),
			cf:       check.FileInfoName(check.StringHasPrefix("XXX")),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix",
			),
		},
		{
			ID:       testhelper.MkID("FileInfoSize - good"),
			cf:       check.FileInfoSize(check.Int64EQ(0)),
			fileName: "testdata/IsAFile",
		},
		{
			ID:       testhelper.MkID("FileInfoSize - bad"),
			cf:       check.FileInfoSize(check.Int64EQ(99)),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the check on the size of '",
				"' failed: ",
				"the value ",
				" must be equal to 99"),
		},
		{
			ID:       testhelper.MkID("FileInfoPerm - good"),
			cf:       check.FileInfoPerm(check.FilePermEQ(0600)),
			fileName: fileWithSetPerms,
		},
		{
			ID:       testhelper.MkID("FileInfoPerm - bad"),
			cf:       check.FileInfoPerm(check.FilePermEQ(0666)),
			fileName: fileWithSetPerms,
			ExpErr: testhelper.MkExpErr(
				"the check on the permissions of '",
				"' failed: ",
				"the permissions ",
				" should be equal to 0666"),
		},
		{
			ID:       testhelper.MkID("FileInfoMode - good"),
			cf:       check.FileInfoMode(os.ModeDir),
			fileName: "testdata/IsADirectory",
		},
		{
			ID:       testhelper.MkID("FileInfoMode - bad"),
			cf:       check.FileInfoMode(os.ModeDir),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"should have been a directory but was a regular file"),
		},
		{
			ID:       testhelper.MkID("FileInfoModTime - good"),
			cf:       check.FileInfoModTime(check.TimeLT(timeAfterModTime)),
			fileName: fileWithKnownInfo,
		},
		{
			ID:       testhelper.MkID("FileInfoModTime - bad"),
			cf:       check.FileInfoModTime(check.TimeLT(timeBeforeModTime)),
			fileName: fileWithKnownInfo,
			ExpErr: testhelper.MkExpErr(
				"the check on the modification time of '",
				"' failed: ",
				"the time ",
				" should be before"),
		},
		{
			ID: testhelper.MkID("FileInfoAnd - good"),
			cf: check.FileInfoAnd(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				check.FileInfoName(check.StringHasSuffix("File")),
			),
			fileName: "testdata/IsAFile",
		},
		{
			ID: testhelper.MkID("FileInfoAnd - bad"),
			cf: check.FileInfoAnd(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				check.FileInfoName(check.StringHasPrefix("XXX")),
			),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix"),
		},
		{
			ID: testhelper.MkID("FileInfoOr - good"),
			cf: check.FileInfoOr(
				check.FileInfoName(check.StringHasPrefix("XXX")),
				check.FileInfoName(check.StringHasSuffix("File")),
			),
			fileName: "testdata/IsAFile",
		},
		{
			ID: testhelper.MkID("FileInfoOr - bad"),
			cf: check.FileInfoOr(
				check.FileInfoName(check.StringHasPrefix("YYY")),
				check.FileInfoName(check.StringHasPrefix("XXX")),
			),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the check on the name of '",
				"' failed: ",
				"should have 'XXX' as a prefix",
				" or ",
				"should have 'YYY' as a prefix"),
		},
		{
			ID: testhelper.MkID("FileInfoNot - good"),
			cf: check.FileInfoNot(
				check.FileInfoName(check.StringHasPrefix("XXX")),
				"test"),
			fileName: "testdata/IsAFile",
		},
		{
			ID: testhelper.MkID("FileInfoNot - bad"),
			cf: check.FileInfoNot(
				check.FileInfoName(check.StringHasPrefix("IsA")),
				"should not have a prefix of 'IsA'"),
			fileName: "testdata/IsAFile",
			ExpErr:   testhelper.MkExpErr("should not have a prefix of 'IsA'"),
		},
	}

	for _, tc := range testCases {
		fi, err := os.Stat(tc.fileName)
		if err != nil {
			t.Log(tc.IDStr())
			t.Fatalf("\t: cannot get the FileInfo from file: %s err: %s",
				tc.fileName, err)
		}

		err = tc.cf(fi)
		testhelper.CheckExpErr(t, err, tc)
	}
}
