package clone

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
)

// CloneRepo clones the given repository to the cacheDir
func CloneRepo(repoURL, cacheDir string) error {
	err := os.MkdirAll(filepath.Dir(cacheDir), os.ModePerm)
	if err != nil {
		return err
	}

	_, err = git.PlainClone(cacheDir, false, &git.CloneOptions{
		URL: repoURL,
	})
	if err != nil {
		return err
	}

	return nil
}

// IsLuaProject checks if the repository is a Lua project
func IsLuaProject(cacheDir string) (bool, error) {
	var luaFileFound bool
	err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".lua") {
			luaFileFound = true
		}

		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return false, err
	}

	return luaFileFound, nil
}
