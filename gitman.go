package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/manifoldco/promptui"
)

func main() {
	actionPrompt := promptui.Select{
		Label: "Select Git Action",
		Items: []string{"init"},
	}

	_, action, err := actionPrompt.Run()
	if err != nil {
		log.Fatalf("Action prompt failed %v\n", err)
	}

	pathPrompt := promptui.Prompt{
		Label:   "Enter path",
		Default: ".",
	}

	path, err := pathPrompt.Run()
	if err != nil {
		log.Fatalf("Path prompt failed %v\n", err)
	}

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v\n", err)
	}

	if action == "init" {
		initPrompt := promptui.Prompt{
			Label:   "Enter repository name",
			Default: filepath.Base(absPath),
		}

		name, err := initPrompt.Run()
		if err != nil {
			log.Fatalf("Init prompt failed %v\n", err)
		}

		if name != filepath.Base(absPath) {
			absPath = filepath.Join(filepath.Dir(absPath), name)
			err = os.MkdirAll(absPath, 0755)
			if err != nil {
				log.Fatalf("Failed to create directory: %v\n", err)
			}
		}

		cmd := exec.Command("git", "init", absPath)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to initialize git repository: %v\n", err)
		}
		fmt.Println("Git repository initialized successfully!")

		err = addGitignore(absPath)
		if err != nil {
			_ = fmt.Errorf(".gitignore could not be added to the repository: %n", err)
		}

		userPrompt := promptui.Prompt{
			Label:   "Git user.name (leave empty to use default)",
			Default: "",
		}

		user, _ := userPrompt.Run()

		if user != "" {
			cmd := exec.Command("git", "-C", absPath, "config", "user.name", user)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				_ = fmt.Errorf("failed to set git user.name: %v", err)
			}
		}

		emailPrompt := promptui.Prompt{
			Label:   "Git user.email (leave empty to use default)",
			Default: "",
		}

		email, _ := emailPrompt.Run()

		if email != "" {
			cmd := exec.Command("git", "-C", absPath, "config", "user.email", email)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				_ = fmt.Errorf("failed to set git user.email: %v", err)
			}
		}

		remotePrompt := promptui.Prompt{
			Label:     "Remote repository URL (leave empty to skip)",
			Default:   "",
			AllowEdit: true,
		}

		remoteURL, _ := remotePrompt.Run()

		if remoteURL != "" {
			cmd := exec.Command("git", "-C", absPath, "remote", "add", "origin", remoteURL)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			if err := cmd.Run(); err != nil {
				log.Fatalf("Failed to add remote: %v\n", err)
			}
			fmt.Println("Remote added successfully!")
		}
	}
}

func addGitignore(path string) error {
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
	templates = append([]string{"None"}, templates...)
	gitignorePrompt := promptui.Select{
		Label: "Select Gitignore Template",
		Items: templates,
		Searcher: func(input string, index int) bool {
			template := strings.ToLower(templates[index])
			return strings.Contains(template, strings.ToLower(input))
		},
	}

	_, template, err := gitignorePrompt.Run()
	if err != nil {
		return fmt.Errorf("gitignore prompt failed %v", err)
	}

	if template != "None" {
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
	}
	return nil
}
