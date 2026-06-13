package conversion

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/robloxapi/rbxfile"
)

func rbxFileToLuau(file *rbxfile.Instance, location string) error {
	if len(file.Children) == 0 {
		return writeLeaf(file, location)
	}
	return writeParent(file, location)
}

func RBXRootToLuau(file *rbxfile.Instance, dir string) (string, error) {
	if len(file.Children) == 0 {
		if err := writeLeaf(file, dir); err != nil {
			return "", err
		}
		return file.Properties["Name"].String() + ".luau", nil
	}

	if err := os.WriteFile(filepath.Join(dir, "init.luau"), []byte(sourceOf(file)), 0644); err != nil {
		return "", err
	}
	for _, child := range file.Children {
		if err := rbxFileToLuau(child, dir); err != nil {
			return "", err
		}
	}
	return ".", nil
}

func LuauToRBXFile(path string) (*rbxfile.Instance, error) {
	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("stat %s: %w", path, err)
	}

	if info.IsDir() {
		return readDir(path)
	}
	return readFile(path)
}
