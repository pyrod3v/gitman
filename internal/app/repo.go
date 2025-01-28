package app

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/manifoldco/promptui"
)

func InitializeRepo(absPath string) error {
	initPrompt := promptui.Prompt{
		Label:   "Enter repository name",
		Default: filepath.Base(absPath),
	}
	name, err := initPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	if name != filepath.Base(absPath) {
		absPath = filepath.Join(filepath.Dir(absPath), name)
		if err := os.MkdirAll(absPath, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}

	cmd := exec.Command("git", "init", absPath)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("failed to initialize git repository: %v", err)
	}
	fmt.Println("Repository initialized successfully!")

	if err := AddGitignore(absPath); err != nil {
		fmt.Printf(".gitignore could not be added to the repository: %v\n", err)
	}

	userPrompt := promptui.Prompt{
		Label:   "Git user.name (leave empty to use default)",
		Default: "",
	}
	user, err := userPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	if user != "" {
		cmd := exec.Command("git", "-C", absPath, "config", "user.name", user)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to set git user.name: %v\n", err)
		}
	}

	emailPrompt := promptui.Prompt{
		Label:   "Git user.email (leave empty to use default)",
		Default: "",
	}
	email, err := emailPrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	if email != "" {
		cmd := exec.Command("git", "-C", absPath, "config", "user.email", email)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			fmt.Printf("Failed to set git user.email: %v\n", err)
		}
	}

	remotePrompt := promptui.Prompt{
		Label:     "Remote repository URL (leave empty to skip)",
		Default:   "",
		AllowEdit: true,
	}
	remoteURL, err := remotePrompt.Run()
	if err != nil {
		log.Fatalf("Prompt failed: %v\n", err)
	}

	if remoteURL != "" {
		cmd := exec.Command("git", "-C", absPath, "remote", "add", "origin", remoteURL)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		err := cmd.Run()
		if err != nil {
			return fmt.Errorf("failed to add remote: %v", err)
		}
		fmt.Println("Remote added successfully!")
	}
	return nil
}
