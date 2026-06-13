package cmd

import (
	"fmt"
	"path/filepath"

	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/repository"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "initialize a roge repository in the current directory",
	Long: `Initialize a new roge repository in the current working directory.
This sets up the necessary tracking and configuration files to start managing the package.`,
	Example: `  roge init`,
	Run:     runInit,
}

func init() {
	rootCmd.AddCommand(initCmd)
}

func runInit(cmd *cobra.Command, args []string) {
	_, err := repository.Initialize(".")
	if err != nil {
		ux.Fatal("failed to initialize roge repository", err)
	}
	absPath, err := filepath.Abs(".")
	if err != nil {
		ux.Fatal("failed to get absolute path of repository", err)
	}
	fmt.Fprintf(cmd.OutOrStdout(), "initialized empty roge repository in %s\n", absPath)
}
