/* Contains the absolute path to the project. */
package projectpath

import (
	"path/filepath"
	"runtime"
)

var (
	_, f, _, _ = runtime.Caller(0)

	// Root folder of this project.
	Root = filepath.Join(filepath.Dir(f), "../../")
)
