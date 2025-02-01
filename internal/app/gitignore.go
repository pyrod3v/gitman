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

package app

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"
)

var gitignores []string
var gitignoresMutex sync.Mutex

func AddGitignore(path string) error {
	gitignorePrompt := promptui.Select{
		Label: "Select Gitignore Template",
		Items: gitignores,
		Searcher: func(input string, index int) bool {
			template := strings.ToLower(gitignores[index])
			return strings.Contains(template, strings.ToLower(input))
		},
	}

	_, template, err := gitignorePrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	if template == "None" {
		return nil
	}

	gitignorePath := filepath.Join(getConfigDir(), "gitignores", template+".gitignore")
	if content, err := os.ReadFile(gitignorePath); err == nil {
		return os.WriteFile(filepath.Join(path, ".gitignore"), content, 0644)
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

	err = os.WriteFile(filepath.Join(path, ".gitignore"), content, 0644)
	if err != nil {
		return fmt.Errorf("failed to write .gitignore: %v", err)
	}

	fmt.Println(".gitignore added!")
	return nil
}

func LoadGitignores() {
	go func() {
		if err := fetchGitignores(); err != nil {
			log.Panicf("Failed to fetch gitignore templates: %v\n", err)
		}
	}()
	if customGitignores, err := loadCustomGitignores(); err == nil {
		gitignoresMutex.Lock()
		gitignores = append(customGitignores, gitignores...)
		gitignoresMutex.Unlock()
	}
	gitignoresMutex.Lock()
	sort.Strings(gitignores)
	gitignores = append([]string{"None"}, gitignores...) // ensure 'None' is the first option
	gitignoresMutex.Unlock()
}

func loadCustomGitignores() ([]string, error) {
	gitignoreDir := filepath.Join(getConfigDir(), "gitignores")
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

func fetchGitignores() error {
	res, err := http.Get("https://www.toptal.com/developers/gitignore/api/list")
	if err != nil {
		return fmt.Errorf("failed to fetch gitignore template list: %v", err)
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return fmt.Errorf("failed to read response body: %v", err)
	}

	templates := strings.Split(string(body), ",")
	gitignoresMutex.Lock()
	defer gitignoresMutex.Unlock()
	gitignores = append(gitignores, templates...)
	sort.Strings(gitignores)
	return nil
}
