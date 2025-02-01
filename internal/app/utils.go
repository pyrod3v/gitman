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

import(
	"os"
	"strings"
	"path/filepath"
)

func getConfigDir() string {
	if configDir := os.Getenv("XDG_CONFIG_HOME"); configDir != "" {
		return configDir
	}
	homeDir, _ := os.UserHomeDir()
	if strings.Contains(strings.ToLower(os.Getenv("OS")), "windows") {
		return filepath.Join(homeDir, "AppData", "Roaming", "gitman")
	}
	return filepath.Join(homeDir, ".config", "gitman")
}
