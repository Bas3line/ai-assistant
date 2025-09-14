package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	// Create bin directory if it doesn't exist
	binDir := "bin"
	if err := os.MkdirAll(binDir, 0755); err != nil {
		fmt.Printf("Error creating bin directory: %v\n", err)
		os.Exit(1)
	}

	// Get the project root directory
	projectRoot, err := os.Getwd()
	if err != nil {
		fmt.Printf("Error getting current directory: %v\n", err)
		os.Exit(1)
	}

	// Move to project root if we're in build directory
	if filepath.Base(projectRoot) == "build" {
		projectRoot = filepath.Dir(projectRoot)
		if err := os.Chdir(projectRoot); err != nil {
			fmt.Printf("Error changing to project root: %v\n", err)
			os.Exit(1)
		}
	}

	// Define build targets
	targets := []struct {
		name     string
		mainPath string
		output   string
	}{
		{
			name:     "api",
			mainPath: "./cmd/api",
			output:   filepath.Join(binDir, "api"),
		},
	}

	// Add platform-specific extension for Windows
	if runtime.GOOS == "windows" {
		for i := range targets {
			targets[i].output += ".exe"
		}
	}

	fmt.Println("Building AI Assistant...")
	fmt.Printf("Project root: %s\n", projectRoot)
	fmt.Printf("Target OS: %s\n", runtime.GOOS)
	fmt.Printf("Target Arch: %s\n", runtime.GOARCH)
	fmt.Println()

	// Build each target
	for _, target := range targets {
		fmt.Printf("Building %s...\n", target.name)
		
		cmd := exec.Command("go", "build", "-o", target.output, target.mainPath)
		cmd.Env = append(os.Environ(),
			"CGO_ENABLED=0",
			"GOOS="+runtime.GOOS,
			"GOARCH="+runtime.GOARCH,
		)
		
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		
		if err := cmd.Run(); err != nil {
			fmt.Printf("Error building %s: %v\n", target.name, err)
			os.Exit(1)
		}
		
		fmt.Printf("✓ Built %s -> %s\n", target.name, target.output)
	}

	fmt.Println()
	fmt.Println("Build completed successfully!")
	fmt.Println("Binaries are available in the bin/ directory:")
	
	// List built binaries
	files, err := os.ReadDir(binDir)
	if err != nil {
		fmt.Printf("Error reading bin directory: %v\n", err)
		return
	}
	
	for _, file := range files {
		if !file.IsDir() {
			fmt.Printf("  - %s\n", filepath.Join(binDir, file.Name()))
		}
	}
}