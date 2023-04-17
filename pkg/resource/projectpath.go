package resource

import (
	"path/filepath"
	"runtime"
)

var (
	_, f, _, _ = runtime.Caller(0)

	// Path to root folder of this project.
	ProjectRoot = filepath.Join(filepath.Dir(f), "../../")
)
