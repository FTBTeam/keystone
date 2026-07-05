package keystone

import (
	"encoding/json"
	"io"
	"os"
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
