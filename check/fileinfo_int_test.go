package check

import (
	"os"
	"testing"

	"github.com/nickwells/testhelper.mod/testhelper"
)

func TestModeName(t *testing.T) {
	testCases := []struct {
		testhelper.ID
		mode   os.FileMode
		expVal string
	}{
		{
			ID:     testhelper.MkID("file"),
			mode:   0,
			expVal: "a regular file",
		},
		{
			ID:     testhelper.MkID("dir"),
			mode:   os.ModeDir,
			expVal: "a directory",
		},
		{
			ID:     testhelper.MkID("symlink"),
			mode:   os.ModeSymlink,
			expVal: "a symlink",
		},
		{
			ID:     testhelper.MkID("named pipe"),
			mode:   os.ModeNamedPipe,
			expVal: "a named pipe",
		},
		{
			ID:     testhelper.MkID("socket"),
			mode:   os.ModeSocket,
			expVal: "a socket",
		},
		{
			ID:     testhelper.MkID("device"),
			mode:   os.ModeDevice,
			expVal: "a device",
		},
		{
			ID:     testhelper.MkID("non-regular"),
			mode:   os.ModeIrregular,
			expVal: "a non-regular file",
		},
		{
			ID:     testhelper.MkID("named pipe or socket"),
			mode:   os.ModeNamedPipe | os.ModeSocket,
			expVal: "a named pipe or a socket",
		},
	}

	for _, tc := range testCases {
		val := modeName(tc.mode)
		testhelper.DiffString(t, tc.IDStr(), "mode name", val, tc.expVal)
	}
}
