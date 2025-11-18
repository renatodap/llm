package skills

import "github.com/pradord/llm/internal/llm"

// NewCodeReviewer creates a code review skill
func NewCodeReviewer() *Skill {
	return &Skill{
		Name:        "code_reviewer",
		Description: "Analyze code quality and provide detailed feedback",
		SystemPrompt: `You are an expert code reviewer who provides constructive, detailed feedback.

Your review process:
1. Analyze code structure and organization
2. Check for potential bugs and edge cases
3. Evaluate performance and efficiency
4. Assess security vulnerabilities
5. Suggest improvements and best practices
6. Provide examples of better approaches

Focus on being helpful and educational, not just critical.`,
		Tools: []string{
			// Code analysis tools would go here
		},
		Examples: []string{
			"Review this Go function for potential issues",
			"Analyze this API endpoint for security vulnerabilities",
			"Suggest improvements for this database query",
		},
		DefaultModel: llm.ModelClaude35Sonnet,
	}
}
