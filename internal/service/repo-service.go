package service

import (
	"fmt"
	"path/filepath"

	"github.com/Nebula-Pack/repo-emmulator/pkg/clone"
	"github.com/google/uuid"
)

func CloneRepository(repo string) error {
	repoURL := fmt.Sprintf("https://github.com/%s", repo)

	cacheDir := filepath.Join("cache", uuid.New().String())

	err := clone.CloneRepo(repoURL, cacheDir)
	if err != nil {
		return err
	}

	return nil
}
