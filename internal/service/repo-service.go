package service

import (
	"fmt"
	"os"
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

	// Defer deletion of cacheDir after function completes
	defer func() {
		if err := os.RemoveAll(cacheDir); err != nil {
			fmt.Printf("Failed to delete cache directory %s: %s\n", cacheDir, err.Error())
		}
	}()

	return isLua, nil
}
