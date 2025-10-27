package main

import (
	"encoding/json"
	"os"
)

// Plugin represents the Invoice Radar plugin structure
type Plugin struct {
	Schema        string                 `json:"$schema,omitempty"`
	ID            string                 `json:"id"`
	Name          string                 `json:"name"`
	Description   string                 `json:"description"`
	Homepage      string                 `json:"homepage"`
	ConfigSchema  map[string]ConfigField `json:"configSchema,omitempty"`
	CheckAuth     []Step                 `json:"checkAuth,omitempty"`
	StartAuth     []Step                 `json:"startAuth,omitempty"`
	GetConfigOpts []Step                 `json:"getConfigOptions,omitempty"`
	GetDocuments  []Step                 `json:"getDocuments"`
	Autofill      interface{}            `json:"autofill,omitempty"`
}

// ConfigField represents a configuration field in the plugin
type ConfigField struct {
	Type        string      `json:"type"`
	Title       string      `json:"title"`
	Description string      `json:"description,omitempty"`
	Required    bool        `json:"required,omitempty"`
	Default     interface{} `json:"default,omitempty"`
	Options     []ConfigOpt `json:"options,omitempty"`
}

// ConfigOpt represents a configuration option
type ConfigOpt struct {
	Label string `json:"label"`
	Value string `json:"value"`
}

// Step represents a plugin step
type Step struct {
	Action             string                 `json:"action"`
	URL                string                 `json:"url,omitempty"`
	Selector           string                 `json:"selector,omitempty"`
	Attribute          string                 `json:"attribute,omitempty"`
	Timeout            int                    `json:"timeout,omitempty"`
	Variable           string                 `json:"variable,omitempty"`
	Script             string                 `json:"script,omitempty"`
	Value              string                 `json:"value,omitempty"`
	Fields             map[string]interface{} `json:"fields,omitempty"`
	ForEach            []Step                 `json:"forEach,omitempty"`
	Then               []Step                 `json:"then,omitempty"`
	Else               []Step                 `json:"else,omitempty"`
	Document           interface{}            `json:"document,omitempty"`
	Duration           int                    `json:"duration,omitempty"`
	Base64             string                 `json:"base64,omitempty"`
	Transform          string                 `json:"transform,omitempty"`
	Config             string                 `json:"config,omitempty"`
	Option             interface{}            `json:"option,omitempty"`
	Snippet            string                 `json:"snippet,omitempty"`
	Args               map[string]interface{} `json:"args,omitempty"`
	Iframe             interface{}            `json:"iframe,omitempty"`
	WaitForNetworkIdle bool                   `json:"waitForNetworkIdle,omitempty"`
}

// LoadPlugin loads a plugin from a JSON file
func LoadPlugin(filename string) (*Plugin, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var plugin Plugin
	if err := json.Unmarshal(data, &plugin); err != nil {
		return nil, err
	}

	return &plugin, nil
}
