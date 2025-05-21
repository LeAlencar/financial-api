package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Warning: Error loading .env file:", err)
	}

	// Get the absolute path to the migrations directory
	currentDir, err := os.Getwd()
	if err != nil {
		panic(fmt.Sprintf("Failed to get current directory: %v", err))
	}

	migrationsPath := filepath.Join(currentDir, "internal", "infra", "database", "migrations")
	configPath := filepath.Join(migrationsPath, "tern.conf")

	// Verify paths exist
	if _, err := os.Stat(migrationsPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Migrations directory does not exist: %s", migrationsPath))
	}
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic(fmt.Sprintf("Tern config file does not exist: %s", configPath))
	}

	fmt.Printf("Using migrations path: %s\n", migrationsPath)
	fmt.Printf("Using config path: %s\n", configPath)

	cmd := exec.Command(
		"tern",
		"migrate",
		"--migrations",
		migrationsPath,
		"--config",
		configPath,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Migration failed: ", err)
		fmt.Println("Output: ", string(output))
		panic(err)
	}

	fmt.Println("Migrations executed successfully:", string(output))
}
