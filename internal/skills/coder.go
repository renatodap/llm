package skills

import "github.com/pradord/llm/internal/llm"

// NewCoder creates a programming assistant skill
func NewCoder() *Skill {
	return &Skill{
		Name:        "coder",
		Description: "Programming assistant that writes, explains, and debugs code",
		SystemPrompt: `You are an expert programmer proficient in multiple languages.

Your coding approach:
1. Understand the problem or requirement clearly
2. Choose the best programming language and approach
3. Write clean, well-documented, efficient code
4. Follow best practices and design patterns
5. Explain your code clearly
6. Test and debug when needed

Use web_search to look up documentation and best practices.
Use execute_code to test and verify code when appropriate.`,
		Tools: []string{
			"web_search",
			"execute_code",
		},
		Examples: []string{
			"Write a Python function to find prime numbers",
			"Debug this JavaScript code",
			"Explain how this algorithm works",
		},
		DefaultModel: llm.ModelClaude35Sonnet,
	}
}
