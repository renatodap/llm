//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pradord/llm/pkg/llm"
)

// Example 1: Simple Chat
// Demonstrates basic LLM usage with a single prompt
func main() {
	// Get API key from environment
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	// Create LLM client
	client := llm.NewClient(llm.ClientConfig{
		APIKey:       apiKey,
		DefaultModel: llm.ModelClaude35Sonnet, // or ModelGPT4o, ModelLlama3170B, etc.
	})

	// Simple prompt
	prompt := "Explain quantum computing in simple terms"

	// Call LLM
	response, err := client.LLM(context.Background(), prompt)
	if err != nil {
		log.Fatalf("LLM call failed: %v", err)
	}

	// Print response
	fmt.Println("Question:", prompt)
	fmt.Println("\nAnswer:", response)

	// Example with different model
	fmt.Println("\n--- Using different model ---")
	response2, err := client.LLM(
		context.Background(),
		"What are the benefits of quantum computing?",
		llm.WithModel(llm.ModelGPT4o),
	)
	if err != nil {
		log.Fatalf("LLM call failed: %v", err)
	}
	fmt.Println("Answer:", response2)
}

/*
USAGE:

1. Set your API key:
   export OPENROUTER_API_KEY=your_key_here

2. Run:
   go run examples/01_simple_chat.go

3. Or use mock mode (no API key needed):
   USE_REAL_LLM=false go run examples/01_simple_chat.go

OUTPUT:
Question: Explain quantum computing in simple terms

Answer: Quantum computing is a revolutionary approach to computation that...
[Full LLM response]

--- Using different model ---
Answer: Quantum computing offers several key benefits...
[Full LLM response]
*/
