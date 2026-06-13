package cmd

import (
	"fmt"
	"os"
	"strconv"

	"github.com/robinskaba/roge/internal/cmd/internal/utils"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/spf13/cobra"
)

var logCmd = &cobra.Command{
	Use:   "log",
	Short: "list versions of the package",
	Long: `List the version history of the current package from Roblox.
This command fetches all available versions for the package tracked in the current repository.
You must run this inside an initialized roge repository.`,
	Example: `  roge log`,
	Run:     runLog,
}

func init() {
	rootCmd.AddCommand(logCmd)
}

func runLog(cmd *cobra.Command, args []string) {
	cfg := utils.GetAnyCfg()
	ux.RequireApiKey(cfg)
	repo := utils.SafeRepository()

	versions, err := roblox.GetVersions(cfg.ApiKey, repo.Asset.AssetId)
	if err != nil {
		ux.Fatal("failed to retrieve package versions", err)
	}

	out := cmd.OutOrStdout()
	if len(versions) < 1 {
		fmt.Fprintf(out, "no package versions exist")
		os.Exit(0)
	}

	idPaddingWidth := len(strconv.Itoa(versions[0].Id)) // count digits
	for _, v := range versions {
		paddedId := fmt.Sprintf("%0*v", idPaddingWidth, v.Id)
		date := v.Time.Format("2006-01-02")
		clock := v.Time.Format("15:04")

		// row formatting
		logRow := fmt.Sprintf("%s %-10s %-5s", ux.Colored(paddedId, ux.Yellow), date, clock)

		// highlight versions
		if v.Id == repo.Asset.Version {
			logRow += " " + ux.Colored("local", ux.Green)
		}
		if v == versions[0] {
			prefix := " "
			if v.Id == repo.Asset.Version {
				prefix = ", "
			}
			logRow += prefix + ux.Colored("remote", ux.Red)
		}

		fmt.Fprintln(out, logRow)
	}
}
