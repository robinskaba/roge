package repository

import (
	"errors"
	"os"
	"path/filepath"
)

var ErrNotRogeRepo = errors.New("not a roge repository (or any of the parent directories): .roge")

func findRepoLocation() (string, error) {
	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		rogePath := filepath.Join(dir, ".roge")
		info, err := os.Stat(rogePath)
		if err == nil && info.IsDir() {
			return rogePath, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			return "", ErrNotRogeRepo
		}
		dir = parent
	}
}
