package entydad

import (
	"path/filepath"
	"runtime"
)

// packageDir returns the absolute directory of this source file.
func packageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}
