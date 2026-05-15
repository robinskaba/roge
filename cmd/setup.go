package cmd

import (
	"fmt"

	"github.com/robinskaba/roge/internal/pkg"
	"github.com/spf13/cobra"
)

var setupCmd = &cobra.Command{
	Use:   "setup",
	Short: "installs roge and configures PATH",
	Long:  "Moves roge executable to user's program directory and if necessary, adds it to the PATH environmental variable.",
	Run:   runSetup,
}

func init() {
	rootCmd.AddCommand(setupCmd)
}

func runSetup(cmd *cobra.Command, args []string) {
	out := cmd.OutOrStdout()
	fmt.Fprintln(out, "installing roge..")
	installPath, err := pkg.Setup()
	if err != nil {
		fatal("failed to install roge", err)
	}
	fmt.Fprintf(out, "successfully installed roge to %s and added it to path, reopen terminal and verify installation with 'roge --version'\n", installPath)
}
