package typegen

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"

	"github.com/teranos/errors"
)

// CheckResult holds the result of a types check
type CheckResult struct {
	UpToDate    bool
	Differences map[string][]string // language -> files with differences
}

// CompareDirectories compares generated types in tempDir with existing types.
// Returns a CheckResult indicating which files differ.
//
// The subdirectories are the ones getOutputConfig writes to under an
// --output root, not the language names: markdown lands in docs/types, so
// looking for tempDir/markdown found nothing and compareDirectory returned
// no differences for a directory that was never there.
//
//	tempDir/typescript/   -> compares with types/generated/typescript/
//	tempDir/docs/types/   -> compares with docs/types/
func CompareDirectories(tempDir string) (*CheckResult, error) {
	differences := make(map[string][]string)

	// Check TypeScript
	if diffs := compareDirectory(
		filepath.Join(tempDir, "typescript"),
		"types/generated/typescript",
	); len(diffs) > 0 {
		differences["TypeScript"] = diffs
	}

	// Check Markdown
	if diffs := compareDirectory(
		filepath.Join(tempDir, "docs", "types"),
		"docs/types",
	); len(diffs) > 0 {
		differences["Markdown"] = diffs
	}

	return &CheckResult{
		UpToDate:    len(differences) == 0,
		Differences: differences,
	}, nil
}

// compareDirectory compares two directories and returns files with differences.
func compareDirectory(tempDir, existingDir string) []string {
	var diffs []string

	// Check if temp directory exists
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return diffs
	}

	// Walk through temp directory
	filepath.Walk(tempDir, func(tempPath string, info os.FileInfo, err error) error {
		if err != nil {
			return errors.Wrapf(err, "failed to access %s", tempPath)
		}
		if info.IsDir() {
			return nil
		}

		// Skip certain files
		baseName := filepath.Base(tempPath)
		if shouldSkipFile(baseName) {
			return nil
		}

		// Get relative path
		relPath, err := filepath.Rel(tempDir, tempPath)
		if err != nil {
			return errors.Wrapf(err, "failed to get relative path for %s", tempPath)
		}

		// Corresponding existing file
		existingPath := filepath.Join(existingDir, relPath)

		// Check if files differ
		different, err := filesAreDifferent(tempPath, existingPath)
		if err != nil {
			diffs = append(diffs, relPath+" (error: "+err.Error()+")")
		} else if different {
			diffs = append(diffs, relPath)
		}

		return nil
	})

	return diffs
}

// shouldSkipFile returns true if the file should be skipped during comparison.
func shouldSkipFile(basename string) bool {
	skip := []string{
		"README.md",
	}

	for _, s := range skip {
		if basename == s {
			return true
		}
	}

	return false
}

// filesAreDifferent reports whether two files' contents differ.
func filesAreDifferent(file1, file2 string) (bool, error) {
	content1, err := os.ReadFile(file1)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", file1, err)
	}

	content2, err := os.ReadFile(file2)
	if err != nil {
		return false, fmt.Errorf("failed to read %s: %w", file2, err)
	}

	return !bytes.Equal(content1, content2), nil
}
