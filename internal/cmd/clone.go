package cmd

import (
	"fmt"
	"os"

	"github.com/robinskaba/roge/internal/cmd/internal/utils"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/conversion"
	"github.com/robinskaba/roge/internal/repository"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/spf13/cobra"
)

var cloneCmd = &cobra.Command{
	Use:   "clone <asset-id>",
	Short: "clone a package from Roblox",
	Long: `Clone a package from Roblox using its asset ID.
This command downloads the Roblox file, converts it to Luau files, and initializes a new local repository inside a directory named after the package.
An API key must be configured in your environment to perform this action.`,
	Example: `  roge clone 1234567890`,
	Args:    cobra.ExactArgs(1),
	Run:     runClone,
}

func init() {
	rootCmd.AddCommand(cloneCmd)
}

func runClone(cmd *cobra.Command, args []string) {
	cfg := utils.GetAnyCfg()
	ux.RequireApiKey(cfg)

	assetId := args[0]

	out := cmd.OutOrStdout()
	fmt.Fprintf(out, "downloading asset: %s\n", assetId)
	rbxFile, err := roblox.Pull(cfg.ApiKey, assetId, 0)
	if err != nil {
		ux.Fatal("failed to download package from Roblox", err)
	}

	fmt.Fprintln(out, "unpacking asset files..")
	name := rbxFile.Properties["Name"].String()
	outDir := name
	if err = os.MkdirAll(outDir, 0755); err != nil {
		ux.Fatal("failed to create output directory", err)
	}
	_, err = conversion.RBXRootToLuau(rbxFile, outDir)
	if err != nil {
		ux.Fatal("failed to write package as luau files", err)
	}

	fmt.Fprintln(out, "intializing roge repository..")
	asset, err := roblox.GetAsset(cfg.ApiKey, assetId)
	if err != nil {
		ux.Fatal("failed to fetch asset", err)
	}
	repo, err := repository.Initialize(outDir)
	if err != nil {
		ux.Fatal("failed to initialize roge repository", err)
	}
	repo.Asset.AssetId = assetId
	repo.Asset.Version = asset.Version.Id

	// prepopulate local configuration with asset metadata
	repo.Config.ApiKey = cfg.ApiKey
	repo.Config.AuthorId = asset.Creator.Id
	repo.Config.AuthorType = asset.Creator.Type

	if err = repo.Save(); err != nil {
		ux.Fatal("failed to save repository", err)
	}

	fmt.Fprintf(out, "cloned rbxasset://%s to %s\n", assetId, outDir)
}
