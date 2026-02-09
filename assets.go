package entydad

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
)

// CopyStyles copies entydad's CSS assets to the target directory.
// Uses runtime.Caller(0) via packageDir() to discover entydad's package
// directory, same approach as pyeza-golang and centymo.
//
// Files are copied to {targetDir}/entydad/ to keep them namespaced.
//
// Example:
//
//	cssDir := filepath.Join("assets", "css")
//	if err := entydad.CopyStyles(cssDir); err != nil {
//	    log.Printf("Warning: Failed to copy entydad styles: %v", err)
//	}
func CopyStyles(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine entydad package directory")
	}

	srcDir := filepath.Join(dir, "assets", "css")
	dstDir := filepath.Join(targetDir, "entydad")

	copied, err := copyDirFiles(srcDir, dstDir, "*.css")
	if err != nil {
		return fmt.Errorf("failed to copy entydad styles: %w", err)
	}

	if copied == 0 {
		log.Printf("entydad: no CSS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d entydad styles to: %s", copied, dstDir)
	return nil
}

// CopyStaticAssets copies entydad's JavaScript assets to the target directory.
// Uses runtime.Caller(0) via packageDir() to discover entydad's package
// directory, same approach as pyeza-golang and centymo.
//
// Files are copied to {targetDir}/entydad/ to keep them namespaced.
//
// Example:
//
//	jsDir := filepath.Join("assets", "js")
//	if err := entydad.CopyStaticAssets(jsDir); err != nil {
//	    log.Printf("Warning: Failed to copy entydad assets: %v", err)
//	}
func CopyStaticAssets(targetDir string) error {
	dir := packageDir()
	if dir == "" {
		return fmt.Errorf("could not determine entydad package directory")
	}

	srcDir := filepath.Join(dir, "assets", "js")
	dstDir := filepath.Join(targetDir, "entydad")

	copied, err := copyDirFiles(srcDir, dstDir, "*.js")
	if err != nil {
		return fmt.Errorf("failed to copy entydad assets: %w", err)
	}

	if copied == 0 {
		log.Printf("entydad: no JS files found in %s", srcDir)
		return nil
	}

	log.Printf("Copied %d entydad assets to: %s", copied, dstDir)
	return nil
}

// copyDirFiles copies all files matching a glob pattern from srcDir to dstDir.
func copyDirFiles(srcDir, dstDir, pattern string) (int, error) {
	if err := os.MkdirAll(dstDir, 0755); err != nil {
		return 0, fmt.Errorf("failed to create target directory: %w", err)
	}

	files, err := filepath.Glob(filepath.Join(srcDir, pattern))
	if err != nil {
		return 0, fmt.Errorf("failed to list source files: %w", err)
	}

	var copied int
	for _, srcFile := range files {
		data, err := os.ReadFile(srcFile)
		if err != nil {
			log.Printf("Warning: Failed to read %s: %v", srcFile, err)
			continue
		}

		dstFile := filepath.Join(dstDir, filepath.Base(srcFile))
		if err := os.WriteFile(dstFile, data, 0644); err != nil {
			return copied, fmt.Errorf("failed to write %s: %w", dstFile, err)
		}
		copied++
	}

	return copied, nil
}
