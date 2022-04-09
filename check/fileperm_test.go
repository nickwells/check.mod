package check_test

import (
	"io/fs"
	"testing"

	"github.com/nickwells/check.mod/v2/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFilePerm(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fileMode fs.FileMode
		cf       check.ValCk[fs.FileMode]
	}{
		{
			ID:       testhelper.MkID("EQ - good"),
			fileMode: 0o777,
			cf:       check.FilePermEQ(0o777),
		},
		{
			ID:       testhelper.MkID("EQ - bad"),
			fileMode: 0o777,
			cf:       check.FilePermEQ(0o770),
			ExpErr: testhelper.MkExpErr(
				"the permissions (0777) should equal 0770"),
		},
		{
			ID:       testhelper.MkID("HasAll - good"),
			fileMode: 0o777,
			cf:       check.FilePermHasAll(0o111),
		},
		{
			ID:       testhelper.MkID("HasAll - bad - has none"),
			fileMode: 0o666,
			cf:       check.FilePermHasAll(0o111),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have all of the permissions in 0111"),
		},
		{
			ID:       testhelper.MkID("HasAll - bad - has some"),
			fileMode: 0o666,
			cf:       check.FilePermHasAll(0o221),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have all of the permissions in 0221"),
		},
		{
			ID:       testhelper.MkID("HasNone - good"),
			fileMode: 0o666,
			cf:       check.FilePermHasNone(0o111),
		},
		{
			ID:       testhelper.MkID("HasNone - bad - has all"),
			fileMode: 0o777,
			cf:       check.FilePermHasNone(0o111),
			ExpErr: testhelper.MkExpErr("the permissions (0777)" +
				" should have none of the permissions in 0111"),
		},
		{
			ID:       testhelper.MkID("HasNone - bad - has some"),
			fileMode: 0o666,
			cf:       check.FilePermHasNone(0o211),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have none of the permissions in 0211"),
		},
	}

	for _, tc := range testCases {
		err := tc.cf(tc.fileMode)
		testhelper.CheckExpErr(t, err, tc)
	}
}
