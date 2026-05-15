package cmd

import (
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"

	"github.com/robinskaba/roge/internal/config"
	"github.com/robinskaba/roge/internal/repository"
)

func fatal(msg string, err error) {
	fmt.Fprintf(os.Stderr, "%s: %s: %v\n", colored("fatal", Red), msg, err)
	os.Exit(1)
}

func misuse(msg string) {
	fmt.Fprintf(os.Stderr, "%s: %s\n", colored("misuse", Yellow), msg)
	os.Exit(1)
}

func requireApiKey(cfg config.Config) {
	if cfg.ApiKey == "" {
		misuse("Roblox Assets (read+write) and Legacy-Assets (manage) API key is not set; set it with 'roge config set --api-key <API_KEY> --global'")
	}
}

func requireUserId(cfg config.Config) {
	if cfg.UserId == "" {
		misuse("user ID is not set; set it with 'roge config set --user-id <USER_ID> --global'")
	}
}

func getAnyCfg() config.Config {
	cfg := safeGlobalCfg()

	repo, err := repository.Load()
	if err == nil {
		if repo.Config.ApiKey != "" {
			cfg.ApiKey = repo.Config.ApiKey
		}
		if repo.Config.UserId != "" {
			cfg.UserId = repo.Config.UserId
		}
	}

	return cfg
}

func safeGlobalCfg() config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		fatal("failed to load roge configuration", err)
	}
	return cfg
}

func safeRepository() *repository.RogeRepo {
	repo, err := repository.Load()
	if err != nil {
		if errors.Is(err, repository.ErrNotRogeRepo) {
			misuse("not a roge repository (or any of its parent directories): .roge")
		}
		fatal("failed to load roge repository", err)
	}
	return repo
}

func listStruct(str any, out io.Writer) {
	val := reflect.ValueOf(str)
	typ := val.Type()

	for i := 0; i < val.NumField(); i++ {
		fieldName := typ.Field(i).Tag.Get("json")
		fieldValue := val.Field(i).Interface()
		fmt.Fprintf(out, "- %s=%v\n", fieldName, fieldValue)
	}
}

type Color string

const (
	Reset  Color = "\033[0m"
	Red    Color = "\033[31m"
	Green  Color = "\033[32m"
	Yellow Color = "\033[33m"
	Cyan   Color = "\033[36m"
)

func colored(text any, color Color) string {
	return fmt.Sprintf("%s%v%s", color, text, Reset)
}
