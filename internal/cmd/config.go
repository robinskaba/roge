package cmd

import (
	"fmt"
	"os"
	"unicode/utf8"

	"github.com/robinskaba/roge/internal/cmd/internal/utils"
	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/config"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "manage tool configuration",
	Long: `Manage the configuration properties for the tool.
	You can set and list properties either locally (for the current repository)
	or globally (for the user). Local configurations take precedence over global ones.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "set configuration properties",
	Long: `Set configuration properties such as api-key and author-id.
	By default, properties are set locally for the current repository.
	Use the --global flag to set them globally instead.`,
	Example: `  roge config set --api-key YOUR_KEY
		roge config set --author-id 123456 --global
		roge config set --author-id 654321 --group --global`,
	Run: runConfigSet,
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "list configuration properties",
	Long: `List the current configuration properties.
	By default, this lists the local configuration. Use the --global flag
	to list the global configuration.`,
	Example: `  roge config list
		roge config list --global`,
	Run: runConfigList,
}

func init() {
	configSetCmd.Flags().String("api-key", "", "Roblox API key for Assets(read+write), LegacyAssets(manage)")
	configSetCmd.Flags().String("author-id", "", "author ID of the package author (you / your group)")
	configSetCmd.Flags().Bool("group", false, "set a group as a package author")

	configSetCmd.Flags().Bool("global", false, "set to global configuration")
	configSetCmd.Flags().Bool("local", true, "set to the configuration of the repository")

	configListCmd.Flags().Bool("global", false, "list global configuration")
	configListCmd.Flags().Bool("local", true, "list local configuration")

	configCmd.AddCommand(configSetCmd, configListCmd)
	rootCmd.AddCommand(configCmd)
}

func runConfigSet(cmd *cobra.Command, args []string) {
	useLocal := useLocalFromFlags(cmd)
	cfg := getRightConfig(useLocal)

	// load vars
	apiKey, _ := cmd.Flags().GetString("api-key")
	authorId, _ := cmd.Flags().GetString("author-id")
	isAuthorGroup, _ := cmd.Flags().GetBool("group")

	// need variables
	if apiKey == "" && authorId == "" {
		cmd.Help()
		os.Exit(1)
	}

	out := cmd.OutOrStdout()
	if useLocal {
		fmt.Fprintln(out, "writing to local configuration")
	} else {
		fmt.Fprintln(out, "writing to global configuration")
	}

	// set vars
	if apiKey != "" {
		cfg.ApiKey = apiKey

		// limit api key length on print
		formattedApiKey := apiKey
		if utf8.RuneCountInString(apiKey) > 25 {
			formattedApiKey = fmt.Sprintf("%.25s...", apiKey)
		}
		fmt.Fprintf(out, "-> set API key to: %s\n", formattedApiKey)
	}
	if authorId != "" {
		cfg.AuthorId = authorId
		cfg.AuthorType = roblox.CreatorTypeUser
		if isAuthorGroup {
			cfg.AuthorType = roblox.CreatorTypeGroup
		}

		fmt.Fprintf(out, "-> set %s ID of author to: %s\n", cfg.AuthorType, authorId)
	}

	if useLocal {
		repo := utils.SafeRepository() // repository should be cached therefore the same as the one fetched in getRightConfig
		repo.Config = cfg
		if err := repo.Save(); err != nil {
			ux.Fatal("failed to save local configuration", err)
		}
	} else {
		err := config.SaveConfig(cfg)
		if err != nil {
			ux.Fatal("failed to save config", err)
		}
	}
}

func runConfigList(cmd *cobra.Command, args []string) {
	useLocal := useLocalFromFlags(cmd)
	out := cmd.OutOrStdout()
	if useLocal {
		fmt.Fprintln(out, "listing local configuration")
	} else {
		fmt.Fprintln(out, "listing global configuration")
	}

	cfg := getRightConfig(useLocal)

	// list print
	ux.ListStruct(cfg, out)
}

func useLocalFromFlags(cmd *cobra.Command) bool {
	globalFlag, _ := cmd.Flags().GetBool("global")
	localFlag, _ := cmd.Flags().GetBool("local")

	return !globalFlag && localFlag
}

func getRightConfig(useLocal bool) config.Config {
	var cfg config.Config
	if useLocal {
		repo := utils.SafeRepository()
		cfg = repo.Config
	} else {
		cfg = utils.SafeGlobalCfg()
	}
	return cfg
}
