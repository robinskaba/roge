package conversion

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/robloxapi/rbxfile"
)

var ErrMissingPackageTarget error = fmt.Errorf("directory has neither init.luau nor a .luau file with the same name as the directory")

func GetPackageEntry(directory string) (string, error) {
	entries, err := os.ReadDir(directory)
	if err != nil {
		return "", err
	}

	singleFileTargetName := fmt.Sprintf("%s.luau", filepath.Base(directory))
	for _, e := range entries {
		// match as nested project
		if e.Name() == "init.luau" {
			return ".", nil // returns a directory path !!
		}

		// match as single file project
		if e.Name() == singleFileTargetName {
			return singleFileTargetName, nil
		}
	}
	return "", ErrMissingPackageTarget
}

func readFile(path string) (*rbxfile.Instance, error) {
	src, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading %s: %w", path, err)
	}

	name := strings.TrimSuffix(filepath.Base(path), filepath.Ext(path))
	inst := rbxfile.NewInstance("ModuleScript")
	inst.Properties["Name"] = rbxfile.ValueString(name)
	inst.Properties["Source"] = rbxfile.ValueProtectedString(src)
	return inst, nil
}

func readDir(dir string) (*rbxfile.Instance, error) {
	absDir, err := filepath.Abs(dir)
	if err != nil {
		return nil, fmt.Errorf("resolving path %s: %w", dir, err)
	}

	initPath := filepath.Join(absDir, "init.luau")
	src, err := os.ReadFile(initPath)
	if err != nil {
		return nil, fmt.Errorf("directory %s must contain init.luau: %w", dir, err)
	}

	inst := rbxfile.NewInstance("ModuleScript")
	inst.Properties["Name"] = rbxfile.ValueString(filepath.Base(absDir))
	inst.Properties["Source"] = rbxfile.ValueProtectedString(src)

	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, fmt.Errorf("reading dir %s: %w", dir, err)
	}

	sort.Slice(entries, func(i, j int) bool {
		return entries[i].Name() < entries[j].Name()
	})

	for _, entry := range entries {
		// skip hidden files (.git, .roge)
		if strings.HasPrefix(entry.Name(), ".") {
			continue
		}

		childPath := filepath.Join(dir, entry.Name())

		if entry.IsDir() {
			child, err := readDir(childPath)
			if err != nil {
				return nil, err
			}
			inst.Children = append(inst.Children, child)
			continue
		}

		if entry.Name() == "init.luau" || filepath.Ext(entry.Name()) != ".luau" {
			continue
		}

		child, err := readFile(childPath)
		if err != nil {
			return nil, err
		}
		inst.Children = append(inst.Children, child)
	}

	return inst, nil
}

func writeLeaf(file *rbxfile.Instance, location string) error {
	name := file.Properties["Name"].String()
	filePath := filepath.Join(location, name+".luau")
	return os.WriteFile(filePath, []byte(sourceOf(file)), 0644)
}

func writeParent(file *rbxfile.Instance, location string) error {
	name := file.Properties["Name"].String()
	dir := filepath.Join(location, name)

	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("creating dir %s: %w", dir, err)
	}

	initPath := filepath.Join(dir, "init.luau")
	if err := os.WriteFile(initPath, []byte(sourceOf(file)), 0644); err != nil {
		return fmt.Errorf("writing %s: %w", initPath, err)
	}

	for _, child := range file.Children {
		if err := rbxFileToLuau(child, dir); err != nil {
			return err
		}
	}
	return nil
}
