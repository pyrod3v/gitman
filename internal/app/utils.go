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
    "os"
    "path/filepath"
    "strings"
)

func GetConfigDir() string {
    home, err := os.UserHomeDir()
    if err != nil {
        fmt.Printf("Error retrieving home directory: %v\n", err)
        os.Exit(1)
    }
    return filepath.Join(home, ".gitman")
}

func RemoveDuplicates(slice []string) []string {
    seen := make(map[string]struct{})
    result := []string{}
    for _, item := range slice {
        if _, exists := seen[item]; !exists {
            seen[item] = struct{}{}
            result = append(result, item)
        }
    }
    return result
}

func LoadFromDir(dir string, suffix string) ([]string, error) {
    entries, err := os.ReadDir(dir)
    if err != nil {
        if errors.Is(err, os.ErrNotExist) {
            err := os.MkdirAll(dir, 0755)
            if err != nil {
                err = fmt.Errorf("failed to create directory: %v", err)
            }
            return []string{}, err
        }
        return []string{}, err
    }
    var files []string
    for _, entry := range entries {
        if !entry.IsDir() && strings.HasSuffix(entry.Name(), suffix) {
            files = append(files, strings.TrimSuffix(entry.Name(), suffix))
        }
    }
    return files, nil
}
