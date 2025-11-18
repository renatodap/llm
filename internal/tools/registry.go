package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pradord/llm/internal/llm"
)

// Tool is the interface all tools must implement
type Tool interface {
	Name() string
	Description() string
	Parameters() interface{} // JSON schema
	Execute(ctx context.Context, args json.RawMessage) (string, error)
	RequiredModel() llm.Model   // Specific model or empty for none
	ModelType() llm.ModelType   // What type of model this needs
}

// ToolRegistry manages all available tools
type ToolRegistry struct {
	tools map[string]Tool
}

// NewToolRegistry creates a registry
func NewToolRegistry() *ToolRegistry {
	return &ToolRegistry{
		tools: make(map[string]Tool),
	}
}

// Register adds a tool (app-defined OR built-in)
func (r *ToolRegistry) Register(tool Tool) {
	r.tools[tool.Name()] = tool
}

// Get retrieves a tool by name
func (r *ToolRegistry) Get(name string) (Tool, bool) {
	tool, ok := r.tools[name]
	return tool, ok
}

// List returns all tool schemas for LLM
func (r *ToolRegistry) List() []map[string]interface{} {
	schemas := make([]map[string]interface{}, 0, len(r.tools))
	for _, tool := range r.tools {
		schemas = append(schemas, map[string]interface{}{
			"type": "function",
			"function": map[string]interface{}{
				"name":        tool.Name(),
				"description": tool.Description(),
				"parameters":  tool.Parameters(),
			},
		})
	}
	return schemas
}

// Execute runs a tool by name
func (r *ToolRegistry) Execute(ctx context.Context, name string, args json.RawMessage) (string, error) {
	tool, ok := r.Get(name)
	if !ok {
		return "", fmt.Errorf("tool not found: %s", name)
	}
	return tool.Execute(ctx, args)
}
