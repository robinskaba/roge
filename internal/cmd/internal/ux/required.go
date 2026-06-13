package ux

import "github.com/robinskaba/roge/internal/config"

func RequireApiKey(cfg config.Config) {
	if cfg.ApiKey == "" {
		Misuse("Roblox Assets (read+write) and Legacy-Assets (manage) API key is not set; set it with 'roge config set --api-key <API_KEY> --global'")
	}
}

func RequireAuthorId(cfg config.Config) {
	if cfg.AuthorId == "" {
		Misuse("author ID is not set; set it with 'roge config set --author-id <AUTHOR_ID> --global'")
	}
}
