package utils

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

func CleanUpExtension(directory string, extension string) error {
	// list entries in dir
	entries, err := os.ReadDir(directory)
	if err != nil {
		return err
	}

	// enforce extension format as ".something"
	if !strings.Contains(extension, ".") {
		extension = fmt.Sprintf(".%s", extension)
	}

	for _, e := range entries {
		if e.IsDir() {
			// clean a subdir
			if err = CleanUpExtension(e.Name(), extension); err != nil {
				return err
			}
		} else if strings.Contains(e.Name(), extension) {
			// clean specific file
			file := filepath.Join(directory, e.Name())
			if err = os.Remove(file); err != nil {
				return err
			}
		}
	}

	empty, err := isDirEmpty(directory)
	if err != nil {
		return err
	}

	if empty {
		err = os.Remove(directory)
		if err != nil {
			return err
		}
	}

	return nil
}

// moves roge executable to the target directory
func MoveFileToDir(currentPath string, targetDir string) error {
	fileName := filepath.Base(currentPath)
	movedFilePath := filepath.Join(targetDir, fileName)

	currentPath = filepath.Clean(currentPath)
	movedFilePath = filepath.Clean(movedFilePath)
	if currentPath == movedFilePath {
		return fmt.Errorf("cannot move to the same location")
	}

	// create target dir if not exists
	if err := os.MkdirAll(targetDir, 0755); err != nil {
		return err
	}

	// copy roge executable
	binary, err := os.ReadFile(currentPath)
	if err != nil {
		return err
	}
	if err = os.WriteFile(movedFilePath, binary, 0755); err != nil {
		return err
	}

	return nil
}

func isDirEmpty(directory string) (bool, error) {
	file, err := os.Open(directory)
	if err != nil {
		return false, err
	}
	defer file.Close()

	_, err = file.Readdirnames(1)
	if err == io.EOF {
		return true, nil
	}
	return false, err
}
