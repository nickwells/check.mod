package check

import (
	"fmt"
	"os"
	"testing"
)

func TestModeName(t *testing.T) {
	testCases := []struct {
		name   string
		mode   os.FileMode
		expVal string
	}{
		{
			name:   "file",
			mode:   0,
			expVal: "a regular file",
		},
		{
			name:   "dir",
			mode:   os.ModeDir,
			expVal: "a directory",
		},
		{
			name:   "symlink",
			mode:   os.ModeSymlink,
			expVal: "a symlink",
		},
		{
			name:   "named pipe",
			mode:   os.ModeNamedPipe,
			expVal: "a named pipe",
		},
		{
			name:   "socket",
			mode:   os.ModeSocket,
			expVal: "a socket",
		},
		{
			name:   "device",
			mode:   os.ModeDevice,
			expVal: "a device",
		},
		{
			name:   "non-regular",
			mode:   os.ModeIrregular,
			expVal: "a non-regular file",
		},
		{
			name:   "named pipe or socket",
			mode:   os.ModeNamedPipe | os.ModeSocket,
			expVal: "a named pipe or a socket",
		},
	}

	for i, tc := range testCases {
		tcID := fmt.Sprintf("test %d: %s", i, tc.name)
		val := modeName(tc.mode)
		if val != tc.expVal {
			t.Log(tcID)
			t.Log("\t: Expected: " + tc.expVal)
			t.Log("\t:      Got: " + val)
			t.Errorf("\t: unexpected mode name\n")
		}
	}
}
