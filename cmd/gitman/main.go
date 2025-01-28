package main

import (
	"fmt"
	"log"
	"path/filepath"

	"github.com/pyrod3v/gitman/internal/app"
	"github.com/manifoldco/promptui"
)

func main() {
	app.LoadGitignores()

	actionPrompt := promptui.Select{
		Label: "Select Git Action",
		Items: []string{"init", "add gitignore"},
	}
	_, action, err := actionPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	pathPrompt := promptui.Prompt{
		Label:   "Enter path",
		Default: ".",
	}
	path, err := pathPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
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
		fmt.Println("Successfully initalized repository!")
	}
}

