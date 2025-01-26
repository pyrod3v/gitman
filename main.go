package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	loadGitignores()

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
		if err := addGitignore(absPath); err != nil {
			log.Fatalf("Failed to add .gitignore: %v\n", err)
		}
		fmt.Println("Successfully added .gitignore!")
	case "init":
		if err := initializeRepo(absPath); err != nil {
			log.Fatalf("Failed to initialize repository: %v\n", err)
		}
		fmt.Println("Successfully initalized repository!")
	}
}

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
