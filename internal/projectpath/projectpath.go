package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)

	// root dir of this project
	Root = filepath.Join(filepath.Dir(b), "../..")
)
