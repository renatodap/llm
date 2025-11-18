//go:build examples

package main

import (
    "context"
    "fmt"

    "github.com/pradord/llm/pkg/tools"
)

// Example: Custom tool using SimpleTool
func main() {
    // Create a simple tool with just a function
    weatherTool := tools.NewSimpleTool(
        "get_weather",
        "Get current weather for a city",
        func(ctx context.Context, args map[string]interface{}) (string, error) {
            city, _ := args["city"].(string)
            // In a real implementation, you'd call a weather API
            return fmt.Sprintf("Weather in %s: Sunny, 72°F", city), nil
        },
    ).WithParameters(map[string]interface{}{
        "type": "object",
        "properties": map[string]interface{}{
            "city": map[string]interface{}{
                "type":        "string",
                "description": "City name",
            },
        },
        "required": []string{"city"},
    })

    // Register it
    registry := tools.NewToolRegistry()
    registry.Register(weatherTool)

    fmt.Println("Custom tool registered:", weatherTool.Name())
}

// ExampleCustomToolAdvanced shows how to create a custom tool implementing the full interface
func ExampleCustomToolAdvanced() {
    // For more complex tools, implement the Tool interface directly
    // See internal/tools/calculator.go for an example
}
