package check_test

import (
	"os"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
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
			fileMode: 0777,
			cf:       check.FilePermEQ(0777),
		},
		{
			ID:       testhelper.MkID("EQ - bad"),
			fileMode: 0777,
			cf:       check.FilePermEQ(0770),
			ExpErr: testhelper.MkExpErr(
				"the permissions (0777) should be equal to 0770"),
		},
		{
			ID:       testhelper.MkID("HasAll - good"),
			fileMode: 0777,
			cf:       check.FilePermHasAll(0111),
		},
		{
			ID:       testhelper.MkID("HasAll - bad - has none"),
			fileMode: 0666,
			cf:       check.FilePermHasAll(0111),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have all of the permissions in 0111"),
		},
		{
			ID:       testhelper.MkID("HasAll - bad - has some"),
			fileMode: 0666,
			cf:       check.FilePermHasAll(0221),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have all of the permissions in 0221"),
		},
		{
			ID:       testhelper.MkID("HasNone - good"),
			fileMode: 0666,
			cf:       check.FilePermHasNone(0111),
		},
		{
			ID:       testhelper.MkID("HasNone - bad - has all"),
			fileMode: 0777,
			cf:       check.FilePermHasNone(0111),
			ExpErr: testhelper.MkExpErr("the permissions (0777)" +
				" should have none of the permissions in 0111"),
		},
		{
			ID:       testhelper.MkID("HasNone - bad - has some"),
			fileMode: 0666,
			cf:       check.FilePermHasNone(0211),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have none of the permissions in 0211"),
		},
		{
			ID:       testhelper.MkID("Or - bad - none pass"),
			fileMode: 0666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0121),
				check.FilePermHasNone(0112),
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)"+
				" should have none of the permissions in 0211",
				" should have none of the permissions in 0121",
				" should have none of the permissions in 0112",
				" or "),
		},
		{
			ID:       testhelper.MkID("Or - good - one passes"),
			fileMode: 0666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0121),
				check.FilePermHasNone(0111),
			),
		},
		{
			ID:       testhelper.MkID("And - bad - one fails"),
			fileMode: 0666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0111),
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should have none of the permissions in 0211"),
		},
		{
			ID:       testhelper.MkID("And - good - all pass"),
			fileMode: 0666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(01),
				check.FilePermHasNone(010),
				check.FilePermHasNone(0100),
			),
		},
		{
			ID:       testhelper.MkID("Not - bad - one fails"),
			fileMode: 0666,
			cf: check.FilePermNot(
				check.FilePermHasAll(0200),
				"have owner write permission set",
			),
			ExpErr: testhelper.MkExpErr("the permissions (0666)" +
				" should not have owner write permission set"),
		},
		{
			ID:       testhelper.MkID("Not - good"),
			fileMode: 0666,
			cf: check.FilePermNot(
				check.FilePermHasAll(01),
				"should not have other execute permission set",
			),
		},
	}

	for _, tc := range testCases {
		err := tc.cf(tc.fileMode)
		testhelper.CheckExpErr(t, err, tc)
	}
}
