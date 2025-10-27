package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	// Command-line flags
	pluginFile := flag.String("plugin", "", "Path to the plugin JSON file")
	configFile := flag.String("config", "", "Path to configuration JSON file (optional)")
	checkAuthOnly := flag.Bool("check-auth", false, "Only check authentication")
	flag.Parse()

	if *pluginFile == "" {
		fmt.Println("Usage: invoiceradar-plugins -plugin <plugin.json> [-config <config.json>] [-check-auth]")
		fmt.Println("\nOptions:")
		fmt.Println("  -plugin string")
		fmt.Println("        Path to the plugin JSON file (required)")
		fmt.Println("  -config string")
		fmt.Println("        Path to configuration JSON file with plugin settings (optional)")
		fmt.Println("  -check-auth")
		fmt.Println("        Only check if already authenticated, don't run full plugin")
		fmt.Println("\nExample:")
		fmt.Println("  invoiceradar-plugins -plugin plugins/plausible.json")
		os.Exit(1)
	}

	// Load plugin
	log.Printf("Loading plugin from: %s", *pluginFile)
	plugin, err := LoadPlugin(*pluginFile)
	if err != nil {
		log.Fatalf("Failed to load plugin: %v", err)
	}

	log.Printf("Loaded plugin: %s (%s)", plugin.Name, plugin.ID)
	log.Printf("Description: %s", plugin.Description)

	// Load configuration if provided
	config := make(map[string]string)
	if *configFile != "" {
		data, err := os.ReadFile(*configFile)
		if err != nil {
			log.Fatalf("Failed to read config file: %v", err)
		}
		if err := json.Unmarshal(data, &config); err != nil {
			log.Fatalf("Failed to parse config file: %v", err)
		}
		log.Printf("Loaded configuration with %d values", len(config))
	}

	// Create executor
	executor := NewExecutor(config)
	defer executor.Close()

	// Check authentication
	if len(plugin.CheckAuth) > 0 {
		log.Println("\n=== Checking Authentication ===")
		err := executor.ExecuteSteps(plugin.CheckAuth)
		if err != nil {
			log.Printf("Authentication check failed: %v", err)
			
			if !*checkAuthOnly {
				// Start authentication process
				log.Println("\n=== Starting Authentication ===")
				if len(plugin.StartAuth) > 0 {
					err = executor.ExecuteSteps(plugin.StartAuth)
					if err != nil {
						log.Fatalf("Authentication failed: %v", err)
					}
					log.Println("✓ Authentication successful")
				}
			}
		} else {
			log.Println("✓ Already authenticated")
		}
	}

	if *checkAuthOnly {
		log.Println("Check auth completed")
		return
	}

	// Fetch documents
	if len(plugin.GetDocuments) > 0 {
		log.Println("\n=== Fetching Documents ===")
		err = executor.ExecuteSteps(plugin.GetDocuments)
		if err != nil {
			log.Fatalf("Failed to fetch documents: %v", err)
		}
		log.Println("✓ Document fetch completed")
	}

	log.Println("\n=== Plugin execution completed successfully ===")
	
	// Keep browser open for a moment to see results
	log.Println("Browser will remain open for 5 seconds...")
	fmt.Scanln() // Wait for user input before closing
}
