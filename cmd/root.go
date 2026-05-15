package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var Version string = "dev"

var rootCmd = &cobra.Command{
	Use:   "roge",
	Short: "CLI for Roblox's native package versioning system",
	Long: `Roge is a command-line interface for managing packages using Roblox's native package versioning system.

It allows you to clone, pull and push Roblox packages, all from your command line.`,
	Example: `  roge config set --api-key YOUR_KEY --author-id 1234567890 --global
  roge clone 1234567890
  roge push`,
	Version: Version,
}

func Execute() {
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.InitDefaultHelpCmd()

	for _, cmd := range rootCmd.Commands() {
		if cmd.Name() == "help" {
			cmd.Short = "help about any command"
			cmd.Long = "help provides a help for any command"
			break
		}
	}

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
