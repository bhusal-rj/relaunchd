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
