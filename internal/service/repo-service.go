package service

import (
	"fmt"
	"path/filepath"

	"github.com/Nebula-Pack/repo-emmulator/pkg/clone"
	"github.com/google/uuid"
)

// CloneRepository clones the repository and checks if it's a Lua project
func CloneRepository(repo string) (bool, error) {
	repoURL := fmt.Sprintf("https://github.com/%s", repo)
	cacheDir := filepath.Join("cache", uuid.New().String())

	err := clone.CloneRepo(repoURL, cacheDir)
	if err != nil {
		return false, err
	}

	isLua, err := clone.IsLuaProject(cacheDir)
	if err != nil {
		return false, err
	}

	return isLua, nil
}
