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

package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/pyrod3v/gitman/internal/app"
	"github.com/spf13/viper"
)

func main() {
	if _, err := os.Stat(gitman.GetConfigDir()); os.IsNotExist(err) {
		if err := os.MkdirAll(gitman.GetConfigDir(), 0755); err != nil {
			fmt.Printf("Error creating config directory: %v\n", err)
			os.Exit(1)
		}
	}

	configFile := filepath.Join(gitman.GetConfigDir(), "config.yaml")
	if _, err := os.Stat(configFile); err != nil && os.IsNotExist(err) {
		_, err := os.Create(configFile)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}
	}

	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./.gitman")
	viper.AddConfigPath(gitman.GetConfigDir())
	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("error loading config: %w", err))
	}

	viper.SetDefault("CacheGitignores", false)
	viper.WriteConfig()

	gitman.LoadGitignores()

	var action string
	var path string

	actionForm := huh.NewForm(
		huh.NewGroup(
			huh.NewSelect[string]().
				Title("Select Git Action").
				Options(huh.NewOptions("init", "add gitignore")...).
				Value(&action),
		),
	)

	if err := actionForm.Run(); err != nil {
		log.Fatalf("Form failed: %v\n", err)
	}

	pathForm := huh.NewForm(
		huh.NewGroup(
			huh.NewInput().
				Title("Enter path").
				Value(&path).
				Placeholder("."),
		),
	)

	if err := pathForm.Run(); err != nil {
		log.Fatalf("Form failed: %v\n", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v\n", err)
	}

	switch action {
	case "add gitignore":
		if err := gitman.AddGitignore(absPath); err != nil {
			log.Fatalf("Failed to add .gitignore: %v\n", err)
		}
		fmt.Println("Successfully added .gitignore!")
	case "init":
		if err := gitman.InitializeRepo(absPath); err != nil {
			log.Fatalf("Failed to initialize repository: %v\n", err)
		}
		fmt.Println("Successfully initialized repository!")
	}
}
