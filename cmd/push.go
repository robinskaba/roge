package cmd

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/robinskaba/roge/internal/conversion"
	"github.com/robinskaba/roge/internal/roblox"
	"github.com/spf13/cobra"
)

var pushCmd = &cobra.Command{
	Use:   "push",
	Short: "publish a new package version of a module",
	Long: `Publish a new package or update an existing package version on Roblox.
If the package has not been uploaded yet, a new asset is created using your configured user ID.
The entry point of the package is deduced automatically by looking for an init.luau file or a .luau file matching the directory name.
Use flags to override the package name and description.`,
	Example: `  roge push
  roge push --name "My Package" --description "Helpful tools"`,
	Args: cobra.NoArgs,
	Run:  runPush,
}

func init() {
	pushCmd.Flags().String("name", "", "new name of the published package")
	pushCmd.Flags().String("description", "", "new description of the published package")
	rootCmd.AddCommand(pushCmd)
}

func runPush(cmd *cobra.Command, args []string) {
	repo := safeRepository()
	cfg := getAnyCfg()
	requireApiKey(cfg)
	requireUserId(cfg)

	isNewUpload := repo.Asset.AssetId == ""
	out := cmd.OutOrStdout()

	// deduce filepath
	projectDir := filepath.Dir(repo.Path)
	pkgEntry, err := conversion.GetPackageEntry(projectDir)
	if err != nil {
		if errors.Is(err, conversion.ErrMissingPackageTarget) {
			misuse("directory must contain either .luau file matching the name of the directory or an init.luau file")
		}
		fatal("failed to deduce target file", err)
	}

	// prepare file
	fmt.Fprintln(out, "packaging luau files..")
	file, err := conversion.LuauToRBXFile(pkgEntry)
	if err != nil {
		fatal("failed to convert file to an rbx instance", err)
	}
	rbxm, err := conversion.BuildRbxm(file)
	if err != nil {
		fatal("failed to convert an instance to rbxm format", err)
	}

	// set general data
	fmt.Fprintln(out, "configuring package details..")
	var assetId string
	if !isNewUpload {
		assetId = repo.Asset.AssetId
	}
	var authorId string
	if isNewUpload {
		authorId = cfg.UserId
	}

	// set package metadata
	var pkgName string
	var pkgDescription string
	fName, _ := cmd.Flags().GetString("name")
	fDesc, _ := cmd.Flags().GetString("description")
	if isNewUpload {
		pkgName = filepath.Base(projectDir) // default package name is name of the project directory
	}
	if fName != "" {
		pkgName = fName
	}
	if fDesc != "" {
		pkgDescription = fDesc
	}

	// publish
	fmt.Fprintln(out, "pushing to Roblox..")
	assetId, version, err := roblox.Push(roblox.PushConfig{
		ApiKey:      cfg.ApiKey,
		Rbxm:        rbxm,
		AssetId:     assetId,
		AuthorId:    authorId,
		AuthorType:  "USER", // for now always USER, doesnt affect NEW vs UPDATE
		Name:        pkgName,
		Description: pkgDescription,
	})
	if err != nil {
		fatal("failed to push to Roblox", err)
	}

	// load response to versioning
	oldVersion := repo.Asset.Version
	repo.Asset.AssetId = assetId
	repo.Asset.Version = version

	if err = repo.Save(); err != nil {
		fatal("failed to save repository, but the package was published", err)
	}

	if isNewUpload {
		fmt.Fprintf(out, "published to a new package (rbxasset://%s) --> %s\n", repo.Asset.AssetId, pkgName)
	} else {
		localTxt := colored("local", Green)
		remoteTxt := colored("remote", Red)
		fmt.Fprintf(
			out,
			"pushed a new version: %d --> %d %s, %s\n",
			oldVersion, version, localTxt, remoteTxt,
		)
	}
}
