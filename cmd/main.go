package main

import (
	"bhusal-rj/relaunchd/internal/config"
	"fmt"
	"os"
)

const version = "0.1.0-dev"

func main() {
	fmt.Printf("relaunchd %s\n", version)
	fmt.Println("Lightweight process manager and file watcher")

	if len(os.Args) < 2 {
		fmt.Println("Usage: relaunchd <command>")
		fmt.Println("Commands: start,stop,status")
	}
	configPath := "config.yaml"
	// Load the config file
	config, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Println("Error loading config:", err)
		return
	}
	fmt.Println("Config loaded successfully")
	fmt.Println(config)
}
