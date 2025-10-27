package main

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/chromedp"
)

// Executor runs plugin steps using chromedp
type Executor struct {
	ctx       context.Context
	cancel    context.CancelFunc
	variables map[string]interface{}
	config    map[string]string
}

// NewExecutor creates a new plugin executor
func NewExecutor(config map[string]string) *Executor {
	// Create chromedp context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", false), // Show browser for authentication
		chromedp.Flag("disable-blink-features", "AutomationControlled"),
		chromedp.UserAgent("Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/119.0.0.0 Safari/537.36"),
	)

	allocCtx, _ := chromedp.NewExecAllocator(context.Background(), opts...)
	ctx, cancel := chromedp.NewContext(allocCtx, chromedp.WithLogf(log.Printf))

	return &Executor{
		ctx:       ctx,
		cancel:    cancel,
		variables: make(map[string]interface{}),
		config:    config,
	}
}

// Close cleans up the executor
func (e *Executor) Close() {
	e.cancel()
}

// interpolate replaces {{variable}} placeholders with actual values
func (e *Executor) interpolate(text string) string {
	result := text

	// Replace config variables
	for key, value := range e.config {
		result = strings.ReplaceAll(result, fmt.Sprintf("{{config.%s}}", key), value)
	}

	// Replace stored variables
	for key, value := range e.variables {
		result = strings.ReplaceAll(result, fmt.Sprintf("{{%s}}", key), fmt.Sprintf("%v", value))
	}

	return result
}

// ExecuteSteps runs a series of plugin steps
func (e *Executor) ExecuteSteps(steps []Step) error {
	for i, step := range steps {
		log.Printf("Executing step %d: %s", i+1, step.Action)
		if err := e.executeStep(step); err != nil {
			return fmt.Errorf("step %d (%s) failed: %w", i+1, step.Action, err)
		}
	}
	return nil
}

// executeStep executes a single plugin step
func (e *Executor) executeStep(step Step) error {
	switch step.Action {
	case "navigate":
		return e.navigate(step)
	case "waitForElement":
		return e.waitForElement(step)
	case "waitForURL":
		return e.waitForURL(step)
	case "waitForNavigation":
		return e.waitForNavigation(step)
	case "waitForNetworkIdle":
		return e.waitForNetworkIdle(step)
	case "checkElementExists":
		return e.checkElementExists(step)
	case "checkURL":
		return e.checkURL(step)
	case "click":
		return e.click(step)
	case "type":
		return e.typeText(step)
	case "dropdownSelect":
		return e.dropdownSelect(step)
	case "extract":
		return e.extract(step)
	case "extractAll":
		return e.extractAll(step)
	case "downloadPdf":
		return e.downloadPdf(step)
	case "printPdf":
		return e.printPdf(step)
	case "sleep":
		return e.sleep(step)
	case "runJs":
		return e.runJs(step)
	case "if":
		return e.ifCondition(step)
	default:
		log.Printf("Warning: Unsupported action '%s', skipping", step.Action)
		return nil
	}
}

// Navigation steps
func (e *Executor) navigate(step Step) error {
	url := e.interpolate(step.URL)
	log.Printf("Navigating to: %s", url)

	tasks := chromedp.Tasks{
		chromedp.Navigate(url),
	}

	if step.WaitForNetworkIdle {
		tasks = append(tasks, chromedp.WaitReady("body"))
	}

	return chromedp.Run(e.ctx, tasks)
}

