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
    "os"
    "os/exec"
    "path/filepath"

    "github.com/charmbracelet/huh"
)

func InitializeRepo(absPath string) error {
    var name, user, email, remoteURL string

    form := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().
                Title("Enter repository name").
                Value(&name).
                Placeholder(filepath.Base(absPath)),
        ),
    )

    if err := form.Run(); err != nil {
        fmt.Printf("Repository form failed: %v\n", err)
    }

    if name != "" && name != filepath.Base(absPath) {
        absPath = filepath.Join(filepath.Dir(absPath), name)
        if err := os.MkdirAll(absPath, 0755); err != nil {
            return fmt.Errorf("failed to create directory: %v", err)
        }
    }

    cmd := exec.Command("git", "init", absPath)
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        fmt.Printf("Failed to initialize git repository: %v\n", err)
    }
    fmt.Println("Repository initialized successfully!")

    if err := AddGitignore(absPath); err != nil {
        fmt.Printf(".gitignore could not be added to the repository: %v\n", err)
    }

    if err := AddLicense(absPath); err != nil {
        fmt.Printf("License could not be added to the repository: %v\n", err)
    }

    userForm := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().Title("Git user.name (leave empty to use default)").Value(&user),
            huh.NewInput().Title("Git user.email (leave empty to use default)").Value(&email),
        ),
    )

    if err := userForm.Run(); err != nil {
        fmt.Printf("User form failed: %v\n", err)
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

    if email != "" {
        cmd := exec.Command("git", "-C", absPath, "config", "user.email", email)
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr
        err := cmd.Run()
        if err != nil {
            fmt.Printf("Failed to set git user.email: %v\n", err)
        }
    }

    remoteForm := huh.NewForm(
        huh.NewGroup(
            huh.NewInput().Title("Remote repository URL (leave empty to skip)").Value(&remoteURL),
        ),
    )

    if err := remoteForm.Run(); err != nil {
        fmt.Printf("Remote form failed: %v\n", err)
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
