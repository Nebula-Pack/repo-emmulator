package service

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/Nebula-Pack/repo-emmulator/pkg/clone"
	"github.com/google/uuid"
)

type ScanRockspecResponse struct {
	Lua string `json:"lua"`
}

func CloneRepository(repo string, version string) (string, bool, bool, *ScanRockspecResponse, string, error) {
	repoURL := fmt.Sprintf("https://github.com/%s", repo)
	cacheDir := filepath.Join("cache", uuid.New().String())

	fmt.Printf("Cloning repository: %s into %s\n", repoURL, cacheDir)
	err := clone.CloneRepo(repoURL, cacheDir, version)
	if err != nil {
		return repoURL, false, false, nil, "", err
	}

	isLua, hasRockspec, err := clone.CheckProjectFiles(cacheDir)
	fmt.Printf("CheckProjectFiles: isLua=%t, hasRockspec=%t\n", isLua, hasRockspec)
	if err != nil {
		return repoURL, false, false, nil, "", err
	}

	defer func() {
		if err := os.RemoveAll(cacheDir); err != nil {
			fmt.Printf("Failed to delete cache directory %s: %s\n", cacheDir, err.Error())
		}
	}()

	if !isLua {
		return repoURL, false, false, nil, version, nil
	}

	var scanResponse *ScanRockspecResponse
	if hasRockspec {
		var rockspecPath string
		err = filepath.Walk(cacheDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if !info.IsDir() && strings.HasSuffix(info.Name(), ".rockspec") {
				absPath, err := filepath.Abs(path)
				if err != nil {
					return err
				}
				rockspecPath = absPath
				fmt.Printf("Found rockspec file: %s\n", rockspecPath)
				return filepath.SkipDir
			}
			return nil
		})
		if err != nil {
			return repoURL, false, false, nil, "", fmt.Errorf("error finding rockspec file: %w", err)
		}

		if rockspecPath == "" {
			return repoURL, false, false, nil, "", fmt.Errorf("no rockspec file found")
		}

		scanResponse, err = postRockspec(rockspecPath)
		if err != nil {
			return repoURL, false, false, nil, "", err
		}
	}

	if version == "" {
		version, err = clone.ExtractVersion(cacheDir)
		if err != nil {
			return repoURL, false, false, nil, "", err
		}
	}

	fmt.Printf("Final scan response: %+v\n", scanResponse)
	return repoURL, isLua, hasRockspec, scanResponse, version, nil
}

func postRockspec(rockspecPath string) (*ScanRockspecResponse, error) {
	fmt.Printf("Sending absolute rockspec file path: %s\n", rockspecPath)

	reqBody, err := json.Marshal(map[string]string{
		"path": rockspecPath,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request body: %w", err)
	}

	resp, err := http.Post("http://localhost:7777/scan-rockspec", "application/json", bytes.NewBuffer(reqBody))
	if err != nil {
		return nil, fmt.Errorf("failed to send POST request: %w", err)
	}
	defer resp.Body.Close()

	bodyBytes, _ := io.ReadAll(resp.Body)
	fmt.Printf("Response status: %d, body: %s\n", resp.StatusCode, string(bodyBytes))

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to scan rockspec file, status code: %d, response body: %s", resp.StatusCode, string(bodyBytes))
	}

	var scanResponse ScanRockspecResponse
	if err := json.Unmarshal(bodyBytes, &scanResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response body: %w", err)
	}

	fmt.Printf("Decoded scan response: %+v\n", scanResponse)
	return &scanResponse, nil
}
