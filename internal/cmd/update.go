package cmd

import (
	"context"
	"fmt"

	"github.com/robinskaba/roge/internal/app"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/spf13/cobra"
)

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "update roge to the latest version",
	Long:  "Updates roge to the latest version available.",
	Run:   runUpdate,
}

func init() {
	rootCmd.AddCommand(updateCmd)
}

func runUpdate(cmd *cobra.Command, args []string) {
	ctx := context.Background()
	out := cmd.OutOrStdout()

	fmt.Fprintln(out, "checking for the latest version..")
	updated, latest, err := app.Update(ctx, Version)
	if err != nil {
		ux.Fatal("failed to update roge", err)
	}

	if !updated {
		fmt.Fprintln(out, "roge is already at the latest version")
	} else {
		fmt.Fprintf(out, "roge was updated to %s\n", latest)
	}
}
