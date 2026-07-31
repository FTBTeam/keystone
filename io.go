package keystone

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func DirectoryExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return stat.IsDir()
}

func FileExists(path string) bool {
	stat, err := os.Stat(path)
	if err != nil {
		return false
	}

	return !stat.IsDir()
}

func CopyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return nil
}

// ParseJsonFile reads a JSON file from the specified path and unmarshal its content into a variable of type T.
func ParseJsonFile[T any](path string) (T, error) {
	var result T

	file, err := os.Open(path)
	if err != nil {
		return result, err
	}
	defer file.Close()

	fileData, err := io.ReadAll(file)
	if err != nil {
		return result, err
	}

	err = json.Unmarshal(fileData, &result)
	if err != nil {
		return result, err
	}

	return result, nil
}

var ErrPathTraversal = errors.New("path traversal detected")

// EnsurePathWithinRoot verifies that path resolves to a location inside root.
func EnsurePathWithinRoot(path, root string) (string, error) {
	absRoot, err := filepath.Abs(filepath.Clean(root))
	if err != nil {
		return "", fmt.Errorf("pathsafe: resolving root %q: %w", root, err)
	}

	absPath, err := filepath.Abs(filepath.Clean(path))
	if err != nil {
		return "", fmt.Errorf("pathsafe: resolving path %q: %w", path, err)
	}

	rel, err := filepath.Rel(absRoot, absPath)
	if err != nil {
		return "", fmt.Errorf("%w: %q cannot be related to root %q: %v", ErrPathTraversal, path, root, err)
	}

	if rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
		return "", fmt.Errorf("%w: %q resolves to %q, which is outside root %q", ErrPathTraversal, path, absPath, absRoot)
	}

	return absPath, nil
}
