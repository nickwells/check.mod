package check_test

import (
	"os"
	"testing"
	"time"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFileInfo(t *testing.T) {
	fileWithSetPerms := "testdata/IsAFile.PBits0600"
	_ = os.Chmod(fileWithSetPerms, 0o600) // force the file mode

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
		cf       check.ValCk[os.FileInfo]
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
			cf:       check.FileInfoName(check.StringHasPrefix[string]("IsA")),
			fileName: "testdata/IsAFile",
		},
		{
			ID:       testhelper.MkID("FileInfoName - bad"),
			cf:       check.FileInfoName(check.StringHasPrefix[string]("XXX")),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the file name ",
				" is incorrect: ",
				`should have "XXX" as a prefix`,
			),
		},
		{
			ID:       testhelper.MkID("FileInfoSize - good"),
			cf:       check.FileInfoSize(check.ValEQ(int64(0))),
			fileName: "testdata/IsAFile",
		},
		{
			ID:       testhelper.MkID("FileInfoSize - bad"),
			cf:       check.FileInfoSize(check.ValEQ(int64(99))),
			fileName: "testdata/IsAFile",
			ExpErr: testhelper.MkExpErr(
				"the check on the size of ",
				" failed: ",
				"the value ",
				" must equal 99"),
		},
		{
			ID:       testhelper.MkID("FileInfoPerm - good"),
			cf:       check.FileInfoPerm(check.FilePermEQ(0o600)),
			fileName: fileWithSetPerms,
		},
		{
			ID:       testhelper.MkID("FileInfoPerm - bad"),
			cf:       check.FileInfoPerm(check.FilePermEQ(0o666)),
			fileName: fileWithSetPerms,
			ExpErr: testhelper.MkExpErr(
				"the file permissions of ",
				" are incorrect: ",
				"the permissions ",
				" should equal 0666"),
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
				"the modification time of ",
				" is incorrect: ",
				"the time ",
				" must be before"),
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
