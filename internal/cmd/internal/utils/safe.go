package utils

import (
	"errors"

	"github.com/robinskaba/roge/internal/cmd/internal/ux"
	"github.com/robinskaba/roge/internal/config"
	"github.com/robinskaba/roge/internal/repository"
)

func SafeGlobalCfg() config.Config {
	cfg, err := config.LoadConfig()
	if err != nil {
		ux.Fatal("failed to load roge configuration", err)
	}
	return cfg
}

func SafeRepository() *repository.RogeRepo {
	repo, err := repository.Load()
	if err != nil {
		if errors.Is(err, repository.ErrNotRogeRepo) {
			ux.Misuse("not a roge repository (or any of its parent directories): .roge")
		}
		ux.Fatal("failed to load roge repository", err)
	}
	return repo
}
