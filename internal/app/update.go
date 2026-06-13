package app

import (
	"context"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/creativeprojects/go-selfupdate"
)

func Update(ctx context.Context, currentVersion string) (bool, string, error) {
	latest, found, err := selfupdate.DetectLatest(ctx, selfupdate.ParseSlug("robinskaba/roge"))
	if err != nil {
		return false, "", err
	}

	if !found {
		return false, "", fmt.Errorf("latest version for %s/%s could not be found from github repository", runtime.GOOS, runtime.GOARCH)
	}

	if currentVersion != "dev" && latest.LessOrEqual(currentVersion) {
		return false, currentVersion, nil
	}

	exe, err := selfupdate.ExecutablePath()
	if err != nil {
		return false, "", errors.New("could not locate executable path")
	}
	err = selfupdate.UpdateTo(ctx, latest.AssetURL, latest.AssetName, exe)
	if err != nil {
		return false, "", fmt.Errorf("error occured while updating binary: %w", err)
	}

	return true, latest.Version(), nil
}

func TryRemovingOldVersion() {
	// if in executable dir is a same file with .old suffix -> try removing it
	// errors ignored since we do not care if it fails
	self, err := os.Executable()
	if err == nil {
		dir := filepath.Dir(self)
		target := fmt.Sprintf(".%s.old", filepath.Base(self))
		oldFile := filepath.Join(dir, target)
		if _, err := os.Stat(oldFile); err == nil {
			os.Remove(oldFile)
		}
	}
}
