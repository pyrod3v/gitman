package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"sync"

	"github.com/manifoldco/promptui"
)

var gitignores []string
var gitignoresMutex sync.Mutex

func main() {
	go func() {
		if err := fetchGitignores(); err != nil {
			log.Printf("Failed to fetch gitignore templates: %v\n", err)
		}
	}()
	if customGitignores, err := loadCustomGitignores(); err == nil {
		gitignoresMutex.Lock()
		gitignores = append(customGitignores, gitignores...)
		sort.Strings(gitignores)
		gitignoresMutex.Unlock()
	}
	gitignores = append([]string{"None"}, gitignores...)

	actionPrompt := promptui.Select{
		Label: "Select Git Action",
		Items: []string{"init", "add gitignore"},
	}
	_, action, _ := actionPrompt.Run()

	pathPrompt := promptui.Prompt{
		Label:   "Enter path",
		Default: ".",
	}
	path, _ := pathPrompt.Run()

	absPath, err := filepath.Abs(path)
	if err != nil {
		log.Fatalf("Failed to resolve path: %v\n", err)
	}

	switch action {
	case "add gitignore":
		if err := addGitignore(absPath); err != nil {
			log.Fatalf("Failed to add .gitignore: %v\n", err)
		}
	case "init":
		if err := initializeRepo(absPath); err != nil {
			log.Fatalf("Failed to initialize repository: %v\n", err)
		}
	}
}

func initializeRepo(absPath string) error {
	initPrompt := promptui.Prompt{
		Label:   "Enter repository name",
		Default: filepath.Base(absPath),
	}
	name, err := initPrompt.Run()
	if err != nil {
		return fmt.Errorf("init prompt failed: %v", err)
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

	if err := addGitignore(absPath); err != nil {
		fmt.Printf(".gitignore could not be added to the repository: %v\n", err)
	}

	userPrompt := promptui.Prompt{
		Label:   "Git user.name (leave empty to use default)",
		Default: "",
	}
	user, _ := userPrompt.Run()

	if user != "" {
		cmd := exec.Command("git", "-C", absPath, "config", user)
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
	email, _ := emailPrompt.Run()

	if email != "" {
		cmd := exec.Command("git", "-C", absPath, "config", email)
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
	remoteURL, _ := remotePrompt.Run()

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

func addGitignore(path string) error {
	gitignorePrompt := promptui.Select{
		Label: "Select Gitignore Template",
		Items: gitignores,
		Searcher: func(input string, index int) bool {
			template := strings.ToLower(gitignores[index])
			return strings.Contains(template, strings.ToLower(input))
		},
	}

	_, template, _ := gitignorePrompt.Run()
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
	return nil
}
