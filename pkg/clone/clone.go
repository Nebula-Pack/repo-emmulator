package clone

import (
	"os"
	"path/filepath"

	"github.com/go-git/go-git/v5"
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
