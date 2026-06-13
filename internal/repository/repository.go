package repository

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/robinskaba/roge/internal/config"
)

type RogeRepo struct {
	Path   string
	Asset  AssetVersioning
	Config config.Config
}

type AssetVersioning struct {
	AssetId string `json:"asset_id"`
	Version int    `json:"version"`
}

var cachedRepo *RogeRepo

func Initialize(dir string) (*RogeRepo, error) {
	var repo RogeRepo
	path := filepath.Join(dir, ".roge")
	err := os.Mkdir(path, 0755)
	if err != nil {
		return nil, err
	}

	repo = RogeRepo{
		Path:  path,
		Asset: AssetVersioning{},
	}

	repo.Save()

	return &repo, nil
}

func IsInitialized() bool {
	_, err := findRepoLocation()
	return err == nil
}

func Load() (*RogeRepo, error) {
	if cachedRepo != nil {
		return cachedRepo, nil
	}

	var repo RogeRepo

	// get repo
	path, err := findRepoLocation()
	if err != nil {
		return nil, err
	}
	repo.Path = path

	// load versioning
	var versioning AssetVersioning
	file, err := os.ReadFile(filepath.Join(repo.Path, "asset.json"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &versioning)
	if err != nil {
		return nil, err
	}
	repo.Asset = versioning

	// load config
	var config config.Config
	file, err = os.ReadFile(filepath.Join(repo.Path, "config.json"))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(file, &config)
	repo.Config = config

	cachedRepo = &repo
	return &repo, err
}

func (repo *RogeRepo) Save() error {
	file, err := json.MarshalIndent(repo.Asset, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(repo.Path, "asset.json"), file, 0600)
	if err != nil {
		return err
	}

	file, err = json.MarshalIndent(repo.Config, "", "  ")
	if err != nil {
		return err
	}
	err = os.WriteFile(filepath.Join(repo.Path, "config.json"), file, 0600)
	if err != nil {
		return err
	}

	return nil
}
