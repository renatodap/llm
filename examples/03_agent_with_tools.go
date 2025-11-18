//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pradord/llm/pkg/agent"
	"github.com/pradord/llm/pkg/llm"
	"github.com/pradord/llm/pkg/skills"
	"github.com/pradord/llm/pkg/tools"
)

// Example 3: Agent with Tools
// Demonstrates multi-step reasoning with tool calling
func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Create LLM client
	client := llm.New(llm.ClientConfig{
		APIKey:       apiKey,
		DefaultModel: llm.ModelClaude35Sonnet,
	})

	// Create agent executor
	executor := agent.NewExecutor(client)

	// Define available tools
	availableTools := []tools.Tool{
		{
			Name:        "web_search",
			Description: "Search the web for current information",
			Parameters: map[string]interface{}{
				"query": "string",
			},
		},
		{
			Name:        "calculator",
			Description: "Perform mathematical calculations",
			Parameters: map[string]interface{}{
				"expression": "string",
			},
		},
	}

	// Use a research skill
	skill := skills.NewResearcher()

	// Task that requires tools
	task := "Research the current price of Bitcoin and calculate how much 5 BTC would be worth"

	fmt.Println("Task:", task)
	fmt.Println("\n--- Agent working ---\n")

	// Execute with agent
	result, err := executor.Run(context.Background(), skill, task, availableTools)
	if err != nil {
		log.Fatalf("Agent execution failed: %v", err)
	}

	fmt.Println("\nResult:", result)
}

/*
USAGE:

1. Set your API key:
   export OPENROUTER_API_KEY=your_key_here

2. Run:
   go run examples/03_agent_with_tools.go

OUTPUT:
Task: Research the current price of Bitcoin and calculate how much 5 BTC would be worth

--- Agent working ---

Step 1: Calling web_search(query="current bitcoin price")
Step 2: Calling calculator(expression="45000 * 5")
Step 3: Generating final answer

Result: Based on my research, Bitcoin is currently trading at approximately $45,000.
Therefore, 5 BTC would be worth $225,000.

HOW IT WORKS:
1. Agent receives task and available tools
2. Agent decides which tools to use
3. Agent calls tools in sequence
4. Agent synthesizes final answer
5. Max steps limit prevents infinite loops

SKILLS AVAILABLE:
- Researcher: Gathers information systematically
- CodeReviewer: Reviews code with best practices
- Writer: Creates well-structured content
- Custom: Define your own skills
*/
