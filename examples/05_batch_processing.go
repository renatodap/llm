//go:build examples

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/pradord/llm/pkg/llm"
)

// Example 5: Batch Processing
// Demonstrates processing multiple prompts efficiently
func main() {
	apiKey := os.Getenv("OPENROUTER_API_KEY")
	if apiKey == "" {
		log.Fatal("OPENROUTER_API_KEY environment variable is required")
	}

	client := llm.New(llm.ClientConfig{
		APIKey:       apiKey,
		DefaultModel: llm.ModelLlama3170B, // Free model for batch processing
	})

	// List of prompts to process
	prompts := []string{
		"Summarize the key benefits of Go programming language",
		"Explain microservices architecture in 2 sentences",
		"What are the best practices for API design?",
		"Describe the CAP theorem briefly",
		"What is the difference between SQL and NoSQL?",
		"Explain Docker containers in simple terms",
		"What are the principles of REST APIs?",
		"Describe the benefits of test-driven development",
	}

	fmt.Printf("Processing %d prompts in parallel...\n\n", len(prompts))

	startTime := time.Now()

	// Process in parallel
	results := make([]string, len(prompts))
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i, prompt := range prompts {
		wg.Add(1)
		go func(idx int, p string) {
			defer wg.Done()

			response, err := client.LLM(context.Background(), p)
			if err != nil {
				log.Printf("Error processing prompt %d: %v", idx, err)
				return
			}

			mu.Lock()
			results[idx] = response
			mu.Unlock()

			fmt.Printf("✓ Completed prompt %d/%d\n", idx+1, len(prompts))
		}(i, prompt)
	}

	wg.Wait()

	duration := time.Since(startTime)

	fmt.Printf("\n--- All prompts processed in %v ---\n\n", duration)

	// Print results
	for i, result := range results {
		if result != "" {
			fmt.Printf("Q%d: %s\n", i+1, prompts[i])
			fmt.Printf("A%d: %s\n\n", i+1, result[:min(100, len(result))]+"...")
		}
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

/*
USAGE:

1. Set your API key:
   export OPENROUTER_API_KEY=your_key_here

2. Run:
   go run examples/05_batch_processing.go

OUTPUT:
Processing 8 prompts in parallel...

✓ Completed prompt 1/8
✓ Completed prompt 3/8
✓ Completed prompt 2/8
✓ Completed prompt 5/8
✓ Completed prompt 4/8
✓ Completed prompt 7/8
✓ Completed prompt 6/8
✓ Completed prompt 8/8

--- All prompts processed in 3.2s ---

Q1: Summarize the key benefits of Go programming language
A1: Go offers excellent performance, built-in concurrency with goroutines, fast compilation times...

[Results for all 8 prompts]

USE CASES:
- Content generation at scale
- Data analysis and summarization
- Code review automation
- Documentation generation
- Translation services
- Sentiment analysis

TIPS:
- Use free models (Llama, Mixtral) for cost efficiency
- Implement rate limiting for API quotas
- Add retry logic for failed requests
- Consider batch size vs concurrency limits
- Monitor token usage and costs
*/
//go:build examples
