package shared

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/conversion"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/robinskaba/roge/internal/system"
)

type DownloadCfg struct {
	ApiKey   string
	AssetId  string
	RepoPath string
	Version  *int
	Out      io.Writer
}

func RunDownload(downloadCfg DownloadCfg) error {
	fmt.Fprintf(downloadCfg.Out, "downloading asset: %s\n", downloadCfg.AssetId)

	var version int = 0
	if downloadCfg.Version != nil {
		version = *downloadCfg.Version
	}
	rbxFile, err := roblox.Pull(downloadCfg.ApiKey, downloadCfg.AssetId, version)
	if err != nil {
		ux.Fatal("failed to download package from Roblox", err)
	}

	// clean up old luau files
	fmt.Fprintln(downloadCfg.Out, "removing dated luau files..")
	repoDir := filepath.Dir(downloadCfg.RepoPath)
	err = system.CleanUpExtension(repoDir, "luau")
	if err != nil {
		ux.Fatal("failed to clean up old luau files", err)
	}

	fmt.Fprintln(downloadCfg.Out, "unpacking asset files..")
	_, err = conversion.RBXRootToLuau(rbxFile, repoDir)
	if err != nil {
		ux.Fatal("failed to write package as luau files", err)
	}

	return nil
}
