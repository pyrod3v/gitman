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
