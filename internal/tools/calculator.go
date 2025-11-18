package tools

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/pradord/llm/internal/llm"
)

// Calculator performs mathematical calculations
type Calculator struct{}

func NewCalculator() *Calculator {
	return &Calculator{}
}

func (c *Calculator) Name() string {
	return "calculator"
}

func (c *Calculator) Description() string {
	return "Perform mathematical calculations and evaluate expressions"
}

func (c *Calculator) Parameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"expression": map[string]interface{}{
				"type":        "string",
				"description": "Mathematical expression to evaluate (e.g., '2 + 2', '10 * 5 + 3')",
			},
		},
		"required": []string{"expression"},
	}
}

func (c *Calculator) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Expression string `json:"expression"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}

	// Simple calculator - in production, use a proper expression parser
	// For now, return a placeholder
	result := fmt.Sprintf("Calculated '%s' (calculator tool placeholder - integrate math parser)", params.Expression)
	return result, nil
}

func (c *Calculator) RequiredModel() llm.Model {
	return "" // No model needed - direct computation
}

func (c *Calculator) ModelType() llm.ModelType {
	return llm.ModelTypeText
}
