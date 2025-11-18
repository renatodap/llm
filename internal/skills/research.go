package skills

import "github.com/pradord/llm/internal/llm"

// NewResearchAssistant creates a research assistant skill
func NewResearchAssistant() *Skill {
	return &Skill{
		Name:        "research_assistant",
		Description: "Deep research on any topic using web search and analysis",
		SystemPrompt: `You are a research assistant skilled at finding, analyzing, and synthesizing information.

Your research process:
1. Break down the research question into sub-questions
2. Search for relevant sources using web_search
3. Read detailed content if URLs are provided
4. Synthesize findings into a comprehensive report
5. Cite all sources

Always verify information from multiple sources and provide clear, well-structured answers.`,
		Tools: []string{
			"web_search",
		},
		Examples: []string{
			"Research the latest developments in quantum computing",
			"What are the health benefits of intermittent fasting?",
			"Explain the current state of AI alignment research",
		},
		DefaultModel: llm.ModelClaude35Sonnet,
	}
}
