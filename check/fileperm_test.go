package check_test

import (
	"os"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestFilePerm(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		testhelper.ExpErr
		fileMode os.FileMode
		cf       check.FilePerm
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
				"the permissions (0777) should be equal to 0770"),
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
		{
			ID:       testhelper.MkID("Or - bad - none pass"),
			fileMode: 0o666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0o211),
				check.FilePermHasNone(0o121),
				check.FilePermHasNone(0o112),
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)"+
				" should have none of the permissions in 0211",
				" should have none of the permissions in 0121",
				" should have none of the permissions in 0112",
				" or "),
		},
		{
			ID:       testhelper.MkID("Or - good - one passes"),
			fileMode: 0o666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0o211),
				check.FilePermHasNone(0o121),
				check.FilePermHasNone(0o111),
			),
		},
		{
			ID:       testhelper.MkID("And - bad - one fails"),
			fileMode: 0o666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(0o211),
				check.FilePermHasNone(0o111),
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have none of the permissions in 0211"),
		},
		{
			ID:       testhelper.MkID("And - good - all pass"),
			fileMode: 0o666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(0o1),
				check.FilePermHasNone(0o10),
				check.FilePermHasNone(0o100),
			),
		},
		{
			ID:       testhelper.MkID("Not - bad - one fails"),
			fileMode: 0o666,
			cf: check.FilePermNot(
				check.FilePermHasAll(0o200),
				"have owner write permission set",
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should not have owner write permission set"),
		},
		{
			ID:       testhelper.MkID("Not - good"),
			fileMode: 0o666,
			cf: check.FilePermNot(
				check.FilePermHasAll(0o1),
				"should not have other execute permission set",
			),
		},
	}

	for _, tc := range testCases {
		err := tc.cf(tc.fileMode)
		testhelper.CheckExpErr(t, err, tc)
	}
}
