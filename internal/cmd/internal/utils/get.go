package utils

import (
	"github.com/robinskaba/roge/internal/config"
	"github.com/robinskaba/roge/internal/repository"
)

func GetAnyCfg() config.Config {
	cfg := SafeGlobalCfg()

	repo, err := repository.Load()
	if err == nil {
		if repo.Config.ApiKey != "" {
			cfg.ApiKey = repo.Config.ApiKey
		}
		if repo.Config.AuthorId != "" {
			cfg.AuthorId = repo.Config.AuthorId
		}
	}

	return cfg
}