func (e *Executor) waitForElement(step Step) error {
	selector := e.interpolate(step.Selector)
	timeout := time.Duration(step.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(e.ctx, timeout)
	defer cancel()

	return chromedp.Run(ctx, chromedp.WaitVisible(selector))
}

func (e *Executor) waitForURL(step Step) error {
	expectedURL := e.interpolate(step.URL)
	timeout := time.Duration(step.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	ctx, cancel := context.WithTimeout(e.ctx, timeout)
	defer cancel()

	// Wait and check URL periodically
	start := time.Now()
	for time.Since(start) < timeout {
		var currentURL string
		if err := chromedp.Run(ctx, chromedp.Location(&currentURL)); err == nil {
			if strings.Contains(expectedURL, "**") {
				prefix := strings.Split(expectedURL, "**")[0]
				if strings.HasPrefix(currentURL, prefix) {
					return nil
				}
			} else if strings.Contains(currentURL, expectedURL) {
				return nil
			}
		}
		time.Sleep(500 * time.Millisecond)
	}

	return fmt.Errorf("timeout waiting for URL: %s", expectedURL)
}

func (e *Executor) waitForNavigation(step Step) error {
	timeout := time.Duration(step.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 10 * time.Second
	}

	time.Sleep(timeout)
	return nil
}

func (e *Executor) waitForNetworkIdle(step Step) error {
	timeout := time.Duration(step.Timeout) * time.Millisecond
	if timeout == 0 {
		timeout = 15 * time.Second
	}

	time.Sleep(2 * time.Second) // Simplified network idle wait
	return nil
}

// Verification steps
func (e *Executor) checkElementExists(step Step) error {
	selector := e.interpolate(step.Selector)
	var nodes []*cdp.Node

	err := chromedp.Run(e.ctx, chromedp.Nodes(selector, &nodes, chromedp.ByQuery))
	if err != nil || len(nodes) == 0 {
		return fmt.Errorf("element not found: %s", selector)
	}

	log.Printf("✓ Element exists: %s", selector)
	return nil
}

func (e *Executor) checkURL(step Step) error {
	expectedURL := e.interpolate(step.URL)
	var currentURL string

	err := chromedp.Run(e.ctx, chromedp.Location(&currentURL))
	if err != nil {
		return err
	}

	// Simple wildcard matching
	if strings.Contains(expectedURL, "**") {
		prefix := strings.Split(expectedURL, "**")[0]
		if !strings.HasPrefix(currentURL, prefix) {
			return fmt.Errorf("URL does not match: expected %s, got %s", expectedURL, currentURL)
		}
	} else if !strings.Contains(currentURL, expectedURL) {
		return fmt.Errorf("URL does not match: expected %s, got %s", expectedURL, currentURL)
	}

	log.Printf("✓ URL matches: %s", currentURL)
	return nil
}

// Interaction steps
func (e *Executor) click(step Step) error {
	selector := e.interpolate(step.Selector)
	log.Printf("Clicking: %s", selector)
	return chromedp.Run(e.ctx, chromedp.Click(selector))
}

func (e *Executor) typeText(step Step) error {
	selector := e.interpolate(step.Selector)
	value := e.interpolate(step.Value)
	log.Printf("Typing into %s", selector)
	return chromedp.Run(e.ctx, chromedp.SendKeys(selector, value))
}

func (e *Executor) dropdownSelect(step Step) error {
	selector := e.interpolate(step.Selector)
	value := e.interpolate(step.Value)
	log.Printf("Selecting %s in dropdown %s", value, selector)
	return chromedp.Run(e.ctx, chromedp.SetValue(selector, value))
}

// Data extraction steps
func (e *Executor) extract(step Step) error {
	log.Printf("Extracting data into variable: %s", step.Variable)
	// Simplified extraction - would need more complex implementation
	return nil
}

func (e *Executor) extractAll(step Step) error {
	log.Printf("Extracting all items with selector: %s", step.Selector)

	// Execute forEach steps for each item (simplified)
	if len(step.ForEach) > 0 {
		log.Printf("Executing forEach steps")
		return e.ExecuteSteps(step.ForEach)
	}

	return nil
}

// Document retrieval steps
func (e *Executor) downloadPdf(step Step) error {
	url := e.interpolate(step.URL)
	log.Printf("Downloading PDF from: %s", url)
	// In a real implementation, this would download the PDF
	return nil
}

func (e *Executor) printPdf(step Step) error {
	log.Printf("Printing page as PDF")
	// In a real implementation, this would use chromedp to print the page
	var buf []byte
	err := chromedp.Run(e.ctx, chromedp.CaptureScreenshot(&buf))
	if err != nil {
		return err
	}
	log.Printf("Captured page (screenshot as proof of concept)")
	return nil
}

// Miscellaneous steps
func (e *Executor) sleep(step Step) error {
	duration := time.Duration(step.Duration) * time.Millisecond
	log.Printf("Sleeping for %v", duration)
	time.Sleep(duration)
	return nil
}

func (e *Executor) runJs(step Step) error {
	script := e.interpolate(step.Script)
	log.Printf("Running JavaScript")
	var result interface{}
	return chromedp.Run(e.ctx, chromedp.Evaluate(script, &result))
}

func (e *Executor) ifCondition(step Step) error {
	script := e.interpolate(step.Script)
	var result bool

	err := chromedp.Run(e.ctx, chromedp.Evaluate(script, &result))
	if err != nil {
		return err
	}

	if result && len(step.Then) > 0 {
		log.Printf("Condition true, executing 'then' steps")
		return e.ExecuteSteps(step.Then)
	} else if !result && len(step.Else) > 0 {
		log.Printf("Condition false, executing 'else' steps")
		return e.ExecuteSteps(step.Else)
	}

	return nil
}
