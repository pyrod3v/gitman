// Copyright 2025 pyrod3v
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package gitman

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/charmbracelet/huh"
	"github.com/spf13/viper"
)

var licenses []string
var licensesMutex sync.Mutex

func AddLicense(path string) error {
	var selected string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select a license").
				Options(huh.NewOptions(licenses...)...).
				Value(&selected).
				Height(min(10, len(licenses))),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form failed: %v", err)
	}

	if selected == "" {
		return nil
	}

	licensePath := filepath.Join(GetConfigDir(), "licenses", selected)
	if content, err := os.ReadFile(licensePath); err == nil {
		return os.WriteFile(filepath.Join(path, "LICENSE"), content, 0644)
	}

	licensePath = filepath.Join(GetConfigDir(), ".cache", "licenses", selected)
	if content, err := os.ReadFile(licensePath); err == nil {
		return os.WriteFile(filepath.Join(path, "LICENSE"), content, 0644)
	}

	url := fmt.Sprintf("https://api.github.com/licenses/%s", selected)
	res, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to fetch license: %v", err)
	}
	defer res.Body.Close()

	var licenseData struct {
		Body string `json:"body"`
	}
	if err := json.NewDecoder(res.Body).Decode(&licenseData); err != nil {
		return fmt.Errorf("failed to parse license JSON: %v", err)
	}

	content := []byte(licenseData.Body)
	if viper.GetBool("CacheLicenses") {
		dir := filepath.Join(GetConfigDir(), ".cache", "licenses")
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create cache directory: %v", err)
		}
		cacheFile := filepath.Join(dir, selected)
		if err := os.WriteFile(cacheFile, content, 0644); err != nil {
			fmt.Printf("Failed to cache license: %v\n", err)
		}
	}

	err = os.WriteFile(filepath.Join(path, "LICENSE"), content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write LICENSE file: %v", err)
	}

	return nil
}

func LoadLicenses() {
	if customLicenses, err := loadCustomLicenses(); err == nil {
		licensesMutex.Lock()
		licenses = append(customLicenses, licenses...)
		licensesMutex.Unlock()
	}
	go func() {
		fetchedLicenses, err := fetchLicenses()
		if err != nil {
			fmt.Printf("Failed to fetch license templates: %v\n", err)
		}
		licensesMutex.Lock()
		licenses = append(licenses, strings.Fields(strings.Join(fetchedLicenses, "\n"))...)
		sort.Strings(licenses)
		licensesMutex.Unlock()
	}()
}

func loadCustomLicenses() ([]string, error) {
	gitignoreDir := filepath.Join(GetConfigDir(), "licenses")
	entries, err := os.ReadDir(gitignoreDir)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			err := os.MkdirAll(gitignoreDir, 0755)
			if err != nil {
				err = fmt.Errorf("failed to create directory: %v", err)
			}
			return []string{}, err
		}
		return nil, err
	}
	var customGitignores []string
	for _, entry := range entries {
		if !entry.IsDir() {
			customGitignores = append(customGitignores, entry.Name())
		}
	}
	return customGitignores, nil
}

func fetchLicenses() ([]string, error) {
	res, err := http.Get("https://api.github.com/licenses")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch license list: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	var licenseList []struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal(body, &licenseList); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %v", err)
	}

	var licenses []string
	for _, license := range licenseList {
		licenses = append(licenses, license.Key)
	}
	return licenses, nil
}
