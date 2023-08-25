package project

import (
	"path/filepath"
	"runtime"
)

var (
	_, b, _, _ = runtime.Caller(0)
)

// Root returns the project absolute path.
func Root() string {
	return filepath.Join(filepath.Dir(b), "../..")
}
