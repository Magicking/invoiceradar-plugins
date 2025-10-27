package main

import (
	"os"
	"testing"
)

func TestLoadPlugin(t *testing.T) {
	// Test loading the blank plugin
	plugin, err := LoadPlugin("blank-plugin.json")
	if err != nil {
		t.Fatalf("Failed to load blank plugin: %v", err)
	}

	if plugin.Schema == "" {
		t.Error("Plugin schema should not be empty")
	}

	if plugin.ID != "" {
		t.Errorf("Blank plugin ID should be empty, got: %s", plugin.ID)
	}
}

func TestLoadPluginPlausible(t *testing.T) {
	// Test loading a real plugin
	plugin, err := LoadPlugin("plugins/plausible.json")
	if err != nil {
		t.Fatalf("Failed to load plausible plugin: %v", err)
	}

	if plugin.ID != "plausible" {
		t.Errorf("Expected plugin ID 'plausible', got: %s", plugin.ID)
	}

	if plugin.Name != "Plausible" {
		t.Errorf("Expected plugin name 'Plausible', got: %s", plugin.Name)
	}

	if len(plugin.CheckAuth) == 0 {
		t.Error("Plugin should have checkAuth steps")
	}

	if len(plugin.StartAuth) == 0 {
		t.Error("Plugin should have startAuth steps")
	}

	if len(plugin.GetDocuments) == 0 {
		t.Error("Plugin should have getDocuments steps")
	}
}

func TestLoadPluginPostHog(t *testing.T) {
	// Test loading a plugin with configuration schema
	plugin, err := LoadPlugin("plugins/posthog.json")
	if err != nil {
		t.Fatalf("Failed to load posthog plugin: %v", err)
	}

	if plugin.ID != "posthog" {
		t.Errorf("Expected plugin ID 'posthog', got: %s", plugin.ID)
	}

	if len(plugin.ConfigSchema) == 0 {
		t.Error("PostHog plugin should have configSchema")
	}

	if _, exists := plugin.ConfigSchema["region"]; !exists {
		t.Error("PostHog plugin should have 'region' in configSchema")
	}
}

func TestLoadPluginNonExistent(t *testing.T) {
	// Test loading a non-existent plugin
	_, err := LoadPlugin("non-existent-plugin.json")
	if err == nil {
		t.Error("Expected error when loading non-existent plugin")
	}
}

func TestLoadPluginInvalidJSON(t *testing.T) {
	// Create a temporary invalid JSON file
	tmpFile, err := os.CreateTemp("", "invalid-*.json")
	if err != nil {
		t.Fatalf("Failed to create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Write invalid JSON
	tmpFile.WriteString("{invalid json")
	tmpFile.Close()

	// Try to load it
	_, err = LoadPlugin(tmpFile.Name())
	if err == nil {
		t.Error("Expected error when loading invalid JSON")
	}
}
