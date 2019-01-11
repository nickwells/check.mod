package check_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nickwells/check.mod/check"
	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestFilePerm(t *testing.T) {
	testCases := []struct {
		name           string
		fileMode       os.FileMode
		cf             check.FilePerm
		errExpected    bool
		errMustContain []string
	}{
		{
			name:     "EQ - good",
			fileMode: 07777,
			cf:       check.FilePermEQ(0777),
		},
		{
			name:        "EQ - bad",
			fileMode:    07777,
			cf:          check.FilePermEQ(0770),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0777) should be equal to 0770",
			},
		},
		{
			name:     "HasAll - good",
			fileMode: 07777,
			cf:       check.FilePermHasAll(0111),
		},
		{
			name:        "HasAll - bad - has none",
			fileMode:    07666,
			cf:          check.FilePermHasAll(0111),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should have all of the permissions in 0111",
			},
		},
		{
			name:        "HasAll - bad - has some",
			fileMode:    07666,
			cf:          check.FilePermHasAll(0221),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should have all of the permissions in 0221",
			},
		},
		{
			name:     "HasNone - good",
			fileMode: 07666,
			cf:       check.FilePermHasNone(0111),
		},
		{
			name:        "HasNone - bad - has all",
			fileMode:    07777,
			cf:          check.FilePermHasNone(0111),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0777)" +
					" should have none of the permissions in 0111",
			},
		},
		{
			name:        "HasNone - bad - has some",
			fileMode:    07666,
			cf:          check.FilePermHasNone(0211),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should have none of the permissions in 0211",
			},
		},
		{
			name:     "Or - bad - none pass",
			fileMode: 07666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0121),
				check.FilePermHasNone(0112),
			),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should have none of the permissions in 0211",
				" should have none of the permissions in 0121",
				" should have none of the permissions in 0112",
				" OR ",
			},
		},
		{
			name:     "Or - good - one passes",
			fileMode: 07666,
			cf: check.FilePermOr(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0121),
				check.FilePermHasNone(0111),
			),
		},
		{
			name:     "And - bad - one fails",
			fileMode: 07666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(0211),
				check.FilePermHasNone(0111),
			),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should have none of the permissions in 0211",
			},
		},
		{
			name:     "And - good - all pass",
			fileMode: 07666,
			cf: check.FilePermAnd(
				check.FilePermHasNone(01),
				check.FilePermHasNone(010),
				check.FilePermHasNone(0100),
			),
		},
		{
			name:     "Not - bad - one fails",
			fileMode: 07666,
			cf: check.FilePermNot(
				check.FilePermHasAll(0200),
				"should not have owner write permission set",
			),
			errExpected: true,
			errMustContain: []string{
				"the permissions (0666)" +
					" should not have owner write permission set",
			},
		},
		{
			name:     "Not - good",
			fileMode: 07666,
			cf: check.FilePermNot(
				check.FilePermHasAll(01),
				"should not have other execute permission set",
			),
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)

		err := tc.cf(tc.fileMode)
		testhelper.CheckError(t, tcID, err, tc.errExpected, tc.errMustContain)
	}
}
