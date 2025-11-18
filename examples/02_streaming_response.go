//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/pradord/llm/pkg/llm"
)

// Example 2: Streaming Response
// Demonstrates real-time streaming of LLM responses
func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := llm.New(llm.ClientConfig{
		APIKey:       apiKey,
		DefaultModel: llm.ModelClaude35Sonnet,
	})

	prompt := "Write a short poem about programming"

	fmt.Println("Question:", prompt)
	fmt.Println("\nAnswer (streaming):")

	// Stream response
	stream, err := client.LLMStream(context.Background(), prompt)
	if err != nil {
		log.Fatalf("Stream failed: %v", err)
	}

	// Read chunks as they arrive
	for chunk := range stream {
		fmt.Print(chunk)
	}

	fmt.Println("\n\n--- Stream complete ---")
}

/*
USAGE:

1. Set your API key:
   export OPENROUTER_API_KEY=your_key_here

2. Run:
   go run examples/02_streaming_response.go

OUTPUT:
Question: Write a short poem about programming

Answer (streaming):
In lines of code, we weave our dreams,
Through loops and functions, logic streams...
[Text appears word-by-word in real-time]

--- Stream complete ---

NOTE:
- Streaming provides real-time feedback
- Useful for chat interfaces
- Lower perceived latency
- Can be cancelled mid-stream with context
*/
