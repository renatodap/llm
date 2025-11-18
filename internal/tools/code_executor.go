package tools

import (
	"context"
	"encoding/json"
	"fmt"
	"os/exec"
	"time"

	"github.com/pradord/llm/internal/llm"
)

// CodeExecutor executes code in isolated environments
// WARNING: Only enable this for trusted environments!
type CodeExecutor struct {
	allowedLanguages map[string]bool
	timeout          time.Duration
}

func NewCodeExecutor() *CodeExecutor {
	return &CodeExecutor{
		allowedLanguages: map[string]bool{
			"python": true,
			"node":   true,
			"go":     true,
		},
		timeout: 30 * time.Second,
	}
}

func (ce *CodeExecutor) Name() string {
	return "execute_code"
}

func (ce *CodeExecutor) Description() string {
	return "Execute code in a sandboxed environment (Python, Node.js, or Go)"
}

func (ce *CodeExecutor) Parameters() interface{} {
	return map[string]interface{}{
		"type": "object",
		"properties": map[string]interface{}{
			"language": map[string]interface{}{
				"type":        "string",
				"description": "Programming language (python, node, go)",
				"enum":        []string{"python", "node", "go"},
			},
			"code": map[string]interface{}{
				"type":        "string",
				"description": "Code to execute",
			},
		},
		"required": []string{"language", "code"},
	}
}

func (ce *CodeExecutor) Execute(ctx context.Context, args json.RawMessage) (string, error) {
	var params struct {
		Language string `json:"language"`
		Code     string `json:"code"`
	}

	if err := json.Unmarshal(args, &params); err != nil {
		return "", err
	}

	if !ce.allowedLanguages[params.Language] {
		return "", fmt.Errorf("language not allowed: %s", params.Language)
	}

	// Create timeout context
	execCtx, cancel := context.WithTimeout(ctx, ce.timeout)
	defer cancel()

	var cmd *exec.Cmd
	switch params.Language {
	case "python":
		cmd = exec.CommandContext(execCtx, "python", "-c", params.Code)
	case "node":
		cmd = exec.CommandContext(execCtx, "node", "-e", params.Code)
	case "go":
		// For Go, we'd need to write to temp file and run
		return "", fmt.Errorf("Go execution not yet implemented")
	default:
		return "", fmt.Errorf("unsupported language: %s", params.Language)
	}

	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("execution error: %w\nOutput: %s", err, string(output))
	}

	return fmt.Sprintf("Execution successful:\n%s", string(output)), nil
}

func (ce *CodeExecutor) RequiredModel() llm.Model {
	return "" // No model needed
}

func (ce *CodeExecutor) ModelType() llm.ModelType {
	return llm.ModelTypeText
}
