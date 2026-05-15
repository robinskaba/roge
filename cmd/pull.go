package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/robinskaba/roge/internal/conversion"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/robinskaba/roge/internal/utils"
	"github.com/spf13/cobra"
)

var pullCmd = &cobra.Command{
	Use:   "pull",
	Short: "overwrite local files with a version from Roblox",
	Long: `Pull the latest version of the tracked package from Roblox and overwrite the local Luau files.
This command updates your current working directory with the remote state of the package.
It requires an initialized repository with a tracked Asset ID. If you are downloading a package for the first time, use the 'clone' command instead.`,
	Example: `  roge pull`,
	Run:     runPull,
}

func init() {
	rootCmd.AddCommand(pullCmd)
}

func runPull(cmd *cobra.Command, args []string) {
	repo := safeRepository()
	cfg := getAnyCfg()
	requireApiKey(cfg)

	if repo.Asset.AssetId == "" {
		misuse("can not pull without a set asset ID, consider using clone")
	}

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "downloading asset: %s\n", repo.Asset.AssetId)
	rbxFile, err := roblox.Pull(cfg.ApiKey, repo.Asset.AssetId)
	if err != nil {
		fatal("failed to download package from Roblox", err)
	}

	// clean up old luau files
	fmt.Fprintln(out, "removing dated luau files..")
	repoDir := filepath.Dir(repo.Path)
	err = utils.CleanUpExtension(repoDir, "luau")
	if err != nil {
		fatal("failed to clean up old luau files", err)
	}

	fmt.Fprintln(out, "unpacking asset files..")
	_, err = conversion.RBXRootToLuau(rbxFile, repoDir)
	if err != nil {
		fatal("failed to write package as luau files", err)
	}

	asset, err := roblox.GetAsset(cfg.ApiKey, repo.Asset.AssetId)
	if err != nil {
		fatal("failed to fetch asset", err)
	}
	oldVersion := repo.Asset.Version
	repo.Asset.Version = asset.Version.Id

	if err = repo.Save(); err != nil {
		fatal("failed to save versioning", err)
	}

	fmt.Fprintf(out, "from rbxasset://%s (%s)\n", repo.Asset.AssetId, asset.Name)
	fmt.Fprintf(out, " * %s --> %s\n", colored(fmt.Sprintf("%d", asset.Version.Id), Red), colored(fmt.Sprintf("local %d", oldVersion), Yellow))
}
