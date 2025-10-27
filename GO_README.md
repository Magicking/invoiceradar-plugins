# Invoice Radar Plugin Executor

A simple Golang chromedp application that executes Invoice Radar plugins to automate document fetching from various platforms.

## Overview

This application implements a plugin executor that can:
- Load Invoice Radar plugin JSON configurations
- Execute plugin steps using Chrome DevTools Protocol via chromedp
- Handle authentication flows
- Navigate web pages and extract data
- Download documents (invoices, receipts, etc.)

## Prerequisites

- Go 1.21 or later
- Google Chrome or Chromium installed

## Installation

1. Clone the repository:
```bash
git clone https://github.com/Magicking/invoiceradar-plugins.git
cd invoiceradar-plugins
```

2. Build the application:
```bash
go build -o invoiceradar-plugins .
```

## Usage

### Basic Usage

Run a plugin to fetch documents:

```bash
./invoiceradar-plugins -plugin plugins/plausible.json
```

### Check Authentication Only

Check if you're already authenticated without fetching documents:

```bash
./invoiceradar-plugins -plugin plugins/plausible.json -check-auth
```

### With Configuration

Some plugins require configuration (like region, team ID, etc.). Provide a config file:

```bash
./invoiceradar-plugins -plugin plugins/posthog.json -config config.json
```

Example `config.json`:
```json
{
  "region": "eu",
  "teamId": "12345"
}
```

### Command-Line Options

```
Usage: invoiceradar-plugins -plugin <plugin.json> [-config <config.json>] [-check-auth]

Options:
  -plugin string
        Path to the plugin JSON file (required)
  -config string
        Path to configuration JSON file with plugin settings (optional)
  -check-auth
        Only check if already authenticated, don't run full plugin
```

## How It Works

The application follows this workflow:

1. **Load Plugin**: Reads the JSON plugin configuration file
2. **Check Authentication**: Runs `checkAuth` steps to verify if already logged in
3. **Start Authentication** (if needed): Executes `startAuth` steps, showing the browser for user login
4. **Fetch Documents**: Executes `getDocuments` steps to extract and download documents

## Supported Plugin Actions

The executor supports the following plugin actions:

### Navigation
- `navigate` - Navigate to a URL
- `waitForElement` - Wait for an element to appear
- `waitForURL` - Wait for URL to match pattern
- `waitForNavigation` - Wait for page navigation
- `waitForNetworkIdle` - Wait for network to be idle

### Verification
- `checkElementExists` - Check if element exists on page
- `checkURL` - Verify current URL matches expected

### Interaction
- `click` - Click an element
- `type` - Type text into an input field
- `dropdownSelect` - Select option from dropdown

### Data Extraction
- `extract` - Extract single data value
- `extractAll` - Extract multiple items and iterate

### Document Retrieval
- `downloadPdf` - Download PDF from URL
- `printPdf` - Print page as PDF

### Control Flow
- `if` - Conditional execution based on JavaScript expression
- `sleep` - Wait for specified duration

### JavaScript
- `runJs` - Execute custom JavaScript in page context

## Example Plugins

### Plausible Analytics

```bash
./invoiceradar-plugins -plugin plugins/plausible.json
```

This will:
1. Check if you're logged in to Plausible
2. If not, navigate to login page and wait for you to authenticate
3. Navigate to billing/invoices page
4. Extract all invoices and download them

### PostHog (with configuration)

Create `config.json`:
```json
{
  "region": "eu"
}
```

Run:
```bash
./invoiceradar-plugins -plugin plugins/posthog.json -config config.json
```

## Architecture

The application consists of three main components:

### 1. Plugin Loader (`plugin.go`)
- Defines the plugin JSON schema structure
- Loads and parses plugin files

### 2. Executor (`executor.go`)
- Manages chromedp browser context
- Executes plugin steps sequentially
- Handles variable interpolation (`{{variable}}`)
- Implements all plugin actions

### 3. CLI (`main.go`)
- Parses command-line arguments
- Orchestrates the plugin execution flow
- Provides user feedback and logging

## Development

### Building from Source

```bash
go mod download
go build -o invoiceradar-plugins .
```

### Running Tests

```bash
go test ./...
```

### Adding New Actions

To add support for new plugin actions:

1. Add the action case in `executeStep()` in `executor.go`
2. Implement the action handler method
3. Update this README with the new action

## Plugin Development

See the main [README.md](README.md) for comprehensive documentation on creating Invoice Radar plugins.

Key points:
- Plugins are JSON files following the schema
- Use CSS selectors for element targeting
- Support variable interpolation with `{{variable}}`
- Chain steps using `forEach` for iteration

## Limitations

This is a simplified implementation of the Invoice Radar plugin system. Some features are not fully implemented:

- **Data Extraction**: The `extract` and `extractAll` actions have basic implementations
- **Network Response Extraction**: `extractNetworkResponse` is not implemented
- **PDF Downloads**: Downloads are logged but not actually saved to disk
- **Snippets**: `runSnippet` action is not implemented
- **Pagination**: Not yet supported
- **iFrame handling**: Limited support

For production use, these features would need full implementation.

## Browser Behavior

- The browser runs in **non-headless mode** during authentication to allow user interaction
- The browser window remains visible during execution
- After completion, the browser stays open briefly before closing

## Troubleshooting

### "Plugin file not found"
Ensure the plugin path is correct and the file exists.

### "Element not found"
The website structure may have changed. Check if the selectors in the plugin are still valid.

### "Authentication failed"
Make sure you complete the login process within the timeout period (default: 120 seconds).

### Chrome not found
Ensure Chrome or Chromium is installed and in your system PATH.

## Contributing

Contributions are welcome! Please:
1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Submit a pull request

## License

This project follows the same license as the Invoice Radar plugins repository.

## Related Links

- [Invoice Radar](https://invoiceradar.com/)
- [chromedp Documentation](https://github.com/chromedp/chromedp)
- [Plugin Handbook](README.md)
