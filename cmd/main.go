package main

import (
	"bhusal-rj/relaunchd/internal/config"
	manager "bhusal-rj/relaunchd/internal/process"
	"bhusal-rj/relaunchd/internal/watcher"
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

var (
	configPath string
	verbose    bool
)

const version = "0.1.0-dev"

var rootCmd = &cobra.Command{
	Use:     "relaunchd",
	Short:   "Lightweight process manager and file watcher",
	Long:    `relaunchd is a lightweight process manager and file watcher that helps you manage and monitor your applications.`,
	Version: version,
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start a process defined in config",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		fmt.Printf("Starting application '%s'...\n", cfg.Name)

		// Set up signal handling BEFORE starting processes
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

		// Create process manager
		processManager := manager.NewProcessManager(cfg)

		if err := processManager.Start(); err != nil {
			fmt.Println("Error starting process:", err)
			os.Exit(1)
		}

		fmt.Println("Process started successfully.")

		var fileWatcher *watcher.Watcher

		// Only set up file watching if paths are configured
		if len(cfg.Watch.Paths) > 0 {
			log.Println("Setting up the file watcher...")

			fileWatcher, err = watcher.New(&cfg.Watch)
			if err != nil {
				fmt.Println("Error creating watcher:", err)
			} else {
				fileWatcher.SetChangeHandler(processManager.TriggerRestart)

				if err := fileWatcher.Start(); err != nil {
					fmt.Println("Failed to start watcher:", err)
					fileWatcher = nil
				} else {
					log.Println("File watcher started successfully.")
				}
			}
		}

		fmt.Println("Press Ctrl+C to stop the application...")

		// Wait for termination signal
		sig := <-sigChan
		fmt.Printf("\nReceived signal %v, shutting down...\n", sig)

		// Use a context with timeout for clean shutdown
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		// Stop components with a timeout
		done := make(chan bool, 1)
		go func() {
			// Stop watcher if it was started
			if fileWatcher != nil {
				fmt.Println("Stopping the watcher...")
				fileWatcher.Stop()
			}

			// Always stop the process
			fmt.Println("Stopping process...")
			err := processManager.Stop()
			if err != nil {
				fmt.Printf("Error stopping process: %v\n", err)
			} else {
				fmt.Println("Process stopped successfully.")
			}

			done <- true
		}()

		// Wait for clean shutdown or timeout
		select {
		case <-done:
			fmt.Println("Clean shutdown completed.")
		case <-ctx.Done():
			fmt.Println("Shutdown timed out, forcing exit.")
		}
	},
}

var stopCmd = &cobra.Command{
	Use:   "stop",
	Short: "Stop a running process",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		fmt.Printf("Stopping application '%s'...\n", cfg.Name)
		// Implementation for stopping the process will go here
	},
}
var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Generate a default configuration file",
	Long:  `Creates a default config.yaml file in the current directory to help you get started.`,
	Run: func(cmd *cobra.Command, args []string) {
		targetPath := "config.yaml"
		if configPath != "config.yaml" {
			targetPath = configPath
		}

		// Check if file already exists
		if _, err := os.Stat(targetPath); err == nil {
			fmt.Printf("Configuration file already exists at %s\n", targetPath)
			fmt.Print("Do you want to overwrite it? (y/N): ")

			var response string
			fmt.Scanln(&response)

			if response != "y" && response != "Y" {
				fmt.Println("Operation cancelled.")
				return
			}
		}

		// Template for default config
		defaultConfig := `# relaunchd configuration file
# Generated on %s

# Project name (used for identification)
name: "my-application"

# Command configuration
command:
  # The command to execute
  exec: "node server.js"
  
  # Working directory (optional, defaults to ".")
  cwd: "."
  
  # Environment variables
  env:
    NODE_ENV: "development"
    PORT: "3000"

# File watching configuration
watch:
  # Paths to monitor for changes (can use glob patterns)
  paths:
    - "."
  
  # Paths to exclude from monitoring
  exclude:
    - "node_modules/**"
    - "*.log"
    - "*.tmp"
    - ".git/**"

# Process management configuration
process:
  # Whether to run in the background
  background: false
  
  # Maximum number of restarts (0 for unlimited)
  max_restarts: 10
  
  # Delay between restarts (milliseconds)
  restart_delay: 1000
`
		// Format the config with the current date/time
		formattedConfig := fmt.Sprintf(defaultConfig, time.Now().Format("2006-01-02 15:04:05"))

		// Write the config to file
		err := os.WriteFile(targetPath, []byte(formattedConfig), 0644)
		if err != nil {
			fmt.Printf("Error creating config file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("Configuration file created at: %s\n", targetPath)
		fmt.Println("Edit this file to customize your application settings.")
		fmt.Println("\nStart your application with:")
		fmt.Printf("  relaunchd start --config %s\n", targetPath)
	},
}
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show the status of processes",
	Run: func(cmd *cobra.Command, args []string) {
		cfg, err := config.LoadConfig(configPath)
		if err != nil {
			fmt.Println("Error loading config:", err)
			os.Exit(1)
		}
		fmt.Printf("Status for application '%s':\n", cfg.Name)
		// Implementation for checking status will go here
	},
}

func init() {
	// Add persistent flags to the root command
	rootCmd.PersistentFlags().StringVarP(&configPath, "config", "c", "config.yaml", "Path to config file")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "Enable verbose output")

	// Add commands to the root command
	rootCmd.AddCommand(startCmd)
	rootCmd.AddCommand(stopCmd)
	rootCmd.AddCommand(statusCmd)
	rootCmd.AddCommand(initCmd) // Add the new init command

	// Remove the completion command
	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	// Initial display name and description of the application
	fmt.Printf("relaunchd %s\n", version)
	fmt.Println("Lightweight process manager and file watcher")

	// Execute the root commands
	if err := rootCmd.Execute(); err != nil {
		fmt.Println("Error:", err)
		os.Exit(1)
	}
}
