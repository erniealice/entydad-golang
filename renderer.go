package entydad

import (
	"path/filepath"
	"runtime"
)

// TemplatePatterns returns glob patterns for entydad's templates.
// Uses runtime.Caller(0) to discover entydad's package directory,
// same approach as pyeza-golang.
// Consumer apps merge these patterns with pyeza + app patterns when
// initializing the renderer.
func TemplatePatterns() []string {
	dir := packageDir()
	return []string{
		filepath.Join(dir, "templates", "client", "*.html"),
		filepath.Join(dir, "templates", "user", "*.html"),
		filepath.Join(dir, "templates", "location", "*.html"),
		filepath.Join(dir, "templates", "role", "*.html"),
		filepath.Join(dir, "templates", "permission", "*.html"),
		filepath.Join(dir, "templates", "workspace", "*.html"),
		filepath.Join(dir, "templates", "login", "*.html"),
	}
}

// packageDir returns the absolute directory of this source file.
func packageDir() string {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return ""
	}
	return filepath.Dir(filename)
}
