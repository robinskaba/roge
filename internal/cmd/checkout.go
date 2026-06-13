package cmd

import (
	"fmt"
	"strconv"

	"github.com/robinskaba/roge/internal/cmd/internal/shared"
	"github.com/robinskaba/roge/internal/cmd/internal/utils"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/spf13/cobra"
)

var checkoutCmd *cobra.Command = &cobra.Command{
	Use:   "checkout <version>",
	Short: "checkout a specific package version",
	Long: `Checkout a specific version of the package from Roblox.
This command downloads the specified version of the package, converts it to Luau files, and overwrites the local files in the repository.
Use 'roge log' to see available versions and their IDs.`,
	Example: `  roge checkout 5`,
	Args:    cobra.ExactArgs(1),
	Run:     runCheckout,
}

func init() {
	rootCmd.AddCommand(checkoutCmd)
}

func runCheckout(cmd *cobra.Command, args []string) {
	version, err := strconv.Atoi(args[0])
	if err != nil || version == 0 {
		ux.Misuse("version needs to be a positive integer")
	}

	repo := utils.SafeRepository()
	cfg := utils.GetAnyCfg()
	out := cmd.OutOrStdout()

	asset, err := roblox.GetAsset(cfg.ApiKey, repo.Asset.AssetId)
	if err != nil {
		ux.Fatal("failed to fetch asset", err)
	}
	behind := asset.Version.Id - version
	if behind < 0 {
		ux.Misuse("required version of the asset does not exist")
	}

	err = shared.RunDownload(shared.DownloadCfg{
		ApiKey:   cfg.ApiKey,
		AssetId:  repo.Asset.AssetId,
		RepoPath: repo.Path,
		Version:  &version,
		Out:      out,
	})
	if err != nil {
		ux.Fatal("failed to download asset", err)
	}

	repo.Asset.Version = version
	err = repo.Save()
	if err != nil {
		ux.Fatal("failed to save versioning", err)
	}

	var behindStr string
	if behind > 0 {
		behindStr = fmt.Sprintf(" %s", ux.Colored(fmt.Sprintf("[%d behind]", behind), ux.Yellow))
	}

	fmt.Fprintf(out, "loaded asset version %d%s - run 'roge push' to overwrite the latest version\n", version, behindStr)
}
