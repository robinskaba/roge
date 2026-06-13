//go:build windows

package app

import (
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"unsafe"

	"github.com/robinskaba/roge/internal/system"
	"golang.org/x/sys/windows/registry"
)

func Setup() (string, error) {
	installDir, err := windowsProgramDir()
	if err != nil {
		return "", err
	}
	return installDir, installWindows(installDir)
}

func windowsProgramDir() (string, error) {
	// determine program directory
	usrCfgDir, err := os.UserConfigDir()
	if err != nil {
		return "", err
	}
	installDir := filepath.Join(usrCfgDir, "roge", "bin")
	exePath, err := os.Executable()
	if err != nil {
		return "", err
	}
	if err := system.MoveFileToDir(exePath, installDir); err != nil {
		return "", err
	}
	return installDir, nil
}

// performs roge installation for Windows
func installWindows(installDir string) error {
	// open key registry
	k, err := registry.OpenKey(registry.CURRENT_USER, `Environment`, registry.QUERY_VALUE|registry.SET_VALUE)
	if err != nil {
		return err
	}
	defer k.Close()

	// read path variable
	pathVar, _, err := k.GetStringValue("Path")
	if err != nil && err != registry.ErrNotExist {
		return err
	}

	// check if already added to path
	installDirLow := strings.ToLower(installDir)
	for _, path := range strings.Split(pathVar, ";") {
		if strings.ToLower(path) == installDirLow {
			return nil
		}
	}

	// add install dir to path (even if pathVar is empty)
	if pathVar != "" && !strings.HasSuffix(pathVar, ";") {
		pathVar += ";"
	}
	pathVar += installDir

	// set new path variable
	if err = k.SetStringValue("Path", pathVar); err != nil {
		return err
	}

	// notify windows of path change
	user32 := syscall.NewLazyDLL("user32.dll")
	sendMessage := user32.NewProc("SendMessageTimeoutW")
	envPtr, _ := syscall.UTF16PtrFromString("Environment")
	sendMessage.Call(
		0xffff, // HWND_BROADCAST
		0x001a, // WM_SETTINGCHANGE
		0,
		uintptr(unsafe.Pointer(envPtr)),
		0x0002, // SMTO_ABORTIFHUNG
		5000,
		0,
	)

	return nil
}
