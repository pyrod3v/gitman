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
    "fmt"
    "io"
    "net/http"
    "os"
    "path/filepath"
    "sort"
    "strings"
    "sync"

    "github.com/charmbracelet/bubbles/key"
    "github.com/charmbracelet/huh"
    "github.com/spf13/viper"
)

var gitignores []string
var gitignoresMutex sync.Mutex

func AddGitignore(path string) error {
    var selected []string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewMultiSelect[string]().
                Title("Select .gitignore templates").
                Options(huh.NewOptions(gitignores...)...).
                Value(&selected).
                Height(min(10, len(gitignores))),
        ),
    ).WithKeyMap(func(k *huh.KeyMap) *huh.KeyMap {
        k.Quit = key.NewBinding(key.WithKeys("q", "ctrl+c"))
        return k
    }(huh.NewDefaultKeyMap()))

    if err := form.Run(); err != nil {
        return fmt.Errorf("form failed: %v", err)
    }

    if len(selected) == 0 {
        return nil
    }

    var builder strings.Builder
    for _, template := range selected {
        gitignorePath := filepath.Join(GetConfigDir(), "gitignores", template+".gitignore")
        if content, err := os.ReadFile(gitignorePath); err == nil {
            builder.WriteString(string(content))
            builder.WriteString("\n")
            continue
        }

        gitignorePath = filepath.Join(GetConfigDir(), ".cache", "gitignores", template+".gitignore")
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

        if viper.GetBool("CacheGitignores") {
            dir := filepath.Join(GetConfigDir(), ".cache", "gitignores")
            if err := os.MkdirAll(dir, 0755); err != nil {
                return fmt.Errorf("failed to create cache directory: %v", err)
            }
            cacheFile := filepath.Join(dir, template+".gitignore")
            err := os.WriteFile(cacheFile, content, 0644)
            if err != nil {
                fmt.Printf("Failed to cache gitignore template: %v\n", err)
            }
        }

        builder.WriteString(string(content))
        builder.WriteString("\n")
    }

    err := os.WriteFile(filepath.Join(path, ".gitignore"), []byte(builder.String()), 0644)
    if err != nil {
        return fmt.Errorf("failed to write .gitignore: %v", err)
    }

    return nil
}

func LoadGitignores() {
    gitignoresMutex.Lock()
    gitignores = append(gitignores, loadCustomGitignores()...)
    gitignoresMutex.Unlock()

    go func() {
        fetchedGitignores, err := fetchGitignores()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            return
        }
        gitignoresMutex.Lock()
        gitignores = append(gitignores, fetchedGitignores...)
        gitignores = RemoveDuplicates(gitignores)
        sort.Strings(gitignores)
        gitignoresMutex.Unlock()
    }()
}

func loadCustomGitignores() []string {
    customGitignores, err := LoadFromDir(filepath.Join(GetConfigDir(), "gitignores"), ".gitignore")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading custom gitignores: %v", err)
    }

    cachedGitignores, err := LoadFromDir(filepath.Join(GetConfigDir(), ".cache", "gitignores"), ".gitignore")
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error loading cached gitignores: %v", err)
    }

    return append(customGitignores, cachedGitignores...)
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

    return strings.Fields(strings.Join(strings.Split(string(body), ","), "\n")), nil
}
