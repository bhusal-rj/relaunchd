package main

import (
	"bhusal-rj/relaunchd/internal/config"
	"bhusal-rj/relaunchd/internal/watcher"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"

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
		// Implementation for starting the process will go here

		if len(cfg.Watch.Paths) > 0 {
			log.Println("Setting up the file watcher...")
			w, err := watcher.New(&cfg.Watch)

			if err != nil {
				fmt.Println("Error creating watcher:", err)
				os.Exit(1)
			}

			err = w.Start()

			if err != nil {
				fmt.Println("Failed to start watcher:", err)
				os.Exit(1)
			}
			log.Println("File watcher started successfully.")

			//Watch for termination signal
			sigChan := make(chan os.Signal, 1)
			signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
			fmt.Println("Press Ctrl+C to stop the watcher...")
			<-sigChan

			fmt.Println("\n Shutting down the watcher...")
			if w != nil {
				w.Stop()
			}
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
