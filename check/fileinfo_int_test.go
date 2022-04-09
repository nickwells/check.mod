package check

import (
	"io/fs"
	"testing"

	"github.com/nickwells/testhelper.mod/v2/testhelper"
)

func TestModeName(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		mode   fs.FileMode
		expVal string
	}{
		{
			ID:     testhelper.MkID("file"),
			mode:   0,
			expVal: "a regular file",
		},
		{
			ID:     testhelper.MkID("dir"),
			mode:   fs.ModeDir,
			expVal: "a directory",
		},
		{
			ID:     testhelper.MkID("symlink"),
			mode:   fs.ModeSymlink,
			expVal: "a symlink",
		},
		{
			ID:     testhelper.MkID("named pipe"),
			mode:   fs.ModeNamedPipe,
			expVal: "a named pipe",
		},
		{
			ID:     testhelper.MkID("socket"),
			mode:   fs.ModeSocket,
			expVal: "a socket",
		},
		{
			ID:     testhelper.MkID("device"),
			mode:   fs.ModeDevice,
			expVal: "a device",
		},
		{
			ID:     testhelper.MkID("non-regular"),
			mode:   fs.ModeIrregular,
			expVal: "a non-regular file",
		},
		{
			ID:     testhelper.MkID("named pipe or socket"),
			mode:   fs.ModeNamedPipe | fs.ModeSocket,
			expVal: "a named pipe or a socket",
		},
	}

	for _, tc := range testCases {
		val := modeName(tc.mode)
		testhelper.DiffString(t, tc.IDStr(), "mode name", val, tc.expVal)
	}
}
