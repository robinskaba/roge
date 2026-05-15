package cmd

import (
	"fmt"
	"os"

	"github.com/robinskaba/roge/internal/repository"
	"github.com/spf13/cobra"
)

var assetCmd = &cobra.Command{
	Use:   "asset",
	Short: "manage asset configuration",
	Long:  `Manage the configuration and tracking details of the Roblox asset associated with this repository.`,
}

var assetSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set an asset property",
	Long: `Set specific properties for the managed asset.
This allows you to link the local repository to a specific Roblox Asset ID.`,
	Example: `  roge asset set --id 1234567890`,
	Run:     runAssetSet,
}

var assetResetcmd = &cobra.Command{
	Use:   "reset",
	Short: "reset asset configuration",
	Long: `Reset the tracking configuration for the managed asset.
This clears the stored Asset ID and other versioning details, effectively unlinking the local repository from the remote package.`,
	Example: `  roge asset reset`,
	Run:     runAssetReset,
}

var assetViewCmd = &cobra.Command{
	Use:   "view",
	Short: "view details about the managed asset",
	Long: `Display the current tracking configuration for the managed asset.
This outputs the stored properties, such as the associated Asset ID and local version state.`,
	Example: `  roge asset view`,
	Run:     runAssetList,
}

func init() {
	assetSetCmd.Flags().String("id", "", "asset ID of the package managed in this repository")

	assetCmd.AddCommand(assetViewCmd, assetSetCmd, assetResetcmd)
	rootCmd.AddCommand(assetCmd)
}

func runAssetSet(cmd *cobra.Command, args []string) {
	repo := safeRepository()
	assetId, _ := cmd.Flags().GetString("id")
	if assetId == "" {
		cmd.Help()
		os.Exit(1)
	}

	out := cmd.OutOrStdout()
	if assetId != "" {
		repo.Asset.AssetId = assetId
		fmt.Fprintf(out, "set asset ID to %s\n", assetId)
	}

	err := repo.Save()
	if err != nil {
		fatal("failed to save repository", err)
	}
}

func runAssetList(cmd *cobra.Command, args []string) {
	repo := safeRepository()
	out := cmd.OutOrStdout()
	listStruct(repo.Asset, out)
}

func runAssetReset(cmd *cobra.Command, args []string) {
	repo := safeRepository()
	repo.Asset = repository.AssetVersioning{}
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "asset configuration was reset")
	err := repo.Save()
	if err != nil {
		fatal("failed to save repository", err)
	}
}
