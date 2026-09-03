package typegen

import (
	"os"
	"path/filepath"
	"testing"
)

func TestCompareDirectories_Identical(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := t.TempDir()

	content := `export interface User { id: string }
`

	os.WriteFile(filepath.Join(tempDir, "user.ts"), []byte(content), 0644)
	os.WriteFile(filepath.Join(existingDir, "user.ts"), []byte(content), 0644)

	diffs := compareDirectory(tempDir, existingDir)

	if len(diffs) != 0 {
		t.Errorf("Expected no differences for identical files, got: %v", diffs)
	}
}

func TestCompareDirectories_FunctionalChange(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := t.TempDir()

	os.WriteFile(filepath.Join(tempDir, "user.ts"),
		[]byte("export interface User { id: string }\n"), 0644)
	os.WriteFile(filepath.Join(existingDir, "user.ts"),
		[]byte("export interface User { id: number }\n"), 0644)

	diffs := compareDirectory(tempDir, existingDir)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 difference for functional change, got: %v", diffs)
	}
}

// A generated file with no counterpart on disk is a difference, not a silent
// pass: that is how a package dropped from the language lists gets caught.
func TestCompareDirectories_MissingExistingFile(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := t.TempDir()

	os.WriteFile(filepath.Join(tempDir, "user.ts"),
		[]byte("export interface User { id: string }\n"), 0644)

	diffs := compareDirectory(tempDir, existingDir)

	if len(diffs) != 1 {
		t.Errorf("Expected 1 difference for missing existing file, got: %v", diffs)
	}
}

func TestCompareDirectories_SkipsReadme(t *testing.T) {
	tempDir := t.TempDir()
	existingDir := t.TempDir()

	os.WriteFile(filepath.Join(tempDir, "README.md"), []byte("generated\n"), 0644)
	os.WriteFile(filepath.Join(existingDir, "README.md"), []byte("different\n"), 0644)

	diffs := compareDirectory(tempDir, existingDir)

	if len(diffs) != 0 {
		t.Errorf("Expected README.md to be skipped, got: %v", diffs)
	}
}
