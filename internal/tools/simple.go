package tools

import (
	"context"
	"encoding/json"

	"github.com/pradord/llm/internal/llm"
)

// SimpleTool allows users to quickly create custom tools with minimal code
type SimpleTool struct {
	name        string
	description string
	parameters  interface{}
	handler     func(ctx context.Context, args map[string]interface{}) (string, error)
	modelType   llm.ModelType
	model       llm.Model
}

// NewSimpleTool creates a simple tool with just a name, description, and handler
func NewSimpleTool(name, description string, handler func(ctx context.Context, args map[string]interface{}) (string, error)) *SimpleTool {
	return &SimpleTool{
		name:        name,
		description: description,
		handler:     handler,
		modelType:   llm.ModelTypeText,
		parameters: map[string]interface{}{
			"type":       "object",
			"properties": map[string]interface{}{},
		},
	}
}

// WithParameters adds parameter schema to the tool
func (st *SimpleTool) WithParameters(params interface{}) *SimpleTool {
	st.parameters = params
	return st
}

// WithModelType sets the model type requirement
func (st *SimpleTool) WithModelType(mt llm.ModelType) *SimpleTool {
	st.modelType = mt
	return st
}

// WithModel sets a specific required model
func (st *SimpleTool) WithModel(m llm.Model) *SimpleTool {
	st.model = m
	return st
}

// Tool interface implementation
func (st *SimpleTool) Name() string {
	return st.name
}

func (st *SimpleTool) Description() string {
	return st.description
}

func (st *SimpleTool) Parameters() interface{} {
	return st.parameters
}

func (st *SimpleTool) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params map[string]interface{}
	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}
	return st.handler(ctx, params)
}

func (st *SimpleTool) RequiredModel() llm.Model {
	return st.model
}

func (st *SimpleTool) ModelType() llm.ModelType {
	return st.modelType
}
