// Package tools re-exports the internal tools interfaces, registry, and common
// built-in tool constructors so applications can register and use tools
// without importing internal paths.
package tools

import (
    "context"
    i "github.com/pradord/llm/internal/tools"
    l "github.com/pradord/llm/internal/llm"
)

// Re-export core types
type (
    Tool         = i.Tool
    ToolRegistry = i.ToolRegistry
)

// Registry
func NewToolRegistry() *ToolRegistry { return i.NewToolRegistry() }
func (r *ToolRegistry) Register(t Tool) { (*i.ToolRegistry)(r).Register(t) }
func (r *ToolRegistry) Get(name string) (Tool, bool) { return (*i.ToolRegistry)(r).Get(name) }
func (r *ToolRegistry) List() []map[string]interface{} { return (*i.ToolRegistry)(r).List() }

// Built-in tools constructors
func NewWebSearch(client *l.Client, model l.Model) Tool { return i.NewWebSearch(client, model) }
func NewCalculator() Tool { return i.NewCalculator() }
func NewURLFetcher() Tool { return i.NewURLFetcher() }
func NewImageAnalyzer(apiKey string) Tool { return i.NewImageAnalyzer(apiKey) }
func NewImageGenerator(apiKey, baseURL string) Tool { return i.NewImageGenerator(apiKey, baseURL) }
func NewAudioTTS(apiKey string) Tool { return i.NewAudioTTS(apiKey) }
func NewVideoGenerator(apiKey string) Tool { return i.NewVideoGenerator(apiKey) }
func NewSimpleTool(name, description string, fn func(ctx context.Context, args map[string]interface{}) (string, error)) *i.SimpleTool { return i.NewSimpleTool(name, description, fn) }

// Metadata loaders/helpers
func LoadToolMetadataDir(dir string) (map[string]i.ToolMetadata, error) { return i.LoadToolMetadataDir(dir) }
func ApplyToolMetadata(reg *ToolRegistry, metas map[string]i.ToolMetadata) error { return i.ApplyToolMetadata((*i.ToolRegistry)(reg), metas) }
