package utils

import (
	"encoding/json"
	"errors"
	"io/fs"
	"os"
	"path"
)

func HasFile(path string) bool {
	_, err := os.Stat(path)
	return !errors.Is(err, os.ErrNotExist)
}

// auto create file
func AutoWriteFile(file_path string, data []byte, mode fs.FileMode) error {
	os.MkdirAll(path.Dir(file_path), fs.ModePerm)
	return os.WriteFile(file_path, data, mode)
}

// read json file
func ReadJsonFile(path string, data any) error {
	file, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(file, data)
}

// write json file
func WriteJsonFile(path string, data any) error {
	file, err := json.Marshal(data)
	if err != nil {
		return err
	}
	return os.WriteFile(path, file, os.ModePerm)
}
