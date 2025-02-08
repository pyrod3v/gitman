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
	"path/filepath"

	"github.com/charmbracelet/huh"
	"github.com/pyrod3v/gitman/internal/app"
)

func main() {
	app.LoadGitignores()

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
		if err := app.AddGitignore(absPath); err != nil {
			log.Fatalf("Failed to add .gitignore: %v\n", err)
		}
		fmt.Println("Successfully added .gitignore!")
	case "init":
		if err := app.InitializeRepo(absPath); err != nil {
			log.Fatalf("Failed to initialize repository: %v\n", err)
		}
		fmt.Println("Successfully initialized repository!")
	}
}
