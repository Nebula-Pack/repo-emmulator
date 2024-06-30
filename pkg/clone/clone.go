package clone

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
)

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

func CheckProjectFiles(cacheDir string) (bool, bool, error) {
	var luaFileFound, rockspecFileFound bool

	err := filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".lua") {
			luaFileFound = true
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".rockspec") {
			rockspecFileFound = true
		}

		if info.IsDir() && info.Name() == "node_modules" {
			return filepath.SkipDir
		}

		return nil
	})

	if err != nil {
		return false, false, err
	}

	return luaFileFound, rockspecFileFound, nil
}

func ExtractVersion(cacheDir string) (string, error) {
	repo, err := git.PlainOpen(cacheDir)
	if err != nil {
		return "", err
	}

	// Try to get the latest tag
	tags, err := repo.Tags()
	if err != nil {
		return "", err
	}

	var latestTag string
	err = tags.ForEach(func(ref *plumbing.Reference) error {
		latestTag = ref.Name().Short()
		return nil
	})

	if latestTag != "" {
		return latestTag, nil
	}

	// Fall back to using the latest commit hash
	head, err := repo.Head()
	if err != nil {
		return "", err
	}

	return head.Hash().String(), nil
}
