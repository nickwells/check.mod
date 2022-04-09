package check

import (
	"io/fs"
	"time"
)

// These aliases are to simplify the migration to the v2.0.0 version
type (
	Duration      = ValCk[time.Duration]
	FileInfo      = ValCk[fs.FileInfo]
	FilePerm      = ValCk[fs.FileMode]
	Float64       = ValCk[float64]
	Int64         = ValCk[int64]
	Int64Slice    = ValCk[[]int64]
	MapStringBool = ValCk[map[string]bool]
	String        = ValCk[string]
	StringSlice   = ValCk[[]string]
	Time          = ValCk[time.Time]
	TimeLocation  = ValCk[time.Location]
)
