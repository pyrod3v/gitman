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
)

var gitignores []string
var gitignoresMutex sync.Mutex

func AddGitignore(path string) error {
	var selectedTemplates []string

	form := huh.NewForm(
		huh.NewGroup(
			huh.NewMultiSelect[string]().
				Title("Select .gitignore templates").
				Options(huh.NewOptions(gitignores...)...).
				Value(&selectedTemplates).
				Height(20),
		),
	)

	if err := form.Run(); err != nil {
		return fmt.Errorf("form failed: %v", err)
	}

	if len(selectedTemplates) == 0 {
		return nil
	}

	var builder strings.Builder
	for _, template := range selectedTemplates {
		gitignorePath := filepath.Join(GetConfigDir(), "gitignores", template+".gitignore")
		if content, err := os.ReadFile(gitignorePath); err == nil {
			builder.WriteString(string(content))
			builder.WriteString("\n")
			continue
		}

		url := fmt.Sprintf("https://www.toptal.com/developers/gitignore/api/%s", template)
		res, err := http.Get(url)
		if err != nil {
			return fmt.Errorf("failed to fetch gitignore template: %v", err)
		}
		defer res.Body.Close()

		content, err := io.ReadAll(res.Body)
		if err != nil {
			return fmt.Errorf("failed to read gitignore template: %v", err)
		}

		builder.WriteString(string(content))
		builder.WriteString("\n")
	}

	err := os.WriteFile(filepath.Join(path, ".gitignore"), []byte(builder.String()), 0644)
	if err != nil {
		return fmt.Errorf("failed to write .gitignore: %v", err)
	}

	fmt.Println(".gitignore added!")
	return nil
}

func LoadGitignores() {
	if customGitignores, err := loadCustomGitignores(); err == nil {
		gitignoresMutex.Lock()
		gitignores = append(customGitignores, gitignores...)
		gitignoresMutex.Unlock()
	}
	go func() {
		fetchedGitignores, err := fetchGitignores()
		if err != nil {
			fmt.Printf("Failed to fetch gitignore templates: %v\n", err)
		}
		gitignoresMutex.Lock()
		// the fetched gtitignore list has newlines, so remove them
		gitignores = append(gitignores, strings.Fields(strings.Join(fetchedGitignores, "\n"))...)
		sort.Strings(gitignores)
		gitignoresMutex.Unlock()
	}()
}

func loadCustomGitignores() ([]string, error) {
	gitignoreDir := filepath.Join(GetConfigDir(), "gitignores")
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
		if !entry.IsDir() && strings.HasSuffix(entry.Name(), ".gitignore") {
			customGitignores = append(customGitignores, strings.TrimSuffix(entry.Name(), ".gitignore"))
		}
	}
	return customGitignores, nil
}

func fetchGitignores() ([]string, error) {
	res, err := http.Get("https://www.toptal.com/developers/gitignore/api/list")
	if err != nil {
		return nil, fmt.Errorf("failed to fetch gitignore template list: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %v", err)
	}

	return strings.Split(string(body), ","), nil
}
