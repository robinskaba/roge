//go:build linux || darwin

package pkg

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/robinskaba/roge/internal/utils"
)

func Setup() (string, error) {
	installDir, err := unixProgramDir()
	if err != nil {
		return "", err
	}
	return installDir, installUnix(installDir)
}

func unixProgramDir() (string, error) {
	// determine program directory
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	installDir := filepath.Join(home, ".local", "bin")
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	if err = utils.MoveFileToDir(exePath, installDir); err != nil {
		return "", err
	}
	return installDir, nil
}

// performs roge installation for UNIX systems
func installUnix(installDir string) error {
	// add to path
	if strings.Contains(os.Getenv("PATH"), installDir) {
		return nil // already in path
	}

	shellConfigs := []string{".zshrc", ".bashrc", ".profile"}
	pathLine := "\n# roge cli\nexport PATH=\"$HOME/.local/bin:$PATH\"\n"

	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	for _, configName := range shellConfigs {
		configPath := filepath.Join(home, configName)
		if _, err := os.Stat(configPath); err != nil {
			continue // config file does not exist -> skipping
		}

		content, err := os.ReadFile(configPath)
		if err != nil {
			continue // failed to read file -> skipping
		}

		if strings.Contains(string(content), ".local/bin") {
			continue // already has unix program dir in path
		}

		f, err := os.OpenFile(configPath, os.O_APPEND|os.O_WRONLY, 0644)
		if err != nil {
			continue // failed to write to config -> skipping
		}
		_, err = f.WriteString(pathLine)
		f.Close()
		if err != nil {
			return err
		}
	}

	return nil
}
