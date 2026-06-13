package cmd

import (
	"fmt"

	"github.com/robinskaba/roge/internal/cmd/internal/shared"
	"github.com/robinskaba/roge/internal/cmd/internal/utils"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/roblox"
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
	repo := utils.SafeRepository()
	cfg := utils.GetAnyCfg()
	ux.RequireApiKey(cfg)

	if repo.Asset.AssetId == "" {
		ux.Misuse("can not pull without a set asset ID, consider using clone")
	}

	out := cmd.OutOrStdout()
	err := shared.RunDownload(shared.DownloadCfg{
		ApiKey:   cfg.ApiKey,
		AssetId:  repo.Asset.AssetId,
		RepoPath: repo.Path,
		Out:      out,
	})
	if err != nil {
		ux.Fatal("failed to download asset", err)
	}

	asset, err := roblox.GetAsset(cfg.ApiKey, repo.Asset.AssetId)
	if err != nil {
		ux.Fatal("failed to fetch asset", err)
	}
	oldVersion := repo.Asset.Version
	repo.Asset.Version = asset.Version.Id

	if err = repo.Save(); err != nil {
		ux.Fatal("failed to save versioning", err)
	}

	fmt.Fprintf(out, "from rbxasset://%s (%s)\n", repo.Asset.AssetId, asset.Name)
	fmt.Fprintf(out, " * %s --> %s\n", ux.Colored(fmt.Sprintf("%d", asset.Version.Id), ux.Red), ux.Colored(fmt.Sprintf("local %d", oldVersion), ux.Yellow))
}
