package main

import (
	"testing"
)

func TestNewExecutor(t *testing.T) {
	config := map[string]string{
		"region": "eu",
		"teamId": "123",
	}

	executor := NewExecutor(config)
	if executor == nil {
		t.Fatal("NewExecutor returned nil")
	}
	defer executor.Close()

	if executor.config["region"] != "eu" {
		t.Errorf("Expected region 'eu', got: %s", executor.config["region"])
	}

	if executor.config["teamId"] != "123" {
		t.Errorf("Expected teamId '123', got: %s", executor.config["teamId"])
	}
}

func TestInterpolate(t *testing.T) {
	config := map[string]string{
		"region": "eu",
		"teamId": "12345",
	}

	executor := NewExecutor(config)
	defer executor.Close()

	// Test config interpolation
	result := executor.interpolate("https://{{config.region}}.example.com")
	expected := "https://eu.example.com"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test multiple config values
	result = executor.interpolate("https://{{config.region}}.example.com/team/{{config.teamId}}")
	expected = "https://eu.example.com/team/12345"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test with stored variables
	executor.variables["invoice"] = "INV-001"
	result = executor.interpolate("Invoice ID: {{invoice}}")
	expected = "Invoice ID: INV-001"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}

	// Test no interpolation needed
	result = executor.interpolate("https://example.com")
	expected = "https://example.com"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestInterpolateComplexVariables(t *testing.T) {
	executor := NewExecutor(nil)
	defer executor.Close()

	// Test with object variable
	type Invoice struct {
		ID    string
		Total string
	}
	executor.variables["invoice.id"] = "INV-123"
	executor.variables["invoice.total"] = "$100.00"

	result := executor.interpolate("Invoice {{invoice.id}} total: {{invoice.total}}")
	expected := "Invoice INV-123 total: $100.00"
	if result != expected {
		t.Errorf("Expected %s, got %s", expected, result)
	}
}

func TestExecuteStepsEmpty(t *testing.T) {
	executor := NewExecutor(nil)
	defer executor.Close()

	// Test with empty steps
	err := executor.ExecuteSteps([]Step{})
	if err != nil {
		t.Errorf("Expected no error for empty steps, got: %v", err)
	}
}

func TestExecuteStepsSleep(t *testing.T) {
	executor := NewExecutor(nil)
	defer executor.Close()

	// Test sleep step
	steps := []Step{
		{
			Action:   "sleep",
			Duration: 100, // 100ms
		},
	}

	err := executor.ExecuteSteps(steps)
	if err != nil {
		t.Errorf("Expected no error for sleep step, got: %v", err)
	}
}

func TestExecuteStepUnsupported(t *testing.T) {
	executor := NewExecutor(nil)
	defer executor.Close()

	// Test unsupported action - should not fail, just log warning
	step := Step{
		Action: "unsupportedAction",
	}

	err := executor.executeStep(step)
	if err != nil {
		t.Errorf("Unsupported action should not return error, got: %v", err)
	}
}
